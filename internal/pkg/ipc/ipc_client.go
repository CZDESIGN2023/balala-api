package ipc

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/gorilla/websocket"
	"go-cs/internal/conf"
	"go-cs/internal/utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"math/rand"
	"sync"
	"time"
)

// DataHandler 是一个处理数据的函数类型，接受一个频道名和一个消息作为参数。
type DataHandler func(channel string, msg proto.Message)

type IpcClient struct {
	pullConns []*websocket.Conn
	pushConns []*websocket.Conn
	pushUrls  []string
	pullUrls  []string
	//name     string
	outMutex      sync.Mutex
	done          chan bool
	handlersMutex sync.Mutex
	handlers      map[string][]DataHandler // 支持一个channel绑定多个处理函数
	log           *log.Helper
}

// NewIpcClient 连接失败的时候不应该直接退出, 允许某个连接不能用
func NewIpcClient(conf *conf.Data, logger log.Logger) *IpcClient {
	moduleName := "IpcClient"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	i := &IpcClient{
		//name:     name,
		pushUrls:  conf.Ipc.PushUrl,
		pullUrls:  conf.Ipc.PullUrl,
		pushConns: make([]*websocket.Conn, len(conf.Ipc.PushUrl)),
		pullConns: make([]*websocket.Conn, len(conf.Ipc.PullUrl)),
		handlers:  make(map[string][]DataHandler),
		//maxRetry: maxRetry,
		//dataChans: make(map[string][]chan *ChanData),
		done: make(chan bool),
		log:  hlog,
	}
	return i
}

// RegisterHandler 注册回调
func (i *IpcClient) RegisterHandler(channel string, handler DataHandler) {
	i.handlersMutex.Lock()
	defer i.handlersMutex.Unlock()
	i.handlers[channel] = append(i.handlers[channel], handler)
}

// reconnect 当某个连接失败的时候,我们可以调用这个函数重新并替换
func (i *IpcClient) reconnect(ctx context.Context, oldIndex int, urls []string, conns []*websocket.Conn) (*websocket.Conn, error) {
	conn := conns[oldIndex]

	// close old connection
	if conn != nil {
		conn.Close()
	}

	// dial new connection
	newConn, _, err := websocket.DefaultDialer.DialContext(ctx, urls[oldIndex], nil)
	if err != nil {
		i.log.Errorf("connect to:%v failed to reconnect: %v", urls[oldIndex], err)
		return nil, err
	}

	i.log.Infof("IpcClient connect to url:%v", urls[oldIndex])

	// replace old connection with new connection
	conns[oldIndex] = newConn
	return newConn, nil
}

func (i *IpcClient) reconnectForOut(ctx context.Context, connIdx int) {
	// 尝试重新连接
	i.outMutex.Lock()
	defer i.outMutex.Unlock()
	_, err := i.reconnect(ctx, connIdx, i.pushUrls, i.pushConns)
	if err != nil {
		i.log.Errorf("reconnectForOut to:%v failed to reconnect: %v", i.pushUrls[connIdx], err)
	}
}

// SendData 发送数据到ipc当中
func (i *IpcClient) SendData(ctx context.Context, channel string, msg proto.Message) error {
	anyMsg, err := anypb.New(msg)
	if err != nil {
		return err
	}

	// 创建 ChannelMessage
	milliTimestamp := time.Now().UnixNano() / int64(time.Millisecond)
	channelMsg := &ChanData{
		Channel: channel,
		Data:    anyMsg,
		Time:    &milliTimestamp,
	}

	// 序列化 ChannelMessage
	data, err := proto.Marshal(channelMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal ChannelMessage: %v", err)
	}

	// 尝试3次, 失败就直接返回错误
	const MAX_TRY_TIME = 3
	count := MAX_TRY_TIME
	for count > 0 {
		count--

		select {
		case <-i.done:
			// When the done channel is closed, exit the goroutine
			return nil
		default:

			// 选择一个随机的连接进行发送, 这里不用上锁, 只读
			//rand.Seed(time.Now().UnixNano())
			connIdx := rand.Intn(len(i.pushConns))
			conn := i.pushConns[connIdx]
			if conn != nil {
				err = conn.WriteMessage(websocket.BinaryMessage, data)

				// 发送成功
				if err == nil {
					return nil
				}

				i.log.Errorf("Failed to SendData message: %v, %v", channel, err)
			}
			// 重连
			i.reconnectForOut(ctx, connIdx)
			time.Sleep(time.Duration(MAX_TRY_TIME-count) * time.Second)
		}
	}

	return err
}

// 从单个连接当中读取数据, 如果失败则重新连接, 成功则触发回调
// TODO:这边触发必须是一个数组
func (i *IpcClient) handleConnection(ctx context.Context, index int) {
	for {
		select {
		case <-i.done:
			// When the done channel is closed, exit the goroutine
			return
		default:
			// Continue with normal operation
			var err error
			conn := i.pullConns[index]
			if conn == nil {
				conn, err = i.reconnect(ctx, index, i.pullUrls, i.pullConns)
				if err != nil {
					i.log.Errorf("reconnect 从[%v]拉取消息失败:%v", i.pullUrls[index], err)
					continue
				}
			}
			_, data, err := conn.ReadMessage()
			if err != nil {
				i.log.Errorf("reconnect2 从[%v]拉取消息失败:%v", i.pullUrls[index], err)
				// 等待后重试
				time.Sleep(1 * time.Second)
				// 有点土啊, 相同代码写两遍
				i.pullConns[index].Close()
				i.pullConns[index] = nil
				continue
			}

			// 解析读取到的数据
			channelMsg := &ChanData{}
			if err := proto.Unmarshal(data, channelMsg); err != nil {
				i.log.Errorf("Failed to unmarshal message: %v", err)
				continue
			}

			msg, err := channelMsg.Data.UnmarshalNew()
			if err != nil {
				i.log.Errorf("Failed to unmarshal data: %v", err)
				continue
			}

			// 查找对应的处理函数并调用它
			handlers, ok := i.handlers[channelMsg.Channel]
			if ok {
				for _, handler := range handlers {
					handler(channelMsg.Channel, msg)
				}
			}
		}
	}
}

// 启动一个定时器用于测试连接的情况
// CheckAlive 启动一个定时器用于测试连接的情况
func (i *IpcClient) CheckAlive(ctx context.Context) {

	const ChanneclName = "sys.check_alivk"
	// 发送的包号
	checkAliveMessageID := int64(0)
	// 发送的时间戳, 第一次初始化
	checkAliveSentTime := time.Now()
	// 最后接收的时间
	checkLastReceivedTime := time.Now()

	i.RegisterHandler(ChanneclName, func(channel string, msg proto.Message) {
		// 接收到消息时的时间戳
		checkLastReceivedTime = time.Now()
		// 解析消息内容
		timestampMsg, ok := msg.(*TestCheckAlive)
		if ok {
			// 匹配唯一标识符
			if timestampMsg.Random == checkAliveMessageID {

				// 计算延迟时间
				latency := checkLastReceivedTime.Sub(checkAliveSentTime)
				i.log.Debugf("Received message on channel '%s' with latency: %s", channel, latency)
			}
		}
	})

	// 定时发送时间戳消息
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-i.done:
				// When the done channel is closed, exit the goroutine
				return
			case <-ticker.C:
				now := time.Now()
				// 检测一下最后的收包时间
				latency := checkLastReceivedTime.Sub(checkAliveSentTime)
				if latency > time.Duration(60)*time.Second {
					i.log.Errorf("ipc超过60s没有响应")
				}

				// 创建时间戳消息
				timestampProto := &timestamp.Timestamp{
					Seconds: int64(now.Unix()),
					Nanos:   int32(now.Nanosecond()),
				}

				// 设置发送时间和消息ID
				checkAliveSentTime = now
				checkAliveMessageID = rand.Int63()

				timestampMsg := &TestCheckAlive{
					Timestamp: timestampProto,
					Random:    checkAliveMessageID,
				}

				// 发送时间戳消息到频道 "sys.check_alick"
				err := i.SendData(ctx, ChanneclName, timestampMsg)
				if err != nil {
					i.log.Errorf("Failed to send timestamp message: %v", err)
				} else {
					i.log.Debugf("Sent timestamp message on channel 'sys.check_alick'")
				}
			}
		}
	}()
}

// Start 启动了, 初始化连接
func (i *IpcClient) Start(ctx context.Context) {
	for j := 0; j < len(i.pullConns); j++ {
		go i.handleConnection(ctx, j)
	}
}

func (i *IpcClient) Close() {
	// Close done channel to signal other goroutines to exit
	close(i.done)
	for _, conn := range i.pullConns {
		if conn != nil {
			conn.Close()
		}
	}
	// Use the mutex to safely close all connections
	i.outMutex.Lock()
	defer i.outMutex.Unlock()

	for _, conn := range i.pushConns {
		if conn != nil {
			conn.Close()
		}
	}
}
