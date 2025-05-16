package space_file_info

import (
	shared "go-cs/internal/pkg/domain"
	"time"
)

type FileStatus int32

var (
	FileStatus_Undefined = FileStatus(0)
	FileStatus_Using     = FileStatus(1)
	FileStatus_Deleted   = FileStatus(2)
)

type FileSourceType int32

var (
	FileSourceType_Undefined = FileSourceType(0)
	FileSourceType_Task      = FileSourceType(1)
)

type FileSource struct {
	SourceType FileSourceType `json:"source_type,omitempty"` //关联类型 1:任务
	SourceId   int64          `json:"source_id,omitempty"`   //关联id ,
}

type FileInfo struct {
	FileInfoId int64  `json:"file_info_id,omitempty"` //关联的文件存储ID
	FileName   string `json:"file_name,omitempty"`    //标签id
	FileUri    string `json:"file_uri,omitempty"`     //标签id
	FileSize   int64  `json:"file_size,omitempty"`    //文件大小
}

type SpaceFileInfos []*SpaceFileInfo

type SpaceFileInfo struct {
	shared.AggregateRoot

	Id         int64      `json:"id,omitempty"`
	SpaceId    int64      `json:"space_id,omitempty"`  //空间id
	FileInfo   FileInfo   `json:"file_info,omitempty"` //关联的文件存储ID
	FileSource FileSource `json:"source,omitempty"`
	Status     FileStatus `json:"status,omitempty"`     //状态 0-未定义 1-使用中 2-已删除
	CreatedAt  int64      `json:"created_at,omitempty"` //创建时间
	UpdatedAt  int64      `json:"updated_at,omitempty"` //更新时间
	DeletedAt  int64      `json:"deleted_at,omitempty"` //删除时间
}

func (s *SpaceFileInfo) OnDelete() {
	//触发领域事件, 物理删除文件
	s.DeletedAt = time.Now().Unix()
	s.Status = FileStatus_Deleted

	s.AddDiff(Diff_Deleted)
}
