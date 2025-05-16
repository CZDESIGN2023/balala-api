package file

import (
	"cmp"
	"github.com/spf13/cast"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func AllFilesOfDir(dirPath string) ([]*os.File, error) {
	var files []*os.File
	err := filepath.WalkDir(dirPath, func(path string, info fs.DirEntry, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		files = append(files, file)
		return nil
	})

	return files, err
}

func AllFilePathOfDir(dirPath string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(dirPath, func(path string, info fs.DirEntry, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		paths = append(paths, path)
		return nil
	})

	return paths, err
}

func AllFileNamesOfDir(dirPath string) ([]string, error) {
	var names []string
	err := filepath.WalkDir(dirPath, func(path string, info fs.DirEntry, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		names = append(names, info.Name())
		return nil
	})

	return names, err
}

func SortChunks(chunks []*os.File) {
	slices.SortFunc(chunks, func(i, j *os.File) int {
		iIdx := cast.ToInt64(strings.Split(filepath.Base(i.Name()), "-")[1])
		jIdx := cast.ToInt64(strings.Split(filepath.Base(j.Name()), "-")[1])
		return cmp.Compare(iIdx, jIdx)
	})
}

func SortPath(paths []string) {
	slices.SortFunc(paths, func(i, j string) int {
		iIdx := cast.ToInt64(strings.Split(filepath.Base(i), "-")[1])
		jIdx := cast.ToInt64(strings.Split(filepath.Base(j), "-")[1])
		return cmp.Compare(iIdx, jIdx)
	})
}

func Exist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func IsNotExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return os.IsNotExist(err)
}

/* filePath /space_file/20240528/706/example.png
* return 706
 */
func extractSpaceIdFromPath(filePath string) int64 {
	endIdx := strings.LastIndex(filePath, "/")
	if endIdx == -1 {
		return 0
	}

	startIdx := strings.LastIndex(filePath[:endIdx], "/")
	if startIdx == -1 {
		return 0
	}

	spaceId := filePath[startIdx+1 : endIdx]
	return cast.ToInt64(spaceId)
}
