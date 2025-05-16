package ioutil

import (
	"io"
	"os"
	"path"
)

func WriteFileAtomic(filename string, data []byte, perm os.FileMode) error {
	dir, name := path.Dir(filename), path.Base(filename)
	f, err := os.CreateTemp(dir, name)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	f.Close()
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	} else {
		err = os.Chmod(f.Name(), perm)
	}
	if err != nil {
		os.Remove(f.Name())
		return err
	}
	return os.Rename(f.Name(), filename)
}
