package api

import (
	"context"
	"go-cs/internal/conf"
	"go-cs/internal/pkg/server3"
	"go-cs/internal/server"
	http_api "go-cs/pkg/http-api"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"
)

var _ = Describe("接口调通测试", func() {

	var (
		ctrl *gomock.Controller
		//mockRegistry *mocks.MockRegistryInterface
		Server3Api server3.Server3Interface
		ctx        context.Context
		c          *conf.Etcd
		t          *durationpb.Duration
	)

	BeforeEach(func() {
		t = &durationpb.Duration{
			Seconds: 30,
		}
		c = &conf.Etcd{
			Endpoints: []string{"10.5.27.115:2379"},
			Timeout:   t,
		}
		registry := server.NewEtcdClient(c)
		registryInterface := http_api.NewRegistryInterface(registry)
		//ctrl = gomock.NewController(GinkgoT())
		//mockRegistry = mocks.NewMockRegistryInterface(ctrl)
		Server3Api = server3.NewServer3Api(registryInterface)
		ctx = context.Background() // 创建一个新的 context
	})

	AfterEach(func() {
		ctrl.Finish()
	})
	Context("GetServiceUrl", func() {
		It("[server3] SendGift接口测试", func() {
			//模拟传参
			fromUserId := 1532537177
			toUserId := 1532537177
			itemId := 2
			count := 2
			row, err := Server3Api.SendGift(ctx, int64(fromUserId), int64(toUserId), int64(itemId), int32(count))
			//现阶段主要是调通 能打通接口 之后要针对接口细节测试再加
			Expect(err).ToNot(HaveOccurred())
			Expect(row).NotTo(Equal(nil))
		})
	})
})
