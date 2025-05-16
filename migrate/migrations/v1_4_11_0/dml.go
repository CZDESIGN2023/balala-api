package v1_4_11_0

import (
	"context"
	"fmt"
	"go-cs/internal/domain/space/repo"
	member "go-cs/internal/domain/space_member"
	space_member_repo "go-cs/internal/domain/space_member/repo"
	version "go-cs/internal/domain/space_work_version"
	version_repo "go-cs/internal/domain/space_work_version/repo"
	"go-cs/internal/dwh/model/dws"
	"go-cs/pkg/stream"
	"gorm.io/gorm"
	"slices"
)

type DML struct {
	db          *gorm.DB
	spaceRepo   repo.SpaceRepo
	versionRepo version_repo.SpaceWorkVersionRepo
	memberRepo  space_member_repo.SpaceMemberRepo
}

func NewDML(
	db *gorm.DB,
	spaceRepo repo.SpaceRepo,
	versionRepo version_repo.SpaceWorkVersionRepo,
	memberRepo space_member_repo.SpaceMemberRepo,
) *DML {
	return &DML{
		db:          db,
		spaceRepo:   spaceRepo,
		versionRepo: versionRepo,
		memberRepo:  memberRepo,
	}
}

func (d *DML) HandleData() error {
	spaceIds, err := d.spaceRepo.GetAllSpaceIds()
	if err != nil {
		return err
	}

	slices.Sort(spaceIds)

	for _, spaceId := range spaceIds {
		fmt.Printf("Handle %v", spaceId)
		// 1.空间
		HandleSpace(d.db, spaceId)

		// 2.版本
		versions, _ := d.versionRepo.GetSpaceWorkVersionBySpaceId(context.Background(), spaceId)
		HandleVersion(d.db, spaceId, versions)

		// 3.成员
		members, _ := d.memberRepo.GetSpaceMemberBySpaceId(context.Background(), spaceId)
		HandleMember(d.db, spaceId, members)
	}

	return nil
}

const limit = 5000
const insertBatch = 3000

func HandleSpace(db *gorm.DB, spaceId int64) {

	var finalList []*dws.DwsSpaceWitem1h

	var pos = "2020-02-08 20:00:00"
	for {
		var list []*dws.DwsSpaceWitem1h
		err := db.Where("space_id = ? AND start_date > ?", spaceId, pos).Order("start_date asc").Limit(limit).Find(&list).Error
		if err != nil {
			fmt.Print(err)
			return
		}

		if len(list) == 0 {
			break
		}

		pos = list[len(list)-1].StartDate.Format("2006-01-02 15:04:05")

		for _, item := range list {
			if len(finalList) == 0 {
				finalList = append(finalList, item)
				continue
			}

			last := finalList[len(finalList)-1]
			if last.AbortNum == item.AbortNum &&
				last.ExpireNum == item.ExpireNum &&
				last.Num == item.Num &&
				last.CloseNum == item.CloseNum &&
				last.CompleteNum == item.CompleteNum &&
				last.TodoNum == item.TodoNum {

				last.EndDate = item.EndDate
			} else {
				finalList = append(finalList, item)
			}
		}

		if len(list) < limit {
			break
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("space_id = ?", spaceId).Delete(&dws.DwsSpaceWitem1h{}).Error
		if err != nil {
			return err
		}

		err = tx.CreateInBatches(finalList, insertBatch).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Print(err)
	}
}

func HandleVersion(db *gorm.DB, spaceId int64, versions []*version.SpaceWorkVersion) {
	var finalList []*dws.DwsVersWitem1h

	for _, v := range versions {
		finalList = append(finalList, handleVersion(db, spaceId, v.Id)...)
	}

	versionIds := stream.Map(versions, func(v *version.SpaceWorkVersion) int64 {
		return v.Id
	})

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("space_id = ? AND version_id IN ?", spaceId, versionIds).Delete(&dws.DwsVersWitem1h{}).Error
		if err != nil {
			return err
		}

		err = tx.CreateInBatches(finalList, insertBatch).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Print(err)
	}
}

func handleVersion(db *gorm.DB, spaceId int64, versionId int64) []*dws.DwsVersWitem1h {

	var finalList []*dws.DwsVersWitem1h

	var pos = "2020-02-08 20:00:00"
	for {
		var list []*dws.DwsVersWitem1h
		err := db.Where("space_id = ? AND version_id = ? AND start_date > ?", spaceId, versionId, pos).Order("start_date asc").Limit(limit).Find(&list).Error
		if err != nil {
			fmt.Print(err)
			return nil
		}

		if len(list) == 0 {
			break
		}

		pos = list[len(list)-1].StartDate.Format("2006-01-02 15:04:05")

		for _, item := range list {
			if len(finalList) == 0 {
				finalList = append(finalList, item)
				continue
			}

			last := finalList[len(finalList)-1]
			if last.VersionId == item.VersionId &&
				last.AbortNum == item.AbortNum &&
				last.ExpireNum == item.ExpireNum &&
				last.Num == item.Num &&
				last.CloseNum == item.CloseNum &&
				last.CompleteNum == item.CompleteNum &&
				last.TodoNum == item.TodoNum {

				last.EndDate = item.EndDate
			} else {
				finalList = append(finalList, item)
			}
		}

		if len(list) < limit {
			break
		}
	}

	return finalList
}

func HandleMember(db *gorm.DB, spaceId int64, members []*member.SpaceMember) {
	var finalList []*dws.DwsMbrWitem1h

	for _, v := range members {
		finalList = append(finalList, handleMember(db, spaceId, v.UserId)...)
	}

	userIds := stream.Map(members, func(v *member.SpaceMember) int64 {
		return v.UserId
	})

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("space_id = ? AND user_id IN ?", spaceId, userIds).Delete(&dws.DwsMbrWitem1h{}).Error
		if err != nil {
			return err
		}

		err = tx.CreateInBatches(finalList, insertBatch).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Print(err)
	}
}

func handleMember(db *gorm.DB, spaceId int64, userId int64) []*dws.DwsMbrWitem1h {

	var finalList []*dws.DwsMbrWitem1h

	var pos = "2020-02-08 20:00:00"
	for {
		var list []*dws.DwsMbrWitem1h
		err := db.Where("space_id = ? AND user_id = ? AND start_date > ?", spaceId, userId, pos).Order("start_date asc").Limit(limit).Find(&list).Error
		if err != nil {
			fmt.Print(err)
			return nil
		}

		if len(list) == 0 {
			break
		}

		pos = list[len(list)-1].StartDate.Format("2006-01-02 15:04:05")

		for _, item := range list {
			if len(finalList) == 0 {
				finalList = append(finalList, item)
				continue
			}

			last := finalList[len(finalList)-1]
			if last.UserId == item.UserId &&
				last.AbortNum == item.AbortNum &&
				last.ExpireNum == item.ExpireNum &&
				last.Num == item.Num &&
				last.CloseNum == item.CloseNum &&
				last.CompleteNum == item.CompleteNum &&
				last.TodoNum == item.TodoNum {

				last.EndDate = item.EndDate
			} else {
				finalList = append(finalList, item)
			}
		}

		if len(list) < limit {
			break
		}
	}

	return finalList
}
