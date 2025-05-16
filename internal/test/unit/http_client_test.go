package unit

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go-cs/internal/test/mocks"
	. "go-cs/pkg/http-api"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("HTTPClient", func() {
	var (
		ctrl         *gomock.Controller
		mockRegistry *mocks.MockHttpDiscoveryInterface
		httpClient   HTTPClientInterface
		ctx          context.Context
		server       *httptest.Server
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRegistry = mocks.NewMockHttpDiscoveryInterface(ctrl)
		mockRegistry.EXPECT().GetName().Return("test-service")
		httpClient = NewHTTPClient(mockRegistry)
		ctx = context.Background()

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				params := r.URL.Query()
				if val, ok := params["param1"]; ok && val[0] == "value1" {
					w.Write([]byte("Hello, client"))
				} else {
					http.Error(w, "Invalid query param", http.StatusBadRequest)
				}
			case http.MethodPost:
				var body map[string]string
				err := json.NewDecoder(r.Body).Decode(&body)
				if err != nil {
					http.Error(w, "Invalid request body", http.StatusBadRequest)
					return
				}
				if val, ok := body["key"]; ok && val == "value" {
					w.Write([]byte("Hello, client"))
				} else {
					http.Error(w, "Invalid request body", http.StatusBadRequest)
				}
			default:
				http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			}
		}))
	})

	AfterEach(func() {
		ctrl.Finish()
		server.Close()
	})

	Context("DoGet", func() {

		It("should return response body when request succeeds", func() {
			path := "/test?param1=value1"
			mockRegistry.EXPECT().GetServiceUrl(ctx, path).Return(server.URL+path, nil)
			body, err := httpClient.DoGet(ctx, path)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("Hello, client"))
		})

		It("should return error when GetServiceUrl failed", func() {
			path := "/test?param1=value1"
			mockRegistry.EXPECT().GetServiceUrl(ctx, path).Return("", errors.New("get service url failed"))
			_, err := httpClient.DoGet(ctx, path)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("get service url failed"))
		})

		// Similar tests can be added for Get and Post methods
	})

	Context("Get", func() {

		It("should return response body when request succeeds", func() {
			params := map[string]string{"param1": "value1"}
			mockRegistry.EXPECT().GetServiceUrl(ctx, "/test").Return(server.URL+"/test", nil)
			body, err := httpClient.Get(ctx, "/test", &params)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("Hello, client"))
		})

		It("should return error when GetServiceUrl failed", func() {
			params := map[string]string{"param1": "wrong-value"}
			mockRegistry.EXPECT().GetServiceUrl(ctx, "/test").Return(server.URL+"/test", nil)
			_, err := httpClient.Get(ctx, "/test", &params)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Invalid query param"))
		})
	})

	Context("Post", func() {

		It("should return response body when request succeeds", func() {
			postBody := map[string]string{"key": "value"}
			mockRegistry.EXPECT().GetServiceUrl(ctx, "/test").Return(server.URL+"/test", nil)
			body, err := httpClient.Post(ctx, "/test", postBody)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("Hello, client"))
		})

		It("should return error when GetServiceUrl failed", func() {
			postBody := map[string]string{"key": "wrong-value"}
			mockRegistry.EXPECT().GetServiceUrl(ctx, "/test").Return(server.URL+"/test", nil)
			_, err := httpClient.Post(ctx, "/test", postBody)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Invalid request body"))
		})
	})

})
