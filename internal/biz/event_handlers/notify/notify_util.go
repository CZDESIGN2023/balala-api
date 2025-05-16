package notify

import (
	"context"
	"fmt"
	"go-cs/internal/consts"
	"go-cs/internal/utils/third_platform/raw_sdk/tea_im"
	"go-cs/pkg/stream"
	"regexp"
	"strings"
	"time"

	space_domain "go-cs/internal/domain/space"
	member_domain "go-cs/internal/domain/space_member"
	user_domain "go-cs/internal/domain/user"
)

func splitUser(operator, ownerId int64, directors, followers []int64) (creators, director, follower, all []int64) {
	operators := []int64{operator}
	follower = stream.Of(followers).
		Diff(operators...).
		List()
	creators = stream.Of([]int64{ownerId}).
		Diff(operators...).Diff(followers...).
		List()
	director = stream.Of(directors).
		Diff(operators...).Diff(followers...).Diff(creators...).
		List()

	all = stream.Unique(stream.Concat(creators, director, follower))

	return
}

func parsePlanTime(start, end any) string {
	startI := start.(int64)
	endI := end.(int64)

	if startI == 0 && endI == 0 {
		return "无"
	}

	startTime := time.Unix(startI, 0)
	endTime := time.Unix(endI, 0)

	layout := "2006/01/02"

	startStr := startTime.Format(layout)
	endStr := endTime.Format(layout)

	return fmt.Sprintf("%s ~ %s", startStr, endStr)
}

func parseUserTmp(users ...*user_domain.User) string {
	if len(users) == 0 {
		return ""
	}

	s := stream.Map(users, func(u *user_domain.User) string {
		return fmt.Sprintf(`%s<span class="minor-color">（%s）</span>`, u.UserNickname, u.UserName)
	})

	return strings.Join(s, "、")
}

var htmlTagRegx = regexp.MustCompile("<.*?>")

func cleanHtmlTag(s string) string {
	s = strings.Replace(s, "</p>", "\n", -1)
	s = strings.Replace(s, "<br />", "\n", -1)
	s = htmlTagRegx.ReplaceAllString(s, "")
	return s
}

func MinorColorSpan(content string) string {
	return fmt.Sprintf("<span class=\"minor-color\">%s</span>", content)
}

var minorColorRegx = regexp.MustCompile(`(?i)<span class="minor-color">.*?</span>`)

func parseToImRich(text string) []tea_im.RichItem {
	// 找到所有用户名

	indexList := minorColorRegx.FindAllStringIndex(text, -1)

	curIdx := 0
	endIdx := len(text)

	var items []tea_im.RichItem
	for _, v := range indexList {
		if curIdx != v[0] {
			items = append(items, tea_im.RichItem{
				Text:     cleanHtmlTag(text[curIdx:v[0]]),
				FontSize: 14,
			})
		}
		items = append(items, tea_im.RichItem{
			Text:       cleanHtmlTag(text[v[0]:v[1]]),
			LightColor: "0x999999",
			DarkColor:  "0x999999",
			FontSize:   14,
		})
		curIdx = v[1]
	}

	if curIdx != endIdx {
		items = append(items, tea_im.RichItem{
			Text:     cleanHtmlTag(text[curIdx:endIdx]),
			FontSize: 14,
		})
	}

	return items
}

func parseUserTmp2(users ...*user_domain.User) string {
	if len(users) == 0 {
		return ""
	}

	s := stream.Map(users, func(u *user_domain.User) string {
		return fmt.Sprintf(`%s（%s）`, u.UserNickname, u.UserName)
	})

	return strings.Join(s, "、")
}

type notifyCtx struct {
	space                            *space_domain.Space
	forceNotify                      bool
	memberMap                        map[int64]*member_domain.SpaceMember
	userNotifySwitchGlobalMap        map[int64]*user_domain.UserConfig
	userNotifySwitchThirdPlatformMap map[int64]*user_domain.UserConfig
	userNotifySwitchSpaceMap         map[int64]*user_domain.UserConfig
}

type NotifyCtxOpt func(*notifyCtx)

func WithNotifyCtxForceNotifyOpt() NotifyCtxOpt {
	return func(ctx *notifyCtx) {
		ctx.forceNotify = true
	}
}

func (s *Notify) buildNotifyCtx(space *space_domain.Space, userIds []int64, opts ...NotifyCtxOpt) *notifyCtx {
	memberMap, _ := s.memberRepo.SpaceMemberMapByUserIds(context.Background(), space.Id, userIds)

	notifyMap, _ := s.userRepo.GetUserConfigMapByUserIdsAndKeys(context.Background(), userIds, []string{
		consts.UserConfigKey_NotifySwitchSpace,
		consts.UserConfigKey_NotifySwitchThirdPlatform,
		consts.UserConfigKey_NotifySwitchGlobal,
	})

	userGlobalNotifyMap := stream.MapValue(notifyMap, func(v map[string]*user_domain.UserConfig) *user_domain.UserConfig {
		return v[consts.UserConfigKey_NotifySwitchGlobal]
	})
	userThirdPlatformNotifyConfigMap := stream.MapValue(notifyMap, func(v map[string]*user_domain.UserConfig) *user_domain.UserConfig {
		return v[consts.UserConfigKey_NotifySwitchThirdPlatform]
	})
	userSpaceNotifyConfigMap := stream.MapValue(notifyMap, func(v map[string]*user_domain.UserConfig) *user_domain.UserConfig {
		return v[consts.UserConfigKey_NotifySwitchSpace]
	})

	ctx := &notifyCtx{
		space:                            space,
		memberMap:                        memberMap,
		userNotifySwitchGlobalMap:        userGlobalNotifyMap,
		userNotifySwitchThirdPlatformMap: userThirdPlatformNotifyConfigMap,
		userNotifySwitchSpaceMap:         userSpaceNotifyConfigMap,
	}

	for _, opt := range opts {
		opt(ctx)
	}

	return ctx
}
