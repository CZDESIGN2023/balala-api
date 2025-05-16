package ffmpeg

import (
	"slices"
)

var videoExtList = []string{".mp4", ".m4v", ".mov", ".avi", ".flv", ".mkv", ".wmv", ".rmvb", ".mpg", ".rm"}
var audioExtList = []string{".mp3", ".wav", ".flac", ".m4a", ".ogg", ".aac", ".wma", ".ape"}

func IsSupportedVideoExt(ext string) bool {
	return slices.Contains(videoExtList, ext)
}

func IsSupportedAudioExt(ext string) bool {
	return slices.Contains(audioExtList, ext)
}
