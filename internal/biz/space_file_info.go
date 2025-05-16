package biz

import (
	"context"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/rand"
	"net/url"
	"path"
	"time"

	file_repo "go-cs/internal/domain/file_info/repo"
	space_file_repo "go-cs/internal/domain/space_file_info/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	user_repo "go-cs/internal/domain/user/repo"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type SpaceFileInfoUsecase struct {
	log *log.Helper
	tm  trans.Transaction

	spaceMemberRepo   member_repo.SpaceMemberRepo
	userRepo          user_repo.UserRepo
	spaceFileInfoRepo space_file_repo.SpaceFileInfoRepo
	fileInfoRepo      file_repo.FileInfoRepo
}

func NewSpaceFileInfoUsecase(
	logger log.Logger,
	tm trans.Transaction,

	spaceMemberRepo member_repo.SpaceMemberRepo,
	userRepo user_repo.UserRepo,
	spaceFileInfoRepo space_file_repo.SpaceFileInfoRepo,
	fileInfoRepo file_repo.FileInfoRepo,

) *SpaceFileInfoUsecase {

	moduleName := "SpaceFileInfoUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)
	return &SpaceFileInfoUsecase{
		log: hlog,
		tm:  tm,

		spaceMemberRepo:   spaceMemberRepo,
		userRepo:          userRepo,
		spaceFileInfoRepo: spaceFileInfoRepo,
		fileInfoRepo:      fileInfoRepo,
	}
}

func (s *SpaceFileInfoUsecase) GenerateDownloadFileToken(ctx context.Context, fileInfoId int64) (string, error) {
	pwd := rand.Letters(10)
	err := s.spaceFileInfoRepo.SaveFileDownToken(ctx, pwd, fileInfoId, int64(60*time.Second))
	if err != nil {
		return "", err
	}

	return pwd, nil

}

func (s *SpaceFileInfoUsecase) GetSpaceWorkItemFileUriByToken(ctx context.Context, token string) (*url.URL, string, error) {

	fileInfoId, err := s.spaceFileInfoRepo.GetFileDownToken(ctx, token)
	if err != nil {
		return nil, "", err
	}

	fileInfo, err := s.fileInfoRepo.GetFileInfo(ctx, cast.ToInt64(fileInfoId))
	if err != nil {
		return nil, "", err
	}

	fileUrl, err := url.ParseRequestURI(path.Join("/", fileInfo.Uri))
	if err != nil {
		return nil, "", err
	}

	return fileUrl, fileInfo.Name, nil
}
