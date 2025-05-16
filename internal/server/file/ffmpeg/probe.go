package ffmpeg

import (
	"context"
	"errors"
	"gopkg.in/vansante/go-ffprobe.v2"
	"path/filepath"
	"strings"
)

// Probe 探测输入音视频信息
func Probe(input string) (*ffprobe.ProbeData, error) {
	ext := strings.ToLower(filepath.Ext(input))

	var codecType string
	switch {
	case IsSupportedVideoExt(ext):
		codecType = "v"
	case IsSupportedAudioExt(ext):
		codecType = "a"
	default:
		return nil, errors.New("unsupported file type")
	}

	data, err := ffprobe.ProbeURL(context.Background(), input, []string{
		"-select_streams", codecType,
	}...)
	if err != nil {
		return nil, err
	}

	return data, nil
}
