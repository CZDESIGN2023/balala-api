package data

import (
	"context"
	"testing"
	"time"
)

func TestSpaceWorkItemCommentTagRepo_CreateComment(t *testing.T) {
	// comment, err := SpaceWorkItemCommentRepo.CreateComment(context.Background(), &db.SpaceWorkItemComment{
	// 	Id:         0,
	// 	UserId:     42,
	// 	WorkItemId: 359,
	// 	Content:    "2",
	// })
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(comment)
}

func TestSpaceWorkItemCommentTagRepo_UpdateComment(t *testing.T) {
	// comment, err := SpaceWorkItemCommentRepo.UpdateComment(context.Background(), &db.SpaceWorkItemComment{
	// 	Id:      4,
	// 	Content: "2",
	// })
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(comment)
}

func TestSpaceWorkItemCommentRepo_GetComment(t *testing.T) {
	comment, err := SpaceWorkItemCommentRepo.GetComment(context.Background(), 3)
	if err != nil {
		t.Error(err)
	}
	t.Log(comment)
}

func TestSpaceWorkItemCommentRepo_GetCommentByIds(t *testing.T) {
	comment, err := SpaceWorkItemCommentRepo.GetCommentByIds(context.Background(), []int64{3, 4})
	if err != nil {
		t.Error(err)
	}
	t.Log(comment)
}

func TestSpaceWorkItemCommentRepo_GetCommentByWorkItemId(t *testing.T) {
	comment, err := SpaceWorkItemCommentRepo.GetCommentByWorkItemId(context.Background(), 361)
	if err != nil {
		t.Error(err)
	}
	t.Log(comment)
}

func TestSpaceWorkItemCommentRepo_CountCommentByWorkItemId(t *testing.T) {
	comment, err := SpaceWorkItemCommentRepo.CountCommentByWorkItemId(context.Background(), 361)
	if err != nil {
		t.Error(err)
	}
	t.Log(comment)
}

func TestSpaceWorkItemCommentRepo_CountWorkItemCommentNumByTime(t *testing.T) {
	comment, err := SpaceWorkItemCommentRepo.CountWorkItemCommentNumByTime(context.Background(), 359, time.Unix(1700535614, 0))
	if err != nil {
		t.Error(err)
	}
	t.Log(comment)
}

func TestSpaceWorkItemCommentRepo_CountCommentByWorkItemIds(t *testing.T) {
	comment, err := SpaceWorkItemCommentRepo.CountCommentByWorkItemIds(context.Background(), []int64{359, 361})
	if err != nil {
		t.Error(err)
	}
	t.Log(comment)
}

func TestSpaceWorkItemCommentRepo_DelComment(t *testing.T) {
	res, err := SpaceWorkItemCommentRepo.DelCommentByWorkItemId(context.Background(), 359)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestSpaceWorkItemCommentRepo_IncrUnreadNumForUser(t *testing.T) {
	err := SpaceWorkItemCommentRepo.IncrUnreadNumForUser(context.Background(), 358, []int64{42})
	if err != nil {
		t.Error(err)
	}
}

func TestSpaceWorkItemCommentRepo_UserUnreadNumMapByWorkItemIds(t *testing.T) {
	res, err := SpaceWorkItemCommentRepo.UserUnreadNumMapByWorkItemIds(context.Background(), 42, []int64{})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
