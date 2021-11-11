package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

func prepLogFile(dirName, fileName string) (*os.File, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return new(os.File), fmt.Errorf("get user cache dir: %s", err)
	}

	err = os.Chdir(cacheDir)
	if err != nil {
		Logger.Errorf("changing to cache directory: %s", err)
		return new(os.File), fmt.Errorf("chdir to user cache dir: %s", err)
	}

	if err := createIfNotExist(dirName); err != nil {
		return new(os.File), fmt.Errorf("create subdir in user cache dir: %s", err)
	}

	logFilePath := filepath.Join(cacheDir, dirName, fileName)

	return os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
}

func createIfNotExist(dirName string) error {
	_, err := os.Stat(dirName)
	if err != nil {
		return fmt.Errorf("get \"%s\" stat: %s", dirName, err)
	}

	if os.IsNotExist(err) {
		err = os.Mkdir(dirName, 0755)
		return fmt.Errorf("make \"%s\" dir: %s", dirName, err)
	}

	return nil
}
