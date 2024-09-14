package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

const rootFolder = "todo-list"

func DBFilePath(localPath string) (string, error) {

	appPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to perform os.Executable(): %w", err)
	}

	var rootAppPath string
	for _, f := filepath.Split(appPath); f != ""; {
		if f == rootFolder {
			rootAppPath = appPath
			break
		}
		appPath = filepath.Dir(appPath)
		_, f = filepath.Split(appPath)
	}
	if rootAppPath == "" {
		return "", fmt.Errorf("failed to get root path \"%s\" of the app", rootFolder)
	}

	return filepath.Join(appPath, localPath), nil

}
