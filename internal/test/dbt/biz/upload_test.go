package biz

import (
	"context"
	"go-cs/internal/test/dbt"
	"testing"
)

func TestDownLoadAvatar(t *testing.T) {
	ctx := context.Background()

	space, err := dbt.UC.UploadUsecase.DownloadAvatarImgToLocal(ctx, 7, "https://www.twle.cn/static/i/img1.jpg", dbt.C.FileConfig.LocalPath)
	if err != nil {
		t.Error(err)
	}

	t.Log(space)
}
