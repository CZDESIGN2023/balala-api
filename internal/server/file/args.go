package file

import "mime/multipart"

type SpaceUploadArgs struct {
	SpaceId  int64  `form:"spaceId"`
	Scene    string `form:"scene"`
	SubScene string `form:"subScene"`
}

type VerifyArgs struct {
	SpaceUploadArgs
	FileHash string `form:"fileHash" binding:"required"`
	FileName string `form:"fileName" binding:"required"`
}

type MergeArgs struct {
	SpaceUploadArgs
	ChunkSize int64  `form:"chunkSize"`
	FileName  string `form:"fileName" binding:"required"`
	FileHash  string `form:"fileHash" binding:"required"`
}

type UploadArgs struct {
	SpaceId  int64                   `form:"spaceId"`
	Scene    string                  `form:"scene" binding:"required"`
	SubScene string                  `form:"sub_scene"`
	Files    []*multipart.FileHeader `form:"files" binding:"required"`
}

type Upload2Args struct {
	FileHash    string                `form:"fileHash" binding:"required"`
	FileName    string                `form:"fileName" binding:"required"`
	FileSize    int64                 `form:"fileSize" binding:"required"`
	Index       int64                 `form:"index"`
	ChunkFile   *multipart.FileHeader `form:"chunkFile" binding:"required"`
	ChunkSize   int64                 `form:"chunkSize" binding:"required"`
	ChunkNumber int64                 `form:"chunkNumber" binding:"required"`
	ChunkHash   string                `form:"chunkHash" binding:"required"`
	Finish      bool                  `form:"finish"`
}

type GetDownloadTokenArgs struct {
	Scene   string `form:"scene" binding:"required"`
	SpaceId int64  `form:"spaceId"`
	Id      int64  `form:"id" binding:"required"`
}

type DownloadArgs struct {
	Scene         string `form:"scene" binding:"required"`
	DownloadToken string `form:"downloadToken" binding:"required"`
}

type ProxyArgs struct {
	Url string `form:"url" binding:"required"`
}
