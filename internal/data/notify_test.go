package data

import (
	"context"
	"go-cs/pkg/pprint"
	"testing"
)

func TestNotifyRepo_GetDelOfflineNotify(t *testing.T) {
	notify, err := NotifyRepo.GetDelOfflineNotify(context.Background(), 42)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(notify)
}
