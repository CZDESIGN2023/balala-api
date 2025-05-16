// ipc_client_test.go
package unit

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go-cs/internal/conf"
	. "go-cs/internal/pkg/ipc"
	"google.golang.org/protobuf/proto"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

var _ = Describe("IpcClient", func() {
	var (
		ctrl    *gomock.Controller
		s       *httptest.Server
		confIpc conf.Data_IPC
		c       *IpcClient
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		//mock websocket卡在netconn的Write和Read，先放棄
		//mc = mocks.NewMockConn(ctrl)
		//wc = mocks.NewMockWebConn(ctrl)
		//s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	c, _ := upgrader.Upgrade(w, r, nil)
		//	_ = c.UnderlyingConn().Close() // We immediately close the real connection
		//	websocket.DefaultDialer = &websocket.Dialer{
		//		NetDial: func(network, addr string) (net.Conn, error) {
		//			return mc, nil // and replace it with our mock
		//		}}
		//}))
		var upgrader = websocket.Upgrader{}
		var clients = make(map[*websocket.Conn]bool)
		var broadcast = make(chan []byte)
		var messageCache = make([][]byte, 0)
		//用httptest模擬一個websocket連線
		s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, _ := upgrader.Upgrade(w, r, nil)
			clients[c] = true
			fmt.Println("當前server存儲的信息數", len(messageCache))
			// 当一个新的连接建立时，将缓存的消息发送给这个新的连接
			for _, msg := range messageCache {
				if err := c.WriteMessage(websocket.BinaryMessage, msg); err != nil {
					return
				}
				fmt.Println("--------send message from server in cache--------")
			}
			go func() {
				defer func() {
					delete(clients, c)
					c.Close()
				}()
				for {
					//客戶端有信息傳進來時讀取
					_, msg, err := c.ReadMessage()
					if err != nil {
						return
					}
					fmt.Println("--------receive message from client--------")
					broadcast <- msg
				}
			}()
		}))
		go func() {
			for {
				msg := <-broadcast
				// 将新的消息加入到缓存中
				messageCache = append(messageCache, msg)
				for client := range clients {
					//向客戶端發送消息
					err := client.WriteMessage(websocket.BinaryMessage, msg)
					if err != nil {
						fmt.Printf("error occurred while writing message to client: %v", err)
						client.Close()
						delete(clients, client)
					}
					fmt.Println("--------send message from server--------")
				}
			}
		}()
	})

	AfterEach(func() {
		ctrl.Finish()
		s.Close()
	})

	Context("when sending a message", func() {
		It("should send the message to the server", func() {
			//设置ipc模拟地址
			confIpc = conf.Data_IPC{
				PushUrl: []string{strings.Replace(s.URL, "http://", "ws://", 1)},
				PullUrl: []string{strings.Replace(s.URL, "http://", "ws://", 1)},
			}
			logger := log.NewStdLogger(os.Stdout)
			logger = log.With(logger, "timestamp", log.DefaultTimestamp)
			c = NewIpcClient(&conf.Data{Ipc: &confIpc}, logger)
			//传送信息
			err := c.SendData(context.Background(), "test", &TestChatMessage{
				UserId:  "test",
				Content: "data",
			})
			Expect(err).To(BeNil())
		})
	})

	Context("when receiving a message", func() {
		It("should call the registered handler", func() {
			//设置ipc模拟地址
			confIpc = conf.Data_IPC{
				PushUrl: []string{strings.Replace(s.URL, "http://", "ws://", 1)},
				PullUrl: []string{strings.Replace(s.URL, "http://", "ws://", 1)},
			}
			logger := log.NewStdLogger(os.Stdout)
			logger = log.With(logger, "timestamp", log.DefaultTimestamp)
			c = NewIpcClient(&conf.Data{Ipc: &confIpc}, logger)

			handlerCalled := false
			handler := func(channel string, msg proto.Message) {
				handlerCalled = true
				fmt.Print("收到消息: ", msg)
			}
			testChatMessage := &TestChatMessage{
				UserId:  "test",
				Content: "data",
			}
			//传送信息
			err := c.SendData(context.Background(), "test", testChatMessage)
			Expect(err).To(BeNil())
			c.RegisterHandler("test", handler)
			c.Start(context.Background())
			Eventually(func() bool { return handlerCalled }).Should(BeTrue())
		})
	})
})
