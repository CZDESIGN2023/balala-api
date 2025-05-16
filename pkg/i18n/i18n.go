package i18n

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var defaultLanguage string = "zh-CN"
var bundle *i18n.Bundle
var supported []language.Tag

type ctxI18nLanguageKey struct{}

func InitI18n(i18nPath string) {

	bundle = i18n.NewBundle(language.Make(defaultLanguage))
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	//搜索文件
	fss, err := os.ReadDir(i18nPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	////載入支持的語言包
	for _, v := range fss {
		if !v.IsDir() {
			fsName := v.Name()
			ext := filepath.Ext(fsName)
			language, _ := strings.CutSuffix(v.Name(), ext)
			LoadMessageFile(language, filepath.Join(i18nPath, fsName))
		}
	}
}

// LoadMessageFile 載入指定語言和檔案路徑的訊息檔案。
func LoadMessageFile(lang, filePath string) error {
	tag := language.MustParse(lang)
	supported = append(supported, tag)
	messageData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = bundle.ParseMessageFileBytes(messageData, lang+".toml")
	if err != nil {
		fmt.Println(err)
	}

	return err
}

/*
GetMessage 輸入語系取得i18n文字訊息
msgKey:
依照toml分類的標題串接底下的key
例如[info]類底下的key是msg_success, 應該傳入的msgKey即為:info.msg_success

fmtValues:
訊息格式化時傳入的資料
使用fmt.Sprintf格式化訊息
訊息內容範例: "Hello %s, you are %d years old."
*/
func GetMessage(lang, msgKey string, fmtValues ...interface{}) string {
	localizer := i18n.NewLocalizer(bundle, lang)
	message, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: msgKey},
		TemplateData:   nil,
	})
	if len(fmtValues) == 0 || message == "" {
		return message
	} else {
		fmtMessage := fmt.Sprintf(message, fmtValues...)
		return fmtMessage
	}
}

/*
GetMessageFmtData 輸入語系取得i18n文字訊息

msgKey:
依照toml分類的標題串接底下的key
例如[info]類底下的key是msg_success, 應該傳入的msgKey即為:info.msg_success

fmtData:
訊息格式化時傳入的資料

使用i18n內建的資料樣板格式化訊息
訊息內容範例: "Hello {{.Name}}, you are {{.Age}} years old."
可以傳入map:

	data := map[string]interface{}{
	    "Name": "John",
		"Age": 30,
	}

或是也可以傳入自訂的struct

	type Person struct {
	   Name string
	   Age  int
	}

	data := Person{
	    Name: "John",
	    Age: 30,
	}
*/
func GetMessageFmtData(lang, msgKey string, fmtData interface{}) string {
	localizer := i18n.NewLocalizer(bundle, lang)
	message, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: msgKey},
		TemplateData:   fmtData,
	})
	return message
}

// GetLanguage 取得ctx中語系資訊
func GetLanguage(ctx context.Context) string {

	if lang, ok := ctx.Value(ctxI18nLanguageKey{}).(string); ok {
		return lang
	}

	tr, ok := transport.FromServerContext(ctx) // ctx.Value("req").(*http.Request)
	if !ok {
		// 查無語系設定，可忽略此錯誤，i18N組件會回應預設語系
		return ""
	}

	lang := tr.RequestHeader().Get("Accept-Language")
	if lang == "" {
		// http2的header規定是小寫, 多判斷一次兼容性比較好
		lang = tr.RequestHeader().Get("accept-language")
	}
	return lang
}
