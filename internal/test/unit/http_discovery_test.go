package unit

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go-cs/internal/test/mocks"
	. "go-cs/pkg/http-api"
)

var _ = Describe("HttpDiscovery", func() {
	var (
		ctrl          *gomock.Controller
		mockRegistry  *mocks.MockRegistryInterface
		httpDiscovery HttpDiscoveryInterface
		serviceName   string
		ctx           context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRegistry = mocks.NewMockRegistryInterface(ctrl)
		// 测试服务的名称
		serviceName = "test-service"
		httpDiscovery = NewHttpDiscovery(mockRegistry, serviceName)
		ctx = context.Background() // 创建一个新的 context
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("GetServiceUrl", func() {
		It("should return correct url when there is at least one service instance", func() {
			mockInstances := []*registry.ServiceInstance{
				{
					Name:      serviceName,
					Endpoints: []string{"http://localhost:8080"},
				},
			}
			mockRegistry.EXPECT().GetService(ctx, serviceName).Return(mockInstances, nil)
			url, err := httpDiscovery.GetServiceUrl(context.Background(), "/test")
			Expect(err).To(BeNil())
			Expect(url).To(Equal("http://localhost:8080/test"))
		})

		It("should return error when no service instances available", func() {
			mockRegistry.EXPECT().GetService(ctx, serviceName).Return(nil, nil)
			_, err := httpDiscovery.GetServiceUrl(context.Background(), "/test")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("no test-service service instances available"))
		})

		It("should return error when GetService failed", func() {
			mockRegistry.EXPECT().GetService(ctx, serviceName).Return(nil, errors.New("get service failed"))
			_, err := httpDiscovery.GetServiceUrl(context.Background(), "/test")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("get service failed"))
		})
	})

	//Context("getInstances", func() {
	//	// 使用反射获取未导出的方法并调用它
	//	value := reflect.ValueOf(httpDiscovery)
	//	method := value.MethodByName("getInstances")
	//	result := method.Call([]reflect.Value{reflect.ValueOf(context.Background())})[0].Interface().(string)
	//	Expect(result).To(Equal("http://localhost:8080"))
	//})
})
