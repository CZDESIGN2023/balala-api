package dbt

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/subosito/gotenv"
	"go-cs/internal/biz"
	"go-cs/internal/conf"
	"go-cs/internal/data"
	"go-cs/internal/domain/space_view/repo"
	space_view_service "go-cs/internal/domain/space_view/service"
	"go-cs/internal/service"

	biz_id_domain "go-cs/internal/pkg/biz_id"

	file_repo "go-cs/internal/domain/file_info/repo"
	file_service "go-cs/internal/domain/file_info/service"

	space_file_repo "go-cs/internal/domain/space_file_info/repo"
	space_file_service "go-cs/internal/domain/space_file_info/service"

	tag_repo "go-cs/internal/domain/space_tag/repo"
	tag_service "go-cs/internal/domain/space_tag/service"

	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	workObj_service "go-cs/internal/domain/space_work_object/service"

	space_work_version_repo "go-cs/internal/domain/space_work_version/repo"
	space_work_version_service "go-cs/internal/domain/space_work_version/service"

	space_repo "go-cs/internal/domain/space/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	user_repo "go-cs/internal/domain/user/repo"
	wf_repo "go-cs/internal/domain/work_flow/repo"
	role_repo "go-cs/internal/domain/work_item_role/repo"
	wiz_repo "go-cs/internal/domain/work_item_status/repo"
	witype_repo "go-cs/internal/domain/work_item_type/repo"

	statics_repo "go-cs/internal/domain/statics/repo"

	wf_service "go-cs/internal/domain/work_flow/service"
	role_service "go-cs/internal/domain/work_item_role/service"
	wiz_service "go-cs/internal/domain/work_item_status/service"

	search_repo "go-cs/internal/domain/search/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"
	witem_service "go-cs/internal/domain/work_item/service"

	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type All struct {
	Data          *data.Data
	Service       *Service
	Usecase       *Usecase
	Repo          *Repo
	DomainService *DomainService
}

type Service struct {
	LoginService           *service.LoginService
	UserService            *service.UserService
	ConfigService          *service.ConfigService
	SpaceService           *service.SpaceService
	SpaceMemberService     *service.SpaceMemberService
	SpaceTagService        *service.SpaceTagService
	SpaceWorkItemService   *service.SpaceWorkItemService
	SpaceWorkObjectService *service.SpaceWorkObjectService
	SearchService          *service.SearchService
	WorkbenchService       *service.WorkbenchService
}

type DomainService struct {
	BusinessIdService       *biz_id_domain.BusinessIdService
	WorkFlowService         *wf_service.WorkFlowService
	WorkItemRoleService     *role_service.WorkItemRoleService
	WorkItemStatusService   *wiz_service.WorkItemStatusService
	FileInfoService         *file_service.FileInfoService
	SpaceFileInfoService    *space_file_service.SpaceFileInfoService
	SpaceWorkObjectService  *workObj_service.SpaceWorkObjectService
	SpaceTagService         *tag_service.SpaceTagService
	SpaceWorkVersionService *space_work_version_service.SpaceWorkVersionService
	WorkItemService         *witem_service.WorkItemService
	SpaceViewService        *space_view_service.SpaceViewService
}

type Usecase struct {
	LoginUsecase           *biz.LoginUsecase
	UserUsecase            *biz.UserUsecase
	UploadUsecase          *biz.UploadUsecase
	ConfigUsecase          *biz.ConfigUsecase
	SpaceUsecase           *biz.SpaceUsecase
	SpaceMemberUsecase     *biz.SpaceMemberUsecase
	SpaceTagUsecase        *biz.SpaceTagUsecase
	SpaceWorkObjectUsecase *biz.SpaceWorkObjectUsecase
	SearchUsecase          *biz.SearchUsecase
	SpaceWorkFlowUsecase   *biz.SpaceWorkItemFlowUsecase
	SpaceWorkItemUsecase   *biz.SpaceWorkItemUsecase
	WorkFlowUsecase        *biz.WorkFlowUsecase
	StaticsUsecase         *biz.StaticsUsecase
	SpaceViewUsecase       *biz.SpaceViewUsecase
	RptUsecase             *biz.RptUsecase
}

type Repo struct {
	LoginRepo            biz.LoginRepo
	UserRepo             user_repo.UserRepo
	FileInfoRepo         file_repo.FileInfoRepo
	SpaceFileInfoRepo    space_file_repo.SpaceFileInfoRepo
	SpaceWorkItemRepo    witem_repo.WorkItemRepo
	SpaceMemberRepo      member_repo.SpaceMemberRepo
	ConfigRepo           biz.ConfigRepo
	SpaceRepo            space_repo.SpaceRepo
	SpaceWorkObjectRepo  workObj_repo.SpaceWorkObjectRepo
	SpaceWorkVersionRepo space_work_version_repo.SpaceWorkVersionRepo

	SpaceTagRepo tag_repo.SpaceTagRepo
	StaticsRepo  statics_repo.StaticsRepo
	SearchRepo   search_repo.SearchRepo

	WorkStatusRepo   wiz_repo.WorkItemStatusRepo
	WorkFlowRepo     wf_repo.WorkFlowRepo
	WorkItemTypeRepo witype_repo.WorkItemTypeRepo
	WorkItemRoleRepo role_repo.WorkItemRoleRepo
	SpaceViewRepo    repo.SpaceViewRepo
}

func NewAll(data *data.Data, service *Service, usecase *Usecase, repo *Repo, domainService *DomainService) *All {
	return &All{
		Data:          data,
		Service:       service,
		Usecase:       usecase,
		Repo:          repo,
		DomainService: domainService,
	}
}

func NewDomainService(
	BusinessIdService *biz_id_domain.BusinessIdService,
	WorkFlowService *wf_service.WorkFlowService,
	WorkItemRoleService *role_service.WorkItemRoleService,
	WorkItemStatusService *wiz_service.WorkItemStatusService,
	FileInfoService *file_service.FileInfoService,
	SpaceFileInfoService *space_file_service.SpaceFileInfoService,
	SpaceWorkObjectService *workObj_service.SpaceWorkObjectService,
	SpaceTagService *tag_service.SpaceTagService,
	SpaceWorkVersionService *space_work_version_service.SpaceWorkVersionService,
	WorkItemService *witem_service.WorkItemService,
	spaceViewService *space_view_service.SpaceViewService,
) *DomainService {
	return &DomainService{
		BusinessIdService:       BusinessIdService,
		WorkFlowService:         WorkFlowService,
		WorkItemRoleService:     WorkItemRoleService,
		WorkItemStatusService:   WorkItemStatusService,
		FileInfoService:         FileInfoService,
		SpaceFileInfoService:    SpaceFileInfoService,
		SpaceWorkObjectService:  SpaceWorkObjectService,
		SpaceTagService:         SpaceTagService,
		SpaceWorkVersionService: SpaceWorkVersionService,
		WorkItemService:         WorkItemService,
		SpaceViewService:        spaceViewService,
	}
}

func NewService(
	LoginService *service.LoginService,
	UserService *service.UserService,
	ConfigService *service.ConfigService,
	SpaceService *service.SpaceService,
	SpaceMemberService *service.SpaceMemberService,
	SpaceTagService *service.SpaceTagService,
	SpaceWorkItemService *service.SpaceWorkItemService,
	SpaceWorkObjectService *service.SpaceWorkObjectService,
	SearchService *service.SearchService,
	WorkbenchService *service.WorkbenchService,
) *Service {
	return &Service{
		LoginService:           LoginService,
		UserService:            UserService,
		ConfigService:          ConfigService,
		SpaceService:           SpaceService,
		SpaceMemberService:     SpaceMemberService,
		SpaceTagService:        SpaceTagService,
		SpaceWorkItemService:   SpaceWorkItemService,
		SpaceWorkObjectService: SpaceWorkObjectService,
		SearchService:          SearchService,
		WorkbenchService:       WorkbenchService,
	}
}

func NewUsecase(
	LoginUsecase *biz.LoginUsecase,
	UserUsecase *biz.UserUsecase,
	UploadUsecase *biz.UploadUsecase,
	ConfigUsecase *biz.ConfigUsecase,
	SpaceUsecase *biz.SpaceUsecase,
	SpaceMemberUsecase *biz.SpaceMemberUsecase,
	SpaceTagUsecase *biz.SpaceTagUsecase,
	SpaceWorkObjectUsecase *biz.SpaceWorkObjectUsecase,
	SearchUsecase *biz.SearchUsecase,
	SpaceWorkFlowUsecase *biz.SpaceWorkItemFlowUsecase,
	SpaceWorkItemUsecase *biz.SpaceWorkItemUsecase,
	WorkFlowUsecase *biz.WorkFlowUsecase,
	StaticsUsecase *biz.StaticsUsecase,
	SpaceViewUsecase *biz.SpaceViewUsecase,
	RptUsecase *biz.RptUsecase,
) *Usecase {
	return &Usecase{
		LoginUsecase:           LoginUsecase,
		UserUsecase:            UserUsecase,
		UploadUsecase:          UploadUsecase,
		ConfigUsecase:          ConfigUsecase,
		SpaceUsecase:           SpaceUsecase,
		SpaceMemberUsecase:     SpaceMemberUsecase,
		SpaceTagUsecase:        SpaceTagUsecase,
		SpaceWorkObjectUsecase: SpaceWorkObjectUsecase,
		SearchUsecase:          SearchUsecase,
		SpaceWorkFlowUsecase:   SpaceWorkFlowUsecase,
		WorkFlowUsecase:        WorkFlowUsecase,
		SpaceWorkItemUsecase:   SpaceWorkItemUsecase,
		StaticsUsecase:         StaticsUsecase,
		SpaceViewUsecase:       SpaceViewUsecase,
		RptUsecase:             RptUsecase,
	}
}

func NewRepo(
	LoginRepo biz.LoginRepo,
	UserRepo user_repo.UserRepo,
	FileInfoRepo file_repo.FileInfoRepo,
	SpaceFileInfoRepo space_file_repo.SpaceFileInfoRepo,
	SpaceWorkItemRepo witem_repo.WorkItemRepo,
	SpaceMemberRepo member_repo.SpaceMemberRepo,
	ConfigRepo biz.ConfigRepo,
	SpaceRepo space_repo.SpaceRepo,
	SpaceWorkObjectRepo workObj_repo.SpaceWorkObjectRepo,

	SpaceTagRepo tag_repo.SpaceTagRepo,
	StaticsRepo statics_repo.StaticsRepo,
	SearchRepo search_repo.SearchRepo,
	WorkStatusRepo wiz_repo.WorkItemStatusRepo,
	WorkFlowRepo wf_repo.WorkFlowRepo,
	WorkItemTypeRepo witype_repo.WorkItemTypeRepo,
	SpaceWorkVersionRepo space_work_version_repo.SpaceWorkVersionRepo,
	WorkItemRoleRepo role_repo.WorkItemRoleRepo,
	SpaceViewRepo repo.SpaceViewRepo,
) *Repo {
	return &Repo{
		LoginRepo:            LoginRepo,
		UserRepo:             UserRepo,
		FileInfoRepo:         FileInfoRepo,
		SpaceFileInfoRepo:    SpaceFileInfoRepo,
		SpaceWorkItemRepo:    SpaceWorkItemRepo,
		SpaceMemberRepo:      SpaceMemberRepo,
		ConfigRepo:           ConfigRepo,
		SpaceRepo:            SpaceRepo,
		SpaceWorkObjectRepo:  SpaceWorkObjectRepo,
		SpaceTagRepo:         SpaceTagRepo,
		StaticsRepo:          StaticsRepo,
		SearchRepo:           SearchRepo,
		WorkStatusRepo:       WorkStatusRepo,
		WorkFlowRepo:         WorkFlowRepo,
		WorkItemTypeRepo:     WorkItemTypeRepo,
		SpaceWorkVersionRepo: SpaceWorkVersionRepo,
		WorkItemRoleRepo:     WorkItemRoleRepo,
		SpaceViewRepo:        SpaceViewRepo,
	}
}

var (
	all      *All
	S        *Service
	UC       *Usecase
	R        *Repo
	S_Domain *DomainService

	Data *data.Data
	DB   *gorm.DB
	C    conf.Bootstrap
)

func Init(configFile, envFile string, debug ...bool) {
	err := gotenv.Load(envFile)
	if err != nil {
		panic(err)
	}
	c := kconfig.New(
		kconfig.WithSource(
			file.NewSource(configFile),
			env.NewSource(),
		),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var b conf.Bootstrap

	if err := c.Scan(&b); err != nil {
		panic(err)
	}

	if len(debug) > 0 && debug[0] == true {
		b.Data.Database.Debug = true
		b.Data.DatabaseRo.Debug = true
	}

	all, _, err = wireApp(&b, b.Data, b.Jwt, b.FileConfig, b.Dwh, klog.DefaultLogger, zap.NewExample(), nil)
	if err != nil {
		fmt.Println(err)
	}

	UC, R, Data = all.Usecase, all.Repo, all.Data
	S = all.Service
	S_Domain = all.DomainService
}
