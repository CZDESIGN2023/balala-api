package space_file_info

import "time"

func NewSpaceFileInfo(
	spaceId int64,
	fileInfo FileInfo,
	source FileSource,
) *SpaceFileInfo {
	return &SpaceFileInfo{
		SpaceId:    spaceId,
		FileInfo:   fileInfo,
		FileSource: source,
		Status:     FileStatus_Undefined,
		CreatedAt:  time.Now().Unix(),
	}
}
