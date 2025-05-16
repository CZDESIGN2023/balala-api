package ffmpeg

import (
	"go-cs/pkg/pprint"
	"os"
	"testing"
)

func TestProbe(t *testing.T) {
	dir, _ := os.UserHomeDir()
	abs := dir + "/Downloads/GHTestVideoDemo-master/Video/8.flv"
	probe, err := Probe(abs)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(probe)
}
