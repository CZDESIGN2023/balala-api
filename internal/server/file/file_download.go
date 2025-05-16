package file

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-cs/internal/utils/errs"
	"net/url"
	"path/filepath"
)

func (s *Server) getDownloadToken(ctx context.Context, uid int64, args GetDownloadTokenArgs) (any, error) {
	switch args.Scene {
	default:
		// 判断是不是文件拥有者
		fileInfo, err := s.fileRepo.GetFileInfo(ctx, args.Id)
		if err != nil {
			return nil, err
		}

		if fileInfo.Owner != uid {
			return nil, errors.New("no permission")
		}
	case "space_file":
		//判断是不是空间成员 不是的话不允许下载
		member, err := s.spaceMemberRepo.GetSpaceMember(ctx, args.SpaceId, uid)
		if member == nil || err != nil {
			return "", err
		}
	}

	token, err := s.spaceFileInfoUc.GenerateDownloadFileToken(ctx, args.Id)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Server) download(ctx context.Context, args DownloadArgs) (any, error) {
	switch args.Scene {
	default:
		return nil, errors.New("unknown scene")
	case "space_file":
		fileUrl, fileName, err := s.spaceFileInfoUc.GetSpaceWorkItemFileUriByToken(ctx, args.DownloadToken)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}

		ginCtx := ctx.(*gin.Context)
		ginCtx.Header("Content-Type", "application/octet-stream")
		ginCtx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(fileName))) // 让游览器直接下载文件，而不是预览
		ginCtx.File(filepath.Join(s.confFile.LocalPath, fileUrl.Path))
		return nil, nil
	}
}
