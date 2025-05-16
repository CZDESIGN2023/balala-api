package v1_3_7_2

import (
	"context"
	"encoding/json"
	"fmt"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/conf"
	dim_model "go-cs/internal/dwh/model/dim"
	dwd_model "go-cs/internal/dwh/model/dwd"
	"go-cs/internal/utils/date"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DML struct {
	log  *klog.Helper
	conf *conf.Dwh

	_dwhDb        *gorm.DB
	dwhDbInitOnce sync.Once

	_exDb        *gorm.DB
	exDbInitOnce sync.Once
}

func NewDML(
	logger log.Logger,
	conf *conf.Dwh,

) *DML {
	kLog := log.NewHelper(logger)
	return &DML{
		log:  kLog,
		conf: conf,
	}
}

func (u *DML) dwhDb() *gorm.DB {
	u.dwhDbInitOnce.Do(func() {
		db, err := gorm.Open(mysql.Open(u.conf.Database.Dsn), &gorm.Config{})
		if err != nil {
			u.log.Error(err)
			panic(err)
		}

		u._dwhDb = db
	})
	return u._dwhDb
}

func (u *DML) exDb() *gorm.DB {
	u.exDbInitOnce.Do(func() {
		db, err := gorm.Open(mysql.Open(u.conf.ExternalDataSource.Database.Dsn), &gorm.Config{})
		if err != nil {
			u.log.Error(err)
			panic(err)
		}

		u._exDb = db
	})
	return u._exDb
}

func (u *DML) Start(ctx context.Context) (err error) {

	fmt.Println("---- 同步[用户]数据至 数仓, 开始处理...")
	err = u.syncUser()
	if err != nil {
		return err
	}
	fmt.Println("---- 同步完成")

	fmt.Println("---- 同步[空间]数据至 数仓, 开始处理...")
	err = u.syncSpace()
	if err != nil {
		return err
	}
	fmt.Println("---- 同步完成")

	fmt.Println("---- 同步[空间成员]数据至 数仓, 开始处理...")
	err = u.syncSpaceMember()
	if err != nil {
		return err
	}
	fmt.Println("---- 同步完成")

	fmt.Println("---- 同步[空间模块]数据至 数仓, 开始处理...")
	err = u.syncSpaceWorkObject()
	if err != nil {
		return err
	}
	fmt.Println("---- 同步完成")

	fmt.Println("---- 同步[工作项状态]数据至 数仓, 开始处理...")
	err = u.syncSpaceWorkStatus()
	if err != nil {
		return err
	}
	fmt.Println("---- 同步完成")

	fmt.Println("---- 同步[工作项版本]数据至 数仓, 开始处理...")
	err = u.syncSpaceWorkVersion()
	if err != nil {
		return err
	}
	fmt.Println("---- 同步完成")

	fmt.Println("---- 同步[工作项]数据至 数仓, 开始处理...")
	err = u.syncSpaceWorkItem()
	if err != nil {
		return err
	}
	fmt.Println("---- 同步完成")

	fmt.Println("---- 同步[工作项-节点]数据至 数仓, 开始处理...")
	err = u.syncSpaceWorkItemFlowNode()
	if err != nil {
		return err
	}
	fmt.Println("---- 同步完成")

	return nil
}

func (u *DML) syncUser() error {

	var rows []*db.User

	err := u.exDb().Table("user").Find(&rows).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var dimRows []*dim_model.DimUser
	for _, row := range rows {
		dimRow := &dim_model.DimUser{}
		dimRow.UserId = row.Id
		dimRow.UserName = row.UserName
		dimRow.UserNickName = row.UserNickname
		dimRow.UserPinyin = row.UserPinyin
		dimRow.GmtCreate = cast.ToTime(row.CreatedAt)
		dimRow.GmtModified = cast.ToTime(row.UpdatedAt)
		dimRow.StartDate = cast.ToTime(row.CreatedAt)
		dimRow.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
		dimRows = append(dimRows, dimRow)
	}

	txErr := u.dwhDb().Transaction(func(tx *gorm.DB) error {
		for _, v := range dimRows {
			err = tx.Table("dim_user").Where("user_id = ?", v.UserId).Save(v).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (u *DML) syncSpace() error {

	var rows []*db.Space

	err := u.exDb().Table("space").Find(&rows).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var dimRows []*dim_model.DimSpace
	for _, row := range rows {
		dimRow := &dim_model.DimSpace{}
		dimRow.SpaceId = row.Id
		dimRow.SpaceName = row.SpaceName
		dimRow.GmtCreate = cast.ToTime(row.CreatedAt)
		dimRow.GmtModified = cast.ToTime(row.UpdatedAt)
		dimRow.StartDate = cast.ToTime(row.CreatedAt)
		dimRow.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
		dimRows = append(dimRows, dimRow)
	}

	txErr := u.dwhDb().Transaction(func(tx *gorm.DB) error {
		for _, v := range dimRows {
			err = tx.Table("dim_space").Where("space_id = ?", v.SpaceId).Save(v).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (u *DML) syncSpaceMember() error {

	var rows []*db.SpaceMember

	err := u.exDb().Table("space_member").Find(&rows).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var dimRows []*dwd_model.DwdMember
	for _, row := range rows {
		dimRow := &dwd_model.DwdMember{}
		dimRow.SpaceId = row.SpaceId
		dimRow.MemberId = row.Id
		dimRow.UserId = row.UserId
		dimRow.RoleId = row.RoleId
		dimRow.GmtCreate = cast.ToTime(row.CreatedAt)
		dimRow.GmtModified = cast.ToTime(row.UpdatedAt)
		dimRow.StartDate = cast.ToTime(row.CreatedAt)
		dimRow.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
		dimRows = append(dimRows, dimRow)
	}

	txErr := u.dwhDb().Transaction(func(tx *gorm.DB) error {
		for _, v := range dimRows {
			err = tx.Table("dwd_member").Where("member_id = ? AND end_date = ?", v.MemberId, v.EndDate).Save(v).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (u *DML) syncSpaceWorkObject() error {

	var rows []*db.SpaceWorkObject

	err := u.exDb().Table("space_work_object").Find(&rows).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var dimRows []*dim_model.DimObject
	for _, row := range rows {
		dimRow := &dim_model.DimObject{}
		dimRow.SpaceId = row.SpaceId
		dimRow.ObjectId = row.Id
		dimRow.ObjectName = row.WorkObjectName
		dimRow.GmtCreate = cast.ToTime(row.CreatedAt)
		dimRow.GmtModified = cast.ToTime(row.UpdatedAt)
		dimRow.StartDate = cast.ToTime(row.CreatedAt)
		dimRow.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
		dimRows = append(dimRows, dimRow)
	}

	txErr := u.dwhDb().Transaction(func(tx *gorm.DB) error {
		for _, v := range dimRows {
			err = tx.Table("dim_object").Where("object_id = ?", v.ObjectId).Save(v).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (u *DML) syncSpaceWorkStatus() error {

	var rows []*db.WorkItemStatus

	err := u.exDb().Table("work_item_status").Find(&rows).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var dimRows []*dim_model.DimWitemStatus
	for _, row := range rows {
		dimRow := &dim_model.DimWitemStatus{}
		dimRow.SpaceId = row.SpaceId
		dimRow.StatusId = row.Id
		dimRow.StatusName = row.Name
		dimRow.StatusKey = row.Key
		dimRow.StatusVal = row.Val
		dimRow.StatusType = row.StatusType
		dimRow.GmtCreate = cast.ToTime(row.CreatedAt)
		dimRow.GmtModified = cast.ToTime(row.UpdatedAt)
		dimRow.StartDate = cast.ToTime(row.CreatedAt)
		dimRow.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
		dimRows = append(dimRows, dimRow)
	}

	txErr := u.dwhDb().Transaction(func(tx *gorm.DB) error {
		for _, v := range dimRows {
			err = tx.Table("dim_witem_status").Where("status_id = ?", v.StatusId).Save(v).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (u *DML) syncSpaceWorkVersion() error {

	var rows []*db.SpaceWorkVersion

	err := u.exDb().Table("space_work_version").Find(&rows).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var dimRows []*dim_model.DimVersion
	for _, row := range rows {
		dimRow := &dim_model.DimVersion{}
		dimRow.SpaceId = row.SpaceId
		dimRow.VersionId = row.Id
		dimRow.VersionName = row.VersionName
		dimRow.GmtCreate = cast.ToTime(row.CreatedAt)
		dimRow.GmtModified = cast.ToTime(row.UpdatedAt)
		dimRow.StartDate = cast.ToTime(row.CreatedAt)
		dimRow.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
		dimRows = append(dimRows, dimRow)
	}

	txErr := u.dwhDb().Transaction(func(tx *gorm.DB) error {
		for _, v := range dimRows {
			err = tx.Table("dim_version").Where("version_id = ?", v.VersionId).Save(v).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

type spaceWorkItemDoc struct {
	PlanStartAt    int64    `gorm:"column:plan_start_at" bson:"plan_start_at" json:"plan_start_at"`
	PlanCompleteAt int64    `gorm:"column:plan_complete_at" bson:"plan_complete_at" json:"plan_complete_at"`
	ProcessRate    int32    `gorm:"column:process_rate" bson:"process_rate" json:"process_rate"`
	Remark         string   `gorm:"column:remark" bson:"remark" json:"remark"`
	Describe       string   `gorm:"column:describe" bson:"describe" json:"describe"`
	Priority       string   `gorm:"column:priority" bson:"priority" json:"priority"`
	Tags           []string `gorm:"column:tags" bson:"tags" json:"tags,omitempty"`
	Directors      []string `gorm:"column:directors" bson:"directors" json:"directors,omitempty"`
	Followers      []string `gorm:"column:followers" bson:"followers" json:"followers,omitempty"`
	Participators  []string `gorm:"column:participators" bson:"participators" json:"participators,omitempty"`
	// NodeDirectors  []string `gorm:"column:node_directors" bson:"node_directors" json:"node_directors,omitempty"`
}

func (u *DML) syncSpaceWorkItem() error {

	var rows []*db.SpaceWorkItemV2

	err := u.exDb().Table("space_work_item_v2").Find(&rows).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var dimRows []*dwd_model.DwdWitem
	for _, row := range rows {

		//解析doc json
		var witemDoc *spaceWorkItemDoc
		err := json.Unmarshal([]byte(row.Doc), &witemDoc)
		if err != nil {
			fmt.Println(err)
			continue
		}

		directors, _ := json.Marshal(witemDoc.Directors)
		participators, _ := json.Marshal(witemDoc.Participators)

		dwdWitem := &dwd_model.DwdWitem{}
		dwdWitem.SpaceId = row.SpaceId
		dwdWitem.WorkItemId = row.Id
		dwdWitem.UserId = row.UserId
		dwdWitem.StatusId = row.WorkItemStatusId
		dwdWitem.ObjectId = row.WorkObjectId
		dwdWitem.VersionId = row.VersionId

		dwdWitem.PlanStartAt = witemDoc.PlanStartAt
		dwdWitem.PlanCompleteAt = witemDoc.PlanCompleteAt
		dwdWitem.Priority = witemDoc.Priority
		dwdWitem.Directors = string(directors)
		dwdWitem.Participators = string(participators)

		dwdWitem.GmtCreate = cast.ToTime(row.CreatedAt)
		dwdWitem.GmtModified = cast.ToTime(row.UpdatedAt)
		dwdWitem.StartDate = cast.ToTime(row.CreatedAt)
		dwdWitem.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
		if row.DeletedAt > 0 {
			dwdWitem.EndDate = cast.ToTime(row.DeletedAt)
		}
		dimRows = append(dimRows, dwdWitem)
	}

	txErr := u.dwhDb().Transaction(func(tx *gorm.DB) error {
		for _, v := range dimRows {
			err = tx.Table("dwd_witem").Where("work_item_id = ? AND end_date = ?", v.WorkItemId, v.EndDate).Save(v).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (u *DML) syncSpaceWorkItemFlowNode() error {

	var rows []*db.SpaceWorkItemFlowV2

	err := u.exDb().Table("space_work_item_flow_v2").Find(&rows).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var dimRows []*dwd_model.DwdWitemFlowNode
	for _, row := range rows {
		dwdWitemFlowNode := &dwd_model.DwdWitemFlowNode{}
		dwdWitemFlowNode.SpaceId = row.SpaceId
		dwdWitemFlowNode.WorkItemId = row.WorkItemId
		dwdWitemFlowNode.NodeId = row.Id
		dwdWitemFlowNode.NodeCode = row.FlowNodeCode
		dwdWitemFlowNode.NodeStatus = row.FlowNodeStatus

		dwdWitemFlowNode.PlanStartAt = row.PlanStartAt
		dwdWitemFlowNode.PlanCompleteAt = row.PlanCompleteAt
		dwdWitemFlowNode.Directors = string(row.Directors)
		dwdWitemFlowNode.GmtCreate = cast.ToTime(row.CreatedAt)
		dwdWitemFlowNode.GmtModified = cast.ToTime(row.UpdatedAt)
		dwdWitemFlowNode.StartDate = cast.ToTime(row.CreatedAt)
		dwdWitemFlowNode.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
		if row.DeletedAt > 0 {
			dwdWitemFlowNode.EndDate = cast.ToTime(row.DeletedAt)
		}
		dimRows = append(dimRows, dwdWitemFlowNode)
	}

	txErr := u.dwhDb().Transaction(func(tx *gorm.DB) error {
		for _, v := range dimRows {
			err = tx.Table("dwd_witem_flow_node").Where("work_item_id = ? AND end_date = ?", v.WorkItemId, v.EndDate).Save(v).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}
