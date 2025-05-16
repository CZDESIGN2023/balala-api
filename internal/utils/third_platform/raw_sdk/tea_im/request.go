package tea_im

import (
	"encoding/json"
)

type PlatformType int
type ShowType int //最后一条消息显示内容
type TextDecorate int

const (
	PlatformType_All    PlatformType = 0 //所有平台
	PlatformType_Mobile PlatformType = 1 //1.手机端
	PlatformType_Web    PlatformType = 2 //2.桌面端

	ShowType_Text       ShowType = 0 //0.text
	ShowType_Title      ShowType = 1 //1.title
	ShowType_SubContent ShowType = 2 //2.subContent

	TextDecorate_Underline   TextDecorate = 1
	TextDecorate_Overline    TextDecorate = 1
	TextDecorate_LineThrough TextDecorate = 1
)

type Response struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

type UserInfo struct {
	Id       int    `json:"id"`
	Birthday int    `json:"birthday"`
	Gender   int    `json:"gender"`
	NickName string `json:"nick_name"`
	UserName string `json:"user_name"`
	HeadAddr string `json:"head_addr"`
}

type ChatTokenInfo struct {
	ChatToken  string `json:"chat_token"`
	Sign       string `json:"sign"`
	Id         int64  `json:"id"`
	PfCode     string `json:"pf_code"`
	UserName   string `json:"user_name"`
	UserId     int64  `json:"user_id"`
	HearUrl    string `json:"head_url"`
	Account    string `json:"account"`
	CreateTime int64  `json:"create_time"`
}
type RichItem struct {
	Text       string  `json:"text"`
	TextHeight float64 `json:"text_height,omitempty"`
	FontSize   int     `json:"font_size,omitempty"`
	Decoration int     `json:"decoration,omitempty"`
	LightColor string  `json:"light_color,omitempty"`
	DarkColor  string  `json:"dark_color,omitempty"`
}

// 消息机器人
type RobotMessage struct {
	Icon         string       `json:"icon"`          //图标url地址
	IconSVG      string       `json:"icon_svg"`      //图标svg
	Title        string       `json:"title"`         //主标题 如：任务变更提醒
	SubTitle     string       `json:"sub_title"`     //空间名称
	SubContent   string       `json:"sub_content"`   //任务标题
	Text         string       `json:"text"`          //内容描述
	Url          string       `json:"link_url"`      //地址跳转url地址
	PlatformType PlatformType `json:"platform_type"` //默认选2 0.所以平台 1.手机端 2.桌面端
	Type         ShowType     `json:"type"`          //最后一条消息显示内容 0.text 1.title 2.subContent

	RightRich []RichItem `json:"right_rich"`
	TextRich  []RichItem `json:"text_rich"`
}

type RobotMessageOption func(m *RobotMessage)

func NewRobotMessage(opts ...RobotMessageOption) *RobotMessage {
	m := &RobotMessage{
		PlatformType: PlatformType_All,
		Type:         ShowType_SubContent,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithWebPlatformTypeOption() RobotMessageOption {
	return func(m *RobotMessage) {
		m.PlatformType = PlatformType_Web
	}
}

func WithShowSubContentTypeOption() RobotMessageOption {
	return func(m *RobotMessage) {
		m.Type = ShowType_SubContent
	}
}

func (m *RobotMessage) Clone() *RobotMessage {
	cpy := *m
	return &cpy
}

func (m *RobotMessage) SetShowType(typ ShowType) *RobotMessage {
	m.Type = typ
	return m
}

func (m *RobotMessage) SetIcon(icon string) *RobotMessage {
	m.Icon = icon
	return m
}

func (m *RobotMessage) SetSVGIcon(icon string) *RobotMessage {
	m.IconSVG = icon
	return m
}

func (m *RobotMessage) SetTitle(title string) *RobotMessage {
	m.Title = title
	return m
}

func (m *RobotMessage) SetText(text string) *RobotMessage {
	m.Text = text
	return m
}

func (m *RobotMessage) SetRightRich(items []RichItem) *RobotMessage {
	m.RightRich = items
	return m
}

func (m *RobotMessage) SetTextRich(items []RichItem) *RobotMessage {
	m.TextRich = items

	var text string
	for _, item := range items {
		text += item.Text
	}
	m.Text = text
	return m
}

func (m *RobotMessage) AddTextRich(items ...RichItem) *RobotMessage {
	m.TextRich = append(m.TextRich, items...)

	var text string
	for _, item := range m.TextRich {
		text += item.Text
	}
	m.Text = text
	return m
}

func (m *RobotMessage) SetSubTitle(subTitle string) *RobotMessage {
	m.SubTitle = subTitle
	return m
}

func (m *RobotMessage) SetSubContent(subContent string) *RobotMessage {
	m.SubContent = subContent
	return m
}

func (m *RobotMessage) SetUrl(url string) *RobotMessage {
	m.Url = url
	return m
}

func (m *RobotMessage) Marshal() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}
