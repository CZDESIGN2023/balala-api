# unit/ 目录下包含了单元测试的代码，按照功能模块划分子目录
# integration/ 目录下包含了集成测试的代码，按照功能模块划分子目录
# 每个测试文件中，可以使用 Ginkgo 和 Gomega 等测试框架来编写测试用例和断言逻辑。
## 使用ginkgo -r进行测试
## 使用mockgen 生成接口mockgen方法 
`mockgen -destination=mocks/mock_http_client.go -source=../pkg/server3/http_client.go`