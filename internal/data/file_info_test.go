package data

import (
	"context"
	"go-cs/pkg/pprint"
	"testing"
)

func TestFileInfoRepo_GetFileInfoByIds(t *testing.T) {
	res, err := fileInfoRepo.GetFileInfoByIds(context.Background(), []int64{11, 12})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

type versionSearchResult struct {
	SpaceId   int64 `gorm:"column:space_id" json:"space_id"`
	VersionId int64 `gorm:"column:version_id" json:"version_id"`
}

func Test123(t *testing.T) {
	const bizVersionSearchSql = `
SELECT * FROM (
	SELECT
		*,
		ROW_NUMBER() over ( PARTITION BY version_id ORDER BY _id DESC ) AS ranking
	FROM dim_version
	WHERE start_date <= ? AND end_date >= ?
) AS t
WHERE t.ranking = 1
		`

	start := "2025-02-19 16:00:00"
	end := "2025-02-19 17:00:00"

	var list []*versionSearchResult
	err := gdb.Raw(bizVersionSearchSql, start, end).Scan(&list).Error
	if err != nil {
		t.Error(err)
	}

	pprint.Println(list)
}
