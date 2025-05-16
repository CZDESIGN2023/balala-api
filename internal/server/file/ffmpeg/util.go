package ffmpeg

import "os"

func fileExist(f string) bool {
	_, err := os.Stat(f)
	return !os.IsNotExist(err)
}
