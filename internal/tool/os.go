package tool

import (
	"os"
)

func IsFileExists(path string) (res bool, err error) {
	_, err = os.Stat(path)

	if err == nil {
		return true, err
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
