package notify_message

import (
	"encoding/json"
	"go-cs/api/notify"
	"time"
)

type ObjectType string
type SubjectType string
type ActionType string
type RelationType string

const (
	SubjectType_user SubjectType = "USER"

	ObjectType_space           ObjectType = "SPACE"
	ObjectType_workItem        ObjectType = "WORK_ITEM"
	ObjectType_workItemComment ObjectType = "COMMENT"
	ObjectType_user            ObjectType = "USER"
	ObjectType_workFlow        ObjectType = "WORK_FLOW"
	ObjectType_View            ObjectType = "VIEW"

	ActionType_add    ActionType = "ADD"
	ActionType_edit   ActionType = "EDIT"
	ActionType_delete ActionType = "DELETE"

	Relation_workItemFollower     RelationType = "WORK_ITEM_FOLLOWER"
	Relation_workItemDirector     RelationType = "WORK_ITEM_DIRECTOR"
	Relation_workItemOwner        RelationType = "WORK_ITEM_OWNER"
	Relation_workItemTodo         RelationType = "WORK_ITEM_TODO"
	Relation_workItemCommentAt    RelationType = "COMMENT_AT"
	Relation_workItemCommentRefer RelationType = "COMMENT_REFER"
)

func NewMessage() *Message {
	return &Message{
		Space:    newSpace(),
		Relation: make([]RelationType, 0),
		Notification: &Notification{
			Subject:   newSubject(),
			Object:    newObject(),
			SubObject: newObject(),
			Action:    "",
			Describe:  "",
			Date:      time.Now(),
		},
	}
}

type Message struct {
	//空间信息
	Space *Space `json:"space"`
	//通知对象与事件的人物关系
	Relation []RelationType `json:"relation"`
	//通知类型 type
	Type     notify.Event `json:"type"`
	TypeDesc string       `json:"typeDesc"`
	//通知故事
	Notification *Notification `json:"notification"`
	//关联链接
	RedirectLink string `json:"redirectLink"`
	//是否弹窗通知
	IsPopup bool `json:"isPopup"`
}

func (m *Message) Clone() *Message {
	cpy := *m
	nCpy := *m.Notification
	cpy.Notification = &nCpy
	return &cpy
}

func (m *Message) SetPopup() *Message {
	m.IsPopup = true
	return m
}

func (m *Message) SetType(typ notify.Event) *Message {
	m.Type = typ
	m.TypeDesc = typ.String()
	return m
}

func (m *Message) SetRelation(r RelationType) *Message {
	m.Relation = append(m.Relation, r)
	return m
}

func (m *Message) SetDescribe(desc string) *Message {
	m.Notification.Describe = desc
	return m
}

func (m *Message) SetRedirectLink(link string) *Message {
	m.RedirectLink = link
	return m
}

func (s *Message) String() string {
	r, _ := json.Marshal(s)
	return string(r)
}

// Notification 某人 在 某任务 中 发表评论
type Notification struct {
	//主语 某人 (someone)
	Subject *Subject `json:"subject"`
	//宾语 某事物
	Object *Object `json:"object"`
	//动词 行为
	Action ActionType `json:"action"`
	//补语 补充某事
	SubObject *Object `json:"subObject"`
	//事件信息补充
	Describe string `json:"describe"`
	//时间
	Date time.Time `json:"date"`
}

func newSpace() *Space {
	return &Space{}
}

type Space struct {
	SpaceId   int64  `json:"spaceId"`
	SpaceName string `json:"spaceName"`
}

func newSubject() *Subject {
	return &Subject{}
}

type Subject struct {
	Type SubjectType `json:"type"`
	Data any         `json:"data"`
}

func newObject() *Object {
	return &Object{}
}

type Object struct {
	Type ObjectType `json:"type"`
	Data any        `json:"data"`
}

/** 这里是通用的故事对象信息结构体 **/

type UserData struct {
	Name     string `json:"name"`
	NickName string `json:"nickName"`
	Id       int64  `json:"id"`
	Avatar   string `json:"avatar"`
}

type WorkItemData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Pid  int64  `json:"pid"`
}

type WorkItemCommentData struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Emoji          string `json:"emoji"`
	ReplyCommentId int64  `json:"replyCommentId"`
}

type SpaceData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type WorkFlowData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ViewData struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Key   string `json:"key"`
	Field string `json:"field"`
}
