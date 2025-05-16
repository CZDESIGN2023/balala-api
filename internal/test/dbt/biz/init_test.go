package biz

import (
	"context"
	"fmt"
	"go-cs/internal/dwh/model/dws"
	"go-cs/internal/test/dbt"
	"gorm.io/gorm"
	"testing"
)

func init() {
	// 测试执行目录为当前目录，所以项目配置文件路径这么奇怪
	dbt.Init("../../../../configs/config.yaml", "../../../../.env.local", true)
}

func Test(t *testing.T) {

	spaceIds, err := dbt.R.SpaceRepo.GetAllSpaceIds()
	if err != nil {
		fmt.Print(err)
		return
	}
	for _, spaceId := range spaceIds {
		HandleSpace(dbt.Data.DB(context.Background()), spaceId)

		versions, _ := dbt.R.SpaceWorkVersionRepo.GetSpaceWorkVersionBySpaceId(context.Background(), spaceId)

		for _, version := range versions {
			HandleVersion(dbt.Data.DB(context.Background()), spaceId, version.Id)
		}

		members, _ := dbt.R.SpaceMemberRepo.GetSpaceMemberBySpaceId(context.Background(), spaceId)

		for _, v := range members {
			HandleMember(dbt.Data.DB(context.Background()), spaceId, v.UserId)
		}
	}
}

func Test2(t *testing.T) {
	HandleMember(dbt.Data.DB(context.Background()), 87, 42)
}

func HandleSpace(db *gorm.DB, spaceId int64) {

	var finalList []*dws.DwsSpaceWitem1h

	var pos = "2020-02-08 20:00:00"
	for {
		var list []*dws.DwsSpaceWitem1h
		err := db.Where("space_id = ? AND start_date > ?", spaceId, pos).Order("start_date asc").Limit(1000).Find(&list).Error
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

		if len(list) < 1000 {
			break
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("space_id = ?", spaceId).Delete(&dws.DwsSpaceWitem1h{}).Error
		if err != nil {
			return err
		}

		err = tx.CreateInBatches(finalList, 1000).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Print(err)
	}
}

func HandleVersion(db *gorm.DB, spaceId int64, versionId int64) {

	var finalList []*dws.DwsVersWitem1h

	var pos = "2020-02-08 20:00:00"
	for {
		var list []*dws.DwsVersWitem1h
		err := db.Where("space_id = ? AND version_id = ? AND start_date > ?", spaceId, versionId, pos).Order("start_date asc").Limit(1000).Find(&list).Error
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

		if len(list) < 1000 {
			break
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("space_id = ? AND version_id = ?", spaceId, versionId).Delete(&dws.DwsVersWitem1h{}).Error
		if err != nil {
			return err
		}

		err = tx.CreateInBatches(finalList, 1000).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Print(err)
	}
}

func HandleMember(db *gorm.DB, spaceId int64, userId int64) {

	var finalList []*dws.DwsMbrWitem1h

	var pos = "2020-02-08 20:00:00"
	for {
		var list []*dws.DwsMbrWitem1h
		err := db.Where("space_id = ? AND user_id = ? AND start_date > ?", spaceId, userId, pos).Order("start_date asc").Limit(1000).Find(&list).Error
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

		if len(list) < 1000 {
			break
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("space_id = ? AND user_id = ?", spaceId, userId).Delete(&dws.DwsMbrWitem1h{}).Error
		if err != nil {
			return err
		}

		err = tx.CreateInBatches(finalList, 1000).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Print(err)
	}
}
