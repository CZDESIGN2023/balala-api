package notify

import (
	"context"
	"fmt"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	msg "go-cs/internal/bean/vo/message/notify_message"
	"go-cs/internal/consts"
	space_domain "go-cs/internal/domain/space"
	wVersion_domain "go-cs/internal/domain/space_work_version"
	user_domain "go-cs/internal/domain/user"
	witem_domain "go-cs/internal/domain/work_item"
	"go-cs/internal/domain/work_item/repo"
	"go-cs/internal/utils"
	"go-cs/internal/utils/third_platform/raw_sdk/tea_im"
	"go-cs/pkg/stream"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cast"
)

var fieldNameMap = map[string]string{
	"directors":      "负责人",  //单独处理
	"tags":           "任务标签", //单独处理
	"workObjectId":   "任务所属模块",
	"workItemType":   "任务类型",
	"workItemName":   "任务名称",
	"describe":       "任务描述",
	"workItemStatus": "任务状态",
	"planTime":       "总排期",
	"processRate":    "任务进度",
	"priority":       "任务优先级",
	"versionId":      "版本",
	"remark":         "交付备注",
	"followers":      "关注人",
}

func buildWorkItemMsg(desc string, event notify.Event, user *user_domain.User, space *space_domain.Space, workItem *witem_domain.WorkItem) *msg.Message {
	m := &msg.Message{
		Space: &msg.Space{
			SpaceId:   space.Id,
			SpaceName: space.SpaceName,
		},
		Relation: make([]msg.RelationType, 0),
		Type:     event,
		TypeDesc: event.String(),
		Notification: &msg.Notification{
			Action: msg.ActionType_edit,
			Subject: &msg.Subject{
				Type: msg.SubjectType_user,
				Data: &msg.UserData{
					Name:     user.UserName,
					NickName: user.UserName,
					Id:       user.Id,
					Avatar:   user.Avatar,
				},
			},
			Object: &msg.Object{
				Type: msg.ObjectType_workItem,
				Data: &msg.WorkItemData{
					Id:   workItem.Id,
					Name: workItem.WorkItemName,
					Pid:  workItem.Pid,
				},
			},
			Describe: desc,
			Date:     time.Now(),
		},
	}

	return m
}

func (s *Notify) ChangeWorkItemField(e *event.ChangeWorkItemField) {
	s.log.Infof("ChangeWorkItemField: %+v", e)

	ctx := context.Background()

	workItem, err := s.workItemRepo.GetWorkItem(ctx, e.WorkItem.Id, &repo.WithDocOption{
		PlanTime:  true,
		Directors: true,
		Followers: true,
	}, nil)
	if err != nil {
		s.log.Error(err)
		return
	}

	// 通知 任务创建人 / 当前节点负责人 / 关注人
	creator, directors, followers, _ := splitUser(
		e.Operator,
		workItem.UserId,
		utils.ToInt64Array(workItem.Doc.Directors),
		utils.ToInt64Array(workItem.Doc.Followers),
	)

	operator, err := s.userRepo.GetUserByUserId(ctx, e.Operator)
	if err != nil {
		s.log.Error(err)
		return
	}

	var oldValue any
	var newValue any

	fieldMap := stream.ToMap(e.Updates, func(i int, v event.FieldUpdate) (string, event.FieldUpdate) {
		return v.Field, v
	})

	var vFiled string
	var needPopup bool
	for field, update := range fieldMap {
		vFiled = field
		switch field {
		case "followers":
			oldValue := update.OldValue.([]int64)
			newValue := update.NewValue.([]int64)
			s.notifyFollowersChange(operator, e.Space, workItem, oldValue, newValue)
			return
		case "planTime":

			oldTimes := update.OldValue.([]any)
			newTimes := update.NewValue.([]any)

			oldValue = parsePlanTime(oldTimes[0], oldTimes[1])
			newValue = parsePlanTime(newTimes[0], newTimes[1])
			needPopup = true

		case "describe", "remark":
			oldValue = utils.ClearRichTextToPlanText(cast.ToString(update.OldValue), false)
			newValue = utils.ClearRichTextToPlanText(cast.ToString(update.NewValue), false)
		case "workObjectId":
			_oldVal := update.OldValue.(int64)
			_newVal := update.NewValue.(int64)
			objectMap, _ := s.workObjectRepo.SpaceWorkObjectMapByObjectIds(ctx, []int64{_oldVal, _newVal})

			oldValue = objectMap[_oldVal].WorkObjectName
			newValue = objectMap[_newVal].WorkObjectName

		case "versionId":
			_oldVal := update.OldValue.(int64)
			_newVal := update.NewValue.(int64)

			objectList, _ := s.workVersionRepo.GetSpaceWorkVersionByIds(ctx, []int64{_oldVal, _newVal})
			objectMap := stream.ToMap(objectList, func(i int, v *wVersion_domain.SpaceWorkVersion) (int64, *wVersion_domain.SpaceWorkVersion) {
				return v.Id, v
			})

			oldValue = objectMap[_oldVal].VersionName
			newValue = objectMap[_newVal].VersionName
		case "processRate":
			oldValue = fmt.Sprintf(`%v`, update.OldValue) + "%"
			newValue = fmt.Sprintf(`%v`, update.NewValue) + "%"
		case "priority":
			oldValue = consts.GetWorkItemPriorityName(cast.ToString(update.OldValue))
			newValue = consts.GetWorkItemPriorityName(cast.ToString(update.NewValue))

		default:
			oldValue = update.OldValue
			newValue = update.NewValue
		}
	}

	buildTemplate := func(filed string) string {
		return fmt.Sprintf(
			`%v 变更了 <br /><p class="ellipsis-desc">%v：“%v” -> “%v”</p>`,
			parseUserTmp(operator),
			fieldNameMap[filed],
			oldValue,
			newValue,
		)
	}

	space := e.Space
	desc := buildTemplate(vFiled)

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, operator, space, workItem)
	m.SetRedirectLink(s.makeRedirectLink(m))
	if needPopup {
		m.SetPopup()
	}

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers))
	utils.Go(func() {
		s.Send2(notifyCtx, toOwner, creator...)
		s.Send2(notifyCtx, toDirector, directors...)
		s.Send2(notifyCtx, toFollower, followers...)
	})

	//推送IM
	needPushQl := stream.ContainsKey(fieldMap, func(field string) bool {
		return slices.Contains([]string{"planTime"}, field)
	})

	if needPushQl {
		//IM
		qlToDirector := tea_im.NewRobotMessage().
			SetTitle("任务提醒").SetSubTitle(space.SpaceName).
			SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
			SetIcon(s.makeIconResLink("", IconRes_Task)).
			SetSVGIcon(IconSVGRes_Task).
			SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

		qlToFollower := tea_im.NewRobotMessage().
			SetTitle("关注任务提醒").SetSubTitle(space.SpaceName).
			SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
			SetIcon(s.makeIconResLink("", IconRes_Follow)).
			SetSVGIcon(IconSVGRes_Follow).
			SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

		utils.Go(func() {
			s.pushThirdPlatformMessage2(notifyCtx, qlToDirector, creator)
			s.pushThirdPlatformMessage2(notifyCtx, qlToDirector, directors)
			s.pushThirdPlatformMessage2(notifyCtx, qlToFollower, followers)
		})
	}
}

func (s *Notify) notifyFollowersChange(operator *user_domain.User, space *space_domain.Space, workItem *witem_domain.WorkItem, oldValue, newValue []int64) {
	removedIds := stream.Diff(oldValue, newValue)
	addedIds := stream.Diff(newValue, oldValue)

	removedIds = stream.Remove(removedIds, operator.Id)
	addedIds = stream.Remove(addedIds, operator.Id)

	addDesc := "你已被添加为 关注人"
	removedDesc := "你已被移除 关注人"

	//组织推送内容
	m := buildWorkItemMsg("", notify.Event_ChangeWorkItemField, operator, space, workItem)
	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()

	toRemoved := m.Clone().SetDescribe(removedDesc)
	toAdded := m.Clone().SetDescribe(addDesc)

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(removedIds, addedIds))
	utils.Go(func() {
		s.Send2(notifyCtx, toRemoved, removedIds...)
		s.Send2(notifyCtx, toAdded, addedIds...)
	})

	//IM
	qlToRemoved := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(removedDesc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToAdded := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(addDesc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, qlToRemoved, removedIds)
		s.pushThirdPlatformMessage2(notifyCtx, qlToAdded, addedIds)
	})
}

func (s *Notify) ChangeWorkItemDirector(e *event.ChangeWorkItemDirector) {
	if e.WorkItem.Pid != 0 {
		s.ChangeChildWorkItemDirector(e) //转交到子任务处理
		return
	}

	s.log.Infof("ChangeWorkItemDirector: %+v", e)

	ctx := context.Background()

	space := e.Space
	workItem := e.WorkItem

	// 通知 任务创建人 / 当前节点负责人 / 关注人
	creators, directors, followers, _ := splitUser(
		e.Operator,
		workItem.UserId,
		utils.ToInt64Array(workItem.Doc.Directors),
		utils.ToInt64Array(workItem.Doc.Followers),
	)

	nodeOldUserIds := stream.Flat(stream.Map(e.Nodes, func(v event.NodeDirectorOp) []int64 {
		return v.OldValues
	}))

	nodeNewUserIds := stream.Flat(stream.Map(e.Nodes, func(v event.NodeDirectorOp) []int64 {
		return v.NewValues
	}))

	nodeOldUserIds = stream.Unique(nodeOldUserIds)
	nodeNewUserIds = stream.Unique(nodeNewUserIds)

	nodeAllUserIds := stream.Concat(nodeOldUserIds, nodeNewUserIds)

	allUserIds := stream.Unique(append(nodeAllUserIds, e.Operator))
	userMap, err := s.userRepo.UserMap(ctx, allUserIds)
	if err != nil {
		return
	}

	operator := userMap[e.Operator]

	// 所有节点的人员变化
	nodeAddUserIds := stream.Of(nodeNewUserIds).Diff(nodeOldUserIds...).Remove(e.Operator).List()
	nodeDelUserIds := stream.Of(nodeOldUserIds).Diff(nodeNewUserIds...).Remove(e.Operator).List()

	diffUserIds := stream.Concat(nodeAddUserIds, nodeDelUserIds)
	creators = stream.Diff(creators, diffUserIds)
	directors = stream.Diff(directors, diffUserIds)
	followers = stream.Diff(followers, diffUserIds)

	var desc string
	var descList []string

	var userAddDescListMap = map[int64][]string{}
	var userDelDescListMap = map[int64][]string{}
	for _, v := range e.Nodes {
		oldUsers := stream.Map(v.OldValues, func(v int64) *user_domain.User {
			return userMap[v]
		})
		newUsers := stream.Map(v.NewValues, func(v int64) *user_domain.User {
			return userMap[v]
		})

		descList = append(descList, fmt.Sprintf(
			`<p class="ellipsis-desc"> “%v”负责人: “%v” -> “%v”</p>`,
			v.NodeName, parseUserTmp(oldUsers...), parseUserTmp(newUsers...),
		))

		for _, userId := range stream.Diff(v.NewValues, v.OldValues) {
			if userId == e.Operator {
				continue
			}
			userAddDescListMap[userId] = append(userAddDescListMap[userId], fmt.Sprintf("你已被添加至 “%v” 负责人", v.NodeName))
		}

		for _, userId := range stream.Diff(v.OldValues, v.NewValues) {
			if userId == e.Operator {
				continue
			}
			userDelDescListMap[userId] = append(userDelDescListMap[userId], fmt.Sprintf("你已被移除 “%v” 负责人", v.NodeName))
		}
	}

	desc = fmt.Sprintf(`%v 变更了 `, parseUserTmp(operator)) + strings.Join(descList, "")
	userAddDescMap := stream.MapValue(userAddDescListMap, func(v []string) string {
		return strings.Join(v, "</br> ")
	})
	userDelDescMap := stream.MapValue(userDelDescListMap, func(v []string) string {
		return strings.Join(v, "</br> ")
	})

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, operator, space, workItem)
	m.SetRedirectLink(s.makeRedirectLink(m))

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)
	toChange := m.Clone().SetRelation(msg.Relation_workItemDirector).SetPopup()

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creators, directors, followers, nodeAddUserIds, nodeDelUserIds))
	utils.Go(func() {

		var userIds []int64

		for userId, desc := range userAddDescMap {
			userIds = append(userIds, userId)
			message := toChange.Clone().SetDescribe(desc)
			s.Send2(notifyCtx, message, userId)
		}
		for userId, desc := range userDelDescMap {
			userIds = append(userIds, userId)
			message := toChange.Clone().SetDescribe(desc)
			s.Send2(notifyCtx, message, userId)
		}

		s.Send2(notifyCtx, toOwner, creators...)
		s.Send2(notifyCtx, toDirector, directors...)
		if !e.ViaCreateWorkItem {
			s.Send2(notifyCtx, toFollower, followers...)
		}
	})

	//IM
	qlToCurChange := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToDirector := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToFollower := tea_im.NewRobotMessage().
		SetTitle("关注任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Follow)).
		SetSVGIcon(IconSVGRes_Follow).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, qlToDirector, creators)
		s.pushThirdPlatformMessage2(notifyCtx, qlToDirector, directors)
		if !e.ViaCreateWorkItem {
			s.pushThirdPlatformMessage2(notifyCtx, qlToFollower, followers)
		}

		for userId, desc := range userAddDescMap {
			message := qlToCurChange.Clone().SetTextRich(parseToImRich(desc))
			s.pushThirdPlatformMessage2(notifyCtx, message, []int64{userId})
		}

		for userId, desc := range userDelDescMap {
			message := qlToCurChange.Clone().SetTextRich(parseToImRich(desc))
			s.pushThirdPlatformMessage2(notifyCtx, message, []int64{userId})
		}
	})
}

func (s *Notify) ChangeChildWorkItemDirector(e *event.ChangeWorkItemDirector) {
	s.log.Infof("ChangeChildWorkItemDirector: %+v", e)

	ctx := context.Background()

	space := e.Space
	workItem := e.WorkItem

	// 通知 任务创建人 / 当前节点负责人 / 关注人
	creator, directors, followers, _ := splitUser(
		e.Operator,
		workItem.UserId,
		utils.ToInt64Array(workItem.Doc.Directors),
		utils.ToInt64Array(workItem.Doc.Followers),
	)

	addUserIds := stream.Diff(e.NewValues, append(e.OldValues, e.Operator))
	delUserIds := stream.Diff(e.OldValues, append(e.NewValues, e.Operator))

	userMap, err := s.userRepo.UserMap(ctx, append(stream.Concat(e.OldValues, e.NewValues), e.Operator))
	if err != nil {
		s.log.Error(err)
		return
	}

	buildTemplate := func(operator *user_domain.User, oldUsers, newUsers []*user_domain.User) string {
		format := "%v 变更了 负责人 <br />“%v” -> “%v”"

		return fmt.Sprintf(
			format,
			parseUserTmp(operator),
			parseUserTmp(oldUsers...),
			parseUserTmp(newUsers...))
	}

	oldUsers := stream.Map(e.OldValues, func(v int64) *user_domain.User {
		return userMap[v]
	})
	newUsers := stream.Map(e.NewValues, func(v int64) *user_domain.User {
		return userMap[v]
	})

	operator := userMap[e.Operator]

	desc := buildTemplate(operator, oldUsers, newUsers)
	addDesc := fmt.Sprintf("你已被添加为 负责人")
	delDesc := fmt.Sprintf("你已被移除 负责人")

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, operator, space, workItem)
	m.SetRedirectLink(s.makeRedirectLink(m))

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)
	toAdd := m.Clone().SetRelation(msg.Relation_workItemDirector).SetDescribe(addDesc).SetPopup()
	toDel := m.Clone().SetRelation(msg.Relation_workItemDirector).SetDescribe(delDesc).SetPopup()

	directors = stream.Diff(directors, addUserIds)

	if s.isTodoWorkItem(workItem) {
		toAdd.SetRelation(msg.Relation_workItemTodo)
		toDel.SetRelation(msg.Relation_workItemTodo)
	}

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers, addUserIds, delUserIds))
	utils.Go(func() {
		s.Send2(notifyCtx, toOwner, creator...)
		s.Send2(notifyCtx, toDirector, directors...)
		s.Send2(notifyCtx, toAdd, addUserIds...)
		s.Send2(notifyCtx, toDel, delUserIds...)

		if !e.ViaCreateWorkItem {
			s.Send2(notifyCtx, toFollower, followers...)
		}
	})

	//IM
	qlToAdd := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetText(addDesc).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToDel := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetText(delDesc).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToDirector := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToFollower := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Follow)).
		SetSVGIcon(IconSVGRes_Follow).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, qlToAdd, addUserIds)
		s.pushThirdPlatformMessage2(notifyCtx, qlToDel, delUserIds)
		s.pushThirdPlatformMessage2(notifyCtx, qlToDirector, creator)
		s.pushThirdPlatformMessage2(notifyCtx, qlToDirector, directors)

		if !e.ViaCreateWorkItem {
			s.pushThirdPlatformMessage2(notifyCtx, qlToFollower, followers)
		}
	})
}

func (s *Notify) SetWorkItemFiles(e *event.SetWorkItemFiles) {
	s.log.Infof("ChangeWorkItemDirector: %+v", e)

	ctx := context.Background()

	space := e.Space
	workItem := e.WorkItem

	// 通知 任务创建人 / 当前节点负责人 / 关注人
	creator, directors, followers, _ := splitUser(
		e.Operator,
		workItem.UserId,
		utils.ToInt64Array(workItem.Doc.Directors),
		utils.ToInt64Array(workItem.Doc.Followers),
	)

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	buildTemplate := func() string {
		if len(e.Adds) > 0 {
			return fmt.Sprintf("%v 添加了 <br />附件：“%v” ", parseUserTmp(operator), e.Adds[0].FileName)
		} else {
			list := stream.Map(e.Deletes, func(v *event.SetWorkItemFiles_FileInfo) string {
				return v.FileName
			})

			return fmt.Sprintf("%v 删除了 <br />附件：“%v” ", parseUserTmp(operator), strings.Join(list, "、"))
		}
	}

	desc := buildTemplate()

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, operator, space, workItem)
	m.SetRedirectLink(s.makeRedirectLink(m))

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers))
	utils.Go(func() {
		s.Send2(notifyCtx, toOwner, creator...)
		s.Send2(notifyCtx, toDirector, directors...)
		s.Send2(notifyCtx, toFollower, followers...)
	})

}

func (s *Notify) CreateChildWorkItem(e *event.CreateChildWorkItem) {
	s.log.Infof("CreateChildWorkItem: %+v", e)

	ctx := context.Background()

	space := e.Space
	workItem := e.WorkItem
	childWorkItem := e.ChildWorkItem

	// 通知 任务创建人 / 当前节点负责人 / 关注人
	creator, directors, followers, _ := splitUser(
		e.Operator,
		workItem.UserId,
		utils.ToInt64Array(workItem.Doc.Directors),
		utils.ToInt64Array(workItem.Doc.Followers),
	)

	childCreator, _, _, _ := splitUser(
		e.Operator,
		childWorkItem.UserId,
		utils.ToInt64Array(childWorkItem.Doc.Directors),
		utils.ToInt64Array(childWorkItem.Doc.Followers),
	)

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	buildTemplate := func() string {
		return fmt.Sprintf("%v 创建了子任务：“%v”", parseUserTmp(operator), childWorkItem.WorkItemName)
	}

	desc := buildTemplate()

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, operator, space, workItem)
	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers, childCreator))

	s.Send2(notifyCtx, toOwner, stream.Diff(creator, childCreator)...)
	s.Send2(notifyCtx, toDirector, stream.Diff(directors, childCreator)...)
	s.Send2(notifyCtx, toFollower, stream.Diff(followers, childCreator)...)

	qlToAdd := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToDirector := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToFollower := tea_im.NewRobotMessage().
		SetTitle("关注任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Follow)).
		SetSVGIcon(IconSVGRes_Follow).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, qlToAdd, stream.Diff(creator, childCreator))
		s.pushThirdPlatformMessage2(notifyCtx, qlToDirector, stream.Diff(directors, childCreator))
		s.pushThirdPlatformMessage2(notifyCtx, qlToFollower, stream.Diff(followers, childCreator))
	})

	s.ChangeChildWorkItemDirector(&event.ChangeWorkItemDirector{
		Event:     notify.Event_ChangeWorkItemDirector,
		Space:     e.Space,
		WorkItem:  childWorkItem,
		Operator:  e.Operator,
		NewValues: utils.ToInt64Array(childWorkItem.Doc.Directors),
	})

	s.notifyFollowersChange(operator, e.Space, childWorkItem, nil, utils.StringArrToInt64Arr(childWorkItem.Doc.Followers))
}

func (s *Notify) CreateWorkItem(e *event.CreateWorkItem) {
	s.log.Infof("CreateWorkItem: %+v", e)

	space := e.Space
	workItem := e.WorkItem

	s.ChangeWorkItemDirector(&event.ChangeWorkItemDirector{
		Event:             notify.Event_ChangeWorkItemDirector,
		Space:             space,
		WorkItem:          workItem,
		Operator:          e.Operator,
		Nodes:             e.Nodes,
		ViaCreateWorkItem: true,
	})

	operator, _ := s.userRepo.GetUserByUserId(context.Background(), e.Operator)

	s.notifyFollowersChange(operator, e.Space, e.WorkItem, nil, utils.StringArrToInt64Arr(e.WorkItem.Doc.Followers))
}

func (s *Notify) deleteChildWorkItem(e *event.DeleteChildWorkItem) {
	s.log.Infof("deleteChildWorkItem: %+v", e)

	ctx := context.Background()

	space := e.Space
	workItem := e.WorkItem
	childWorkItem := e.ChildWorkItem

	// 通知 任务创建人 / 当前节点负责人 / 关注人
	operators := []int64{e.Operator}
	followers := stream.Of(utils.ToInt64Array(workItem.Doc.Followers)).
		Diff(operators...).
		List()
	creators := stream.Of([]int64{e.WorkItem.UserId}).
		Diff(operators...).Diff(followers...).
		List()
	directors := stream.Of(utils.ToInt64Array(workItem.Doc.Directors)).
		Diff(creators...).Diff(followers...).Diff(operators...).
		List()

	childFollowers := stream.Of(utils.ToInt64Array(childWorkItem.Doc.Followers)).
		Diff(operators...).
		List()
	childCreator := stream.Of([]int64{childWorkItem.UserId}).
		Diff(operators...).Diff(childFollowers...).
		List()
	childDirectors := stream.Of(utils.ToInt64Array(childWorkItem.Doc.Directors)).
		Diff(operators...).Diff(childFollowers...).Diff(childCreator...).
		List()

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	buildTemplate := func() string {
		return fmt.Sprintf("%v 删除了子任务：“%s”", parseUserTmp(operator), childWorkItem.WorkItemName)
	}

	desc := buildTemplate()

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, operator, space, childWorkItem)
	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)

	toChildOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toChildDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)
	toChildFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creators, directors, followers, childCreator, childDirectors, childFollowers))
	utils.Go(func() {
		s.Send2(notifyCtx, toChildOwner, childCreator...)
		s.Send2(notifyCtx, toChildDirector, childDirectors...)
		s.Send2(notifyCtx, toChildFollower, childFollowers...)

		s.Send2(notifyCtx, toOwner, stream.Diff(creators, childCreator)...)
		s.Send2(notifyCtx, toDirector, stream.Diff(directors, childDirectors)...)
		s.Send2(notifyCtx, toFollower, stream.Diff(followers, childFollowers)...)
	})

	qlToAdd := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToDirector := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToFollower := tea_im.NewRobotMessage().
		SetTitle("关注任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Follow)).
		SetSVGIcon(IconSVGRes_Follow).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, qlToAdd, childCreator)
		s.pushThirdPlatformMessage2(notifyCtx, qlToDirector, childDirectors)
		s.pushThirdPlatformMessage2(notifyCtx, qlToFollower, childFollowers)

		s.pushThirdPlatformMessage2(notifyCtx, qlToAdd, stream.Diff(creators, childCreator))
		s.pushThirdPlatformMessage2(notifyCtx, qlToDirector, stream.Diff(directors, childDirectors))
		s.pushThirdPlatformMessage2(notifyCtx, qlToFollower, stream.Diff(followers, childFollowers))
	})
}

func (s *Notify) DeleteWorkItem(e *event.DeleteWorkItem) {
	s.CooperateDeleteWorkItem(e) //协作通知

	s.log.Infof("DeleteWorkItem: %+v", e)

	if e.WorkItem.Pid != 0 {

		parent, _ := s.workItemRepo.GetWorkItem(context.Background(), e.WorkItem.Pid, nil, nil)
		s.deleteChildWorkItem(&event.DeleteChildWorkItem{
			Event:         notify.Event_DeleteWorkItem,
			Operator:      e.Operator,
			Space:         e.Space,
			WorkItem:      parent,
			ChildWorkItem: e.WorkItem,
		})
		return
	}

	space := e.Space
	workItem := e.WorkItem

	userMap, err := s.userRepo.UserMap(context.Background(), []int64{e.Operator})
	if err != nil {
		return
	}

	operator := userMap[e.Operator]
	desc := fmt.Sprintf("%s已删除任务", parseUserTmp(operator))

	m := buildWorkItemMsg(desc, e.Event, operator, space, workItem)
	m.SetPopup()

	operators := []int64{e.Operator}
	followers := stream.Of(utils.ToInt64Array(workItem.Doc.Followers)).
		Diff(operators...).
		List()
	creators := stream.Of([]int64{workItem.UserId}).
		Diff(operators...).Diff(followers...).
		List()

	participators := utils.ToInt64Array(workItem.Doc.Participators)
	if len(workItem.WorkItemFlowRoles) > 0 {
		for _, v := range workItem.WorkItemFlowRoles {
			participators = append(participators, v.Directors.ToInt64s()...)
		}
		participators = stream.Unique(participators)
	}
	participators = stream.Of(participators).
		Diff(operators...).Diff(followers...).Diff(creators...).
		List()

	toCreator := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toParticipator := m.Clone().SetRelation(msg.Relation_workItemDirector)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creators, participators, followers))
	s.Send2(notifyCtx, toCreator, creators...)
	s.Send2(notifyCtx, toParticipator, participators...)
	s.Send2(notifyCtx, toFollower, followers...)

	qlToAdd := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToDirector := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	qlToFollower := tea_im.NewRobotMessage().
		SetTitle("关注任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Follow)).
		SetSVGIcon(IconSVGRes_Follow).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, qlToAdd, creators)
		s.pushThirdPlatformMessage2(notifyCtx, qlToDirector, participators)
		s.pushThirdPlatformMessage2(notifyCtx, qlToFollower, followers)
	})

	//发送子任务删除信息
	if len(e.SubWorkItems) > 0 {
		for _, subWorkItem := range e.SubWorkItems {
			s.deleteChildWorkItem(&event.DeleteChildWorkItem{
				Event:         notify.Event_DeleteWorkItem,
				Operator:      e.Operator,
				Space:         e.Space,
				WorkItem:      e.WorkItem,
				ChildWorkItem: subWorkItem,
			})
		}
	}
}

func (s *Notify) ChangeWorkItemTag(e *event.ChangeWorkItemTag) {
	s.log.Infof("ChangeWorkItemTag: %+v", e)

	ctx := context.Background()

	// 通知 任务创建人 / 当前节点负责人 / 关注人
	creator, directors, followers, _ := splitUser(
		e.Operator,
		e.WorkItem.UserId,
		utils.ToInt64Array(e.WorkItem.Doc.Directors),
		utils.ToInt64Array(e.WorkItem.Doc.Followers),
	)

	allTagIds := stream.Unique(stream.Concat(e.OldValues, e.NewValues))

	userMap, _ := s.userRepo.UserMap(ctx, []int64{e.Operator})
	tagMap, _ := s.tagRepo.TagMap(ctx, allTagIds)

	operator := userMap[e.Operator]

	buildTemplate := func() string {
		oldTagNames := stream.Concat(stream.Map(e.OldValues, func(v int64) string {
			return tagMap[v].TagName
		}))
		newTagNames := stream.Concat(stream.Map(e.NewValues, func(v int64) string {
			return tagMap[v].TagName
		}))
		return fmt.Sprintf(
			"%v 变更了 任务标签：“%v” -> “%v”",
			parseUserTmp(operator),
			strings.Join(oldTagNames, "、"),
			strings.Join(newTagNames, "、"),
		)
	}

	workItem := e.WorkItem
	space := e.Space
	desc := buildTemplate()

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, operator, space, workItem)
	m.SetRedirectLink(s.makeRedirectLink(m))

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers))
	s.Send2(notifyCtx, toOwner, creator...)
	s.Send2(notifyCtx, toDirector, directors...)
	s.Send2(notifyCtx, toFollower, followers...)
}

func (s *Notify) WorkItemExpired(e *event.WorkItemExpired) {
	s.log.Infof("WorkItemExpired: %+v", e)

	ctx := context.Background()

	workItem, err := s.workItemRepo.GetWorkItem(ctx, e.WorkItemId, &repo.WithDocOption{
		Directors: true,
		Followers: true,
	}, nil)
	if err != nil {
		return
	}

	// 仅在工作日推送
	spaceCfg, _ := s.spaceRepo.GetSpaceConfig(ctx, workItem.SpaceId)
	if spaceCfg == nil || !spaceCfg.WorkingDay.IsWorkingDay(time.Now().Weekday()) {
		return
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil {
		return
	}

	buildTemplate := func() string {
		return fmt.Sprintf("任务已逾期")
	}

	desc := buildTemplate()

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, &user_domain.User{}, space, workItem)
	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()

	//关注人/当前负责人/创建人
	userIds := stream.Concat(utils.ToInt64Array(workItem.Doc.Directors), utils.ToInt64Array(workItem.Doc.Followers), []int64{workItem.UserId})
	userIds = stream.Remove(stream.Unique(userIds), e.Operator)

	notifyCtx := s.buildNotifyCtx(space, userIds)
	utils.Go(func() {
		s.Send2(notifyCtx, m, userIds...)
	})

	imMsg := tea_im.NewRobotMessage().
		SetShowType(tea_im.ShowType_Text).
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetText(desc).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMsg, userIds)
	})
}

func (s *Notify) WorkItemFlowNodeExpired(e *event.WorkItemFlowNodeExpired) {
	s.log.Infof("WorkItemFlowNodeExpired: %+v", e)

	ctx := context.Background()

	workItem, err := s.workItemRepo.GetWorkItem(ctx, e.WorkItemId, &repo.WithDocOption{
		Followers: true,
	}, nil)
	if err != nil {
		return
	}

	// 仅在工作日推送
	spaceCfg, _ := s.spaceRepo.GetSpaceConfig(ctx, workItem.SpaceId)
	if !spaceCfg.WorkingDay.IsWorkingDay(time.Now().Weekday()) {
		return
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil {
		return
	}

	buildTemplate := func() string {
		switch {
		case workItem.IsWorkFlowMainTask():
			return fmt.Sprintf("[%v]节点已逾期", e.NodeName)
		case workItem.IsStateFlowMainTask():
			return fmt.Sprintf("[%v]状态已逾期", e.NodeName)
		default:
			return ""
		}
	}

	desc := buildTemplate()

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, &user_domain.User{}, space, workItem)
	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()

	//关注人 /节点负责人
	userIds := stream.Concat(e.NodeDirectors, utils.ToInt64Array(workItem.Doc.Followers))
	userIds = stream.Remove(stream.Unique(userIds), e.Operator)

	notifyCtx := s.buildNotifyCtx(space, userIds)
	utils.Go(func() {
		s.Send2(notifyCtx, m, userIds...)
	})

	imMsg := tea_im.NewRobotMessage().
		SetShowType(tea_im.ShowType_Text).
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetText(desc).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMsg, userIds)
	})
}

func (s *Notify) RemindWork(e *event.RemindWork) {
	s.log.Infof("RemindWork: %+v", e)
	space := e.Space
	workItem := e.WorkItem

	userMap, _ := s.userRepo.UserMap(context.Background(), []int64{e.Operator})
	operator := userMap[e.Operator]

	buildTemplate := func() string {
		workItemName := e.WorkItem.WorkItemName
		return fmt.Sprintf("%s 催了你尽快处理：“%s”", parseUserTmp(operator), workItemName)
	}

	desc := buildTemplate()

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, operator, space, e.WorkItem)
	m.SetRelation(msg.Relation_workItemTodo)
	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()

	notifyCtx := s.buildNotifyCtx(space, e.TargetIds)
	utils.Go(func() {
		s.Send2(notifyCtx, m, e.TargetIds...)
	})

	imMsg := tea_im.NewRobotMessage().
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetSubContent(workItem.WorkItemName).SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeWorkItemRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMsg, e.TargetIds)
	})
}
