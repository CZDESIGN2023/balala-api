package data

import (
	"errors"
	"go-cs/pkg/stream"
	"gorm.io/gorm"
	"time"
)

type RptRepo struct {
	data *DwhData
}

func NewRptRepo(data *DwhData) *RptRepo {
	return &RptRepo{
		data: data,
	}
}

func (r *RptRepo) SpaceMemberMap(spaceIds []int64, endDate time.Time) (map[int64][]int64, error) {
	end := endDate.Format("2006-01-02 15:04:05")

	const sql = `
SELECT * FROM (
	SELECT 
		*,
		ROW_NUMBER() over ( PARTITION BY space_id, user_id ORDER BY _id DESC ) AS ranking 
	FROM dwd_member
	WHERE start_date < ? AND end_date >= ? 
	AND space_id IN ?
) AS t 
WHERE t.ranking = 1
`

	type spaceMemberSearchResult struct {
		SpaceId int64 `gorm:"column:space_id"`
		UserId  int64 `gorm:"column:user_id"`
	}

	//获取符合时间段的空间成员
	var spaceMemberSearchResults []*spaceMemberSearchResult
	err := r.data.Db().Raw(sql, end, end, spaceIds).Scan(&spaceMemberSearchResults).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	//按空间分组
	spaceMap := stream.GroupBy(spaceMemberSearchResults, func(v *spaceMemberSearchResult) int64 {
		return v.SpaceId
	})

	ret := stream.MapValue(spaceMap, func(v []*spaceMemberSearchResult) []int64 {
		return stream.Unique(stream.Map(v, func(v *spaceMemberSearchResult) int64 {
			return v.UserId
		}))
	})

	return ret, nil
}
