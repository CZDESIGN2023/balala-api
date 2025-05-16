package biz

import (
	"context"
	v1 "go-cs/api/space_work_item/v1"
	"go-cs/internal/test/dbt"
	"testing"
)

func TestAddComment(t *testing.T) {
	res, err := dbt.UC.SpaceWorkItemUsecase.AddComment(context.Background(), 21, &v1.AddWorkItemCommentRequest{
		WorkItemId: 618,
		Content:    "2",
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestCommentList(t *testing.T) {
	res, err := dbt.UC.SpaceWorkItemUsecase.CommentList(context.Background(), 42, &v1.WorkItemCommentListRequest{
		WorkItemId: 618,
		Size:       1,
		Order:      "DESC",
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
