package space_work_item_comment

import (
	"encoding/json"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"slices"
	"strconv"
)

type ReferUserIds []int64

func (r ReferUserIds) ToJsonString() string {
	return utils.ToJSON(r)
}

func (r ReferUserIds) FormJsonString(jsonStr string) ReferUserIds {
	json.Unmarshal([]byte(jsonStr), &r)
	return r
}

type Emojis []EmojiInfo

type EmojiInfo struct {
	Id      string   `json:"id"`
	UserIds []string `json:"user_ids"`
}

func (e Emojis) ToJson() string {
	return utils.ToJSON(e)
}

func (e Emojis) FromJson(data string) Emojis {
	v := Emojis{}
	json.Unmarshal([]byte(data), &v)
	return v
}

type SpaceWorkItemComments []*SpaceWorkItemComment

type SpaceWorkItemComment struct {
	shared.AggregateRoot

	Id             int64
	UserId         int64
	WorkItemId     int64
	Content        string
	ReferUserIds   ReferUserIds
	Emojis         Emojis
	ReplyCommentId int64
	CreatedAt      int64
	UpdatedAt      int64
	DeletedAt      int64
}

func (s *SpaceWorkItemComment) UpdateContent(content string, ids ReferUserIds, oper shared.Oper) {
	oldContent := s.Content

	s.Content = content
	s.ReferUserIds = ids

	s.AddDiff(Diff_ReferUserIds, Diff_Content)

	s.AddMessage(oper, &domain_message.UpdateComment{
		WorkItemId:     s.WorkItemId,
		CommentId:      s.Id,
		UserId:         s.UserId,
		OldContent:     oldContent,
		NewContent:     content,
		ReferUserIds:   ids,
		ReplyCommentId: s.ReplyCommentId,
	})
}

func (s *SpaceWorkItemComment) AddEmoji(emojiId string, userId int64, oper shared.Oper) {
	if emojiId == "" || userId <= 0 {
		return
	}

	uidStr := strconv.FormatInt(userId, 10)

	idx := slices.IndexFunc(s.Emojis, func(info EmojiInfo) bool {
		return info.Id == emojiId
	})

	if idx < 0 {
		s.Emojis = append(s.Emojis, EmojiInfo{Id: emojiId, UserIds: []string{uidStr}})
	} else {
		if slices.Contains(s.Emojis[idx].UserIds, uidStr) {
			return
		}
		s.Emojis[idx].UserIds = append(s.Emojis[idx].UserIds, uidStr)
	}

	s.AddDiff(Diff_Emojis)

	s.AddMessage(oper, &domain_message.AddCommentEmoji{
		WorkItemId:     s.WorkItemId,
		CommentId:      s.Id,
		UserId:         s.UserId,
		Content:        s.Content,
		ReferUserIds:   s.ReferUserIds,
		ReplyCommentId: s.ReplyCommentId,
		Emoji:          emojiId,
	})
}
func (s *SpaceWorkItemComment) RemoveEmoji(emojiId string, userId int64, oper shared.Oper) {
	uidStr := strconv.FormatInt(userId, 10)

	idx := slices.IndexFunc(s.Emojis, func(info EmojiInfo) bool {
		return info.Id == emojiId
	})

	if idx < 0 {
		return
	}

	if !slices.Contains(s.Emojis[idx].UserIds, uidStr) {
		return
	}

	s.Emojis[idx].UserIds = stream.Remove(s.Emojis[idx].UserIds, uidStr)

	if len(s.Emojis[idx].UserIds) == 0 {
		s.Emojis = slices.Delete(s.Emojis, idx, idx+1)
	}

	s.AddDiff(Diff_Emojis)

	s.AddMessage(oper, &domain_message.RemoveCommentEmoji{
		WorkItemId:     s.WorkItemId,
		CommentId:      s.Id,
		UserId:         s.UserId,
		Content:        s.Content,
		ReferUserIds:   s.ReferUserIds,
		ReplyCommentId: s.ReplyCommentId,
		Emoji:          emojiId,
	})
}

func (s *SpaceWorkItemComment) OnDelete(oper shared.Oper) {
	s.AddMessage(oper, &domain_message.DeleteComment{
		WorkItemId:     s.WorkItemId,
		CommentId:      s.Id,
		UserId:         s.UserId,
		Content:        s.Content,
		ReferUserIds:   s.ReferUserIds,
		ReplyCommentId: s.ReplyCommentId,
	})
}
