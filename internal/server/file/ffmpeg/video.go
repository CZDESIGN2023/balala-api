package ffmpeg

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"path/filepath"
	"strings"
)

func ExtractFirstFrame(input string, output string) bool {
	ext := strings.ToLower(filepath.Ext(input))
	if !IsSupportedVideoExt(ext) {
		return false
	}
	_ = ffmpeg.Input(input, ffmpeg.KwArgs{"ss": 0}).
		Filter("select", ffmpeg.Args{"eq(n,0)"}).
		Output(output, ffmpeg.KwArgs{"vframes": 1, "q:v": 10}).Run() // 流式输出

	return fileExist(output)
}
