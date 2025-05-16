package v1_4_13_0

import (
	"context"
	"fmt"
	"go-cs/internal/domain/space/repo"
	space_member_repo "go-cs/internal/domain/space_member/repo"
	"go-cs/internal/domain/space_view"
	view_repo "go-cs/internal/domain/space_view/repo"
	view_service "go-cs/internal/domain/space_view/service"
	version_repo "go-cs/internal/domain/space_work_version/repo"
	status_repo "go-cs/internal/domain/work_item_status/repo"
	"gorm.io/gorm"
	"slices"
)

type DML struct {
	db          *gorm.DB
	spaceRepo   repo.SpaceRepo
	versionRepo version_repo.SpaceWorkVersionRepo
	memberRepo  space_member_repo.SpaceMemberRepo
	viewRepo    view_repo.SpaceViewRepo
	statusRepo  status_repo.WorkItemStatusRepo
	viewService *view_service.SpaceViewService
}

func NewDML(
	db *gorm.DB,
	spaceRepo repo.SpaceRepo,
	versionRepo version_repo.SpaceWorkVersionRepo,
	memberRepo space_member_repo.SpaceMemberRepo,
	viewRepo view_repo.SpaceViewRepo,
	statusRepo status_repo.WorkItemStatusRepo,
	viewService *view_service.SpaceViewService,
) *DML {
	return &DML{
		db:          db,
		spaceRepo:   spaceRepo,
		versionRepo: versionRepo,
		memberRepo:  memberRepo,
		viewRepo:    viewRepo,
		statusRepo:  statusRepo,
		viewService: viewService,
	}
}

func (d *DML) HandleData() error {
	spaceIds, err := d.spaceRepo.GetAllSpaceIds()
	if err != nil {
		return err
	}

	slices.Sort(spaceIds)

	view := make([]*space_view.SpaceGlobalView, 0)
	for _, spaceId := range spaceIds {
		fmt.Printf("Handle %v", spaceId)
		// 1.空间
		view = append(view, d.HandleSpace(spaceId)...)
	}

	err = d.viewRepo.CreateGlobalViews(context.Background(), view)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

const limit = 5000
const insertBatch = 3000

func (d *DML) HandleSpace(spaceId int64) []*space_view.SpaceGlobalView {
	info, err := d.statusRepo.GetWorkItemStatusInfo(context.Background(), spaceId)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	view, err := d.viewService.InitSpacePublicView(context.Background(), spaceId, info.Items.GetProcessingStatus())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return view
}
