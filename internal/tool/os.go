package tool

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type isExistsResult struct {
	Exist bool
	IsDir bool
}

func IsExists(path string) (ret isExistsResult, err error) {
	stat, err := os.Stat(path)
	// 存在
	if err == nil {
		ret.Exist = true
		// 是目录
		if stat.IsDir() {
			ret.IsDir = true
		} else {
			ret.IsDir = false
		}
		return
	}
	// 如果是不存在错误
	if os.IsNotExist(err) {
		ret.Exist = false
		err = nil
		return
	}

	return
}

// GetDirByIndex 获取指定路径 指定索引的目录
func GetDirByIndex(path string, index int) (result, fullPath string, err error) {
	sep := string(filepath.Separator)
	path = strings.Trim(path, sep)
	arr := strings.Split(path, sep)
	if index >= len(arr) {
		err = fmt.Errorf("get dir from \"%s\" failed, path length is %d, but index is %d, index out. ", path, len(arr), index)
		return
	}
	result = arr[index]
	fullPath = strings.Join(arr[:index+1], sep)
	return
}

func RemoveEmptyDirAll(path string, deleteRoot bool) (ret []string, err error) {
	exists, err := IsExists(path)
	if err != nil {
		return
	}
	if !exists.Exist || !exists.IsDir {
		return
	}

	// 获取目录下的所有空目录
	checkDirList := make([]string, 0)
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			checkDirList = append(checkDirList, path)
		}
		return nil
	})
	if err != nil {
		return
	}

	endIndex := 1
	if deleteRoot {
		endIndex = 0
	}
	for i := len(checkDirList) - 1; i >= endIndex; i-- {
		dir, err := ioutil.ReadDir(checkDirList[i])
		if err != nil {
			return ret, err
		}
		if len(dir) == 0 {
			if err = os.Remove(checkDirList[i]); err != nil {
				ret = append(ret, checkDirList[i])
				return ret, err
			}
		}
	}
	return
}
