package file_info

import shared "go-cs/internal/pkg/domain"

type FileInfos []*FileInfo

type FileInfo struct {
	shared.AggregateRoot

	Id        int64  `json:"id,omitempty"`
	Hash      string `json:"hash,omitempty"`
	Name      string `json:"name,omitempty"`
	Typ       int32  `json:"typ,omitempty"` // 文件類型 (1-图片，2-视频，3-音频，4-文本)
	Size      int64  `json:"size,omitempty"`
	Uri       string `json:"uri,omitempty"` // 前端讀取檔案用, 動態產生: /path/name?sign=hmac-sha256簽名(附加過期timestamp)
	Pwd       string `json:"pwd,omitempty"`
	Cover     string `json:"cover,omitempty"`
	Status    int32  `json:"status,omitempty"` // 0-未定義, 1-初始化, 2-上传中, 3-成功, 4-失败, 5-处理中（转码等）, 6-待删
	Owner     int64  `json:"owner,omitempty"`  // 上傳者id
	Meta      string `json:"meta,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"` //创建时间
	UpdatedAt int64  `json:"updated_at,omitempty"` //更新时间
	DeletedAt int64  `json:"deleted_at,omitempty"` //删除时间
	// 以下欄位不應回傳到前端, 於api返回值清空
	UploadTyp    int32  `json:"upload_typ,omitempty"`    // 上傳服務器類型 (0-local 1-fastdfs，2-s3)
	UploadDomain string `json:"upload_domain,omitempty"` // 上傳服務器
	UploadMd5    string `json:"upload_md5,omitempty"`    // 上傳完成md5
	UploadPath   string `json:"upload_path,omitempty"`   // 上傳完成檔案路徑
}

func (s *FileInfo) IsOwner(userId int64) bool {
	return s.Owner == userId
}
