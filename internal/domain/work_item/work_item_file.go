package work_item

type FileInfo struct {
	FileInfoId int64  `json:"file_info_id,omitempty"` //关联的文件存储ID
	FileName   string `json:"file_name,omitempty"`    //标签id
	FileUri    string `json:"file_uri,omitempty"`     //标签id
	FileSize   int64  `json:"file_size,omitempty"`    //文件大小
}

type WorkItemFiles []*WorkItemFile

type WorkItemFile struct {
	Id         int64    `json:"id,omitempty"`
	WorkItemId int64    `json:"work_item_id,omitempty"` //工作项id
	SpaceId    int64    `json:"space_id,omitempty"`     //空间id
	FileInfo   FileInfo `json:"file_info,omitempty"`    //关联的文件存储ID
	CreatedAt  int64    `json:"created_at,omitempty"`   //创建时间
	UpdatedAt  int64    `json:"updated_at,omitempty"`   //更新时间
	DeletedAt  int64    `json:"deleted_at,omitempty"`   //删除时间
	Status     int32    `json:"status,omitempty"`
}
