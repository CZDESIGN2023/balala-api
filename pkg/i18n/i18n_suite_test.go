package i18n_test

import (
	"context"
	"fmt"
	"go-cs/api/comm"
	"go-cs/internal/utils"
	"go-cs/pkg/i18n"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestI18n(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "I18n Suite")
}

var _ = Describe("Add", func() {

	var ctx context.Context
	ctx = context.Background() // 创建一个新的 context

	i18n.InitI18n("./files")

	It("en基本测试", func() {
		Expect(i18n.GetMessage("en", "test.msg_i18n_test")).To(Equal("This is a i18n test!"))
	})

	It("简体测试", func() {
		Expect(i18n.GetMessage("zh-CN", "test.msg_i18n_test")).To(Equal("这是i18n测试!"))
	})

	It("繁体测试", func() {
		Expect(i18n.GetMessage("zh-TW", "test.msg_i18n_test")).To(Equal("這是i18n測試!"))
	})

	It("測試瀏覽器header帶入的accept-language格式", func() {
		Expect(i18n.GetMessage("zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7", "test.msg_i18n_test2")).To(Equal("This is a i18n test2!"))
	})

	It("測試帶變數i18n", func() {
		s := fmt.Sprintf(i18n.GetMessage("en", "test.msg_i18n_test3"), "測試字串測試字串測試字串")
		Expect(s).To(Equal("This is a i18n variable string: 測試字串測試字串測試字串"))
	})

	It("測試帶變數i18n", func() {
		s := fmt.Sprintf(i18n.GetMessage("zh-CN", "test.msg_i18n_test4"), 123456)
		Expect(s).To(Equal("This is a i18n variable number: 123456"))
	})

	It("測試帶變數i18n", func() {
		s := i18n.GetMessage("zh-CN", "test.msg_i18n_test5", 3.1415926)
		//浮點數默認精度為6故自動四捨五入
		Expect(s).To(Equal("这是i18n测试浮点数变数: 3.141593"))
	})

	It("測試錯誤碼i18n", func() {
		resp := utils.NewCommonErrorReply(ctx, comm.ErrorCode_ERROR_MSG_TEST)
		Expect(resp.GetError().GetMessage()).To(Equal("This is a ErrorCode i18n test!"))
	})

	It("測試帶變數錯誤碼i18n", func() {
		resp := utils.NewCommonErrorReply(ctx, comm.ErrorCode_ERROR_MSG_FORMAT_TEST, "測試字串測試字串測試字串")
		Expect(resp.GetError().GetMessage()).To(Equal("This is a ErrorCode i18n variable string: 測試字串測試字串測試字串"))
	})

	It("測試不存在的錯誤碼", func() {
		resp := utils.NewCommonErrorReply(ctx, 999999)
		// 找不到錯誤訊息時，預期結果: 訊息=錯誤碼
		Expect(resp.GetError().GetMessage()).To(Equal("999999"))
	})

	//It("should return an error when adding a negative number", func() {
	//	_, err := Add(-2, 3)
	//	Expect(err).To(HaveOccurred())
	//})
})
