package data

import (
	"go-cs/internal/data/test"
	"go-cs/internal/domain/space_member/repo"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

func init() {
	test.Init("../../configs/config.yaml", "../../.env.local", true)

	logger := log.DefaultLogger
	data := &Data{
		db:   test.GetDB(),
		dbRo: &ReadonlyGormDB{test.GetDB()},
		rdb:  test.GetRedis(),
		es:   test.GetEs(),
		conf: test.GetConf().Data,
	}
	gdb = data.db

	UserRepo = NewUserRepo(data, logger).(*userRepo)
	SpaceRepo = NewSpaceRepo(data, logger).(*spaceRepo)
	SpaceMemberRepo = NewSpaceMemberRepo(data, logger).(*spaceMemberRepo)
	SearchRepo = NewSearchRepo(data, test.GetConf(), logger).(*searchRepo)
	StatusRepo = NewWorkItemStatusRepo(data, logger).(*workItemStatusRepo)
	SpaceWorkItemRepo = NewSpaceWorkItemRepo(data, logger).(*spaceWorkItemRepo)
	SpaceFileInfoRepo = NewSpaceFileInfoRepo(data, logger).(*spaceFileInfoRepo)
	fileInfoRepo = NewFileInfoRepo(data, logger).(*FileInfoRepo)
	StaticsRepo = NewStaticsRepo(data, StatusRepo, logger).(*staticsRepo)
	StaticsEsRepo = NewStaticsEsRepo(data, StatusRepo, logger).(*staticsEsRepo)
	SpaceTagRepo = NewSpaceTagRepo(data, logger).(*spaceTagRepo)
	SpaceWorkItemCommentRepo = NewSpaceWorkItemCommentRepo(data, logger).(*spaceWorkItemCommentRepo)
	SpaceWorkObjectRepo = NewSpaceWorkObjectRepo(data, logger).(*spaceWorkObjectRepo)
	NotifyRepo = NewNotifyRepo(data, logger).(*notifyRepo)
	UserLoginLogRepo = NewUserLoginLogRepo(data, logger).(*userLoginLogRepo)
	CondfigRepo = NewConfigRepo(data, logger).(*configRepo)
	SpaceMemberRepo = NewSpaceMemberRepo(data, logger)
	WorkFlowRepo = NewWorkFlowRepo(data, logger).(*workFlowRepo)
}

var gdb *gorm.DB

var (
	UserRepo         *userRepo
	UserLoginLogRepo *userLoginLogRepo
	SpaceRepo        *spaceRepo
	SpaceMemberRepo  repo.SpaceMemberRepo
	SearchRepo       *searchRepo
	StatusRepo       *workItemStatusRepo
	// SpaceWorkFlowRepo        *spaceWorkItemFlowRepo
	SpaceWorkItemRepo        *spaceWorkItemRepo
	SpaceFileInfoRepo        *spaceFileInfoRepo
	fileInfoRepo             *FileInfoRepo
	StaticsRepo              *staticsRepo
	StaticsEsRepo            *staticsEsRepo
	SpaceTagRepo             *spaceTagRepo
	SpaceWorkItemCommentRepo *spaceWorkItemCommentRepo
	SpaceWorkObjectRepo      *spaceWorkObjectRepo
	NotifyRepo               *notifyRepo
	CondfigRepo              *configRepo
	WorkFlowRepo             *workFlowRepo
)
