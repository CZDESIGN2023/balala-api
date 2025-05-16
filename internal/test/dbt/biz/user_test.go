package biz

import (
	"context"
	v1 "go-cs/api/user/v1"
	"go-cs/internal/test/dbt"
	"go-cs/pkg/pprint"
	"testing"
)

func TestSetSpaceOrder(t *testing.T) {
	err := dbt.UC.UserUsecase.SetSpaceOrder(context.Background(), 42, &v1.SetSpaceOrderRequest{
		FromIdx: 0,
		ToIdx:   1,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestAllSpaceInfo(t *testing.T) {
	data, err := dbt.UC.UserUsecase.AllSpaceInfo(context.Background(), 42, 42)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(data)
}
