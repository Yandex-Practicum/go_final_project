package storage

import (
	"fmt"
	"path/filepath"
	"todo-list/internal/lib/common"
)

func DBFilePath(localPath string) (string, error) {

	appPath, err := common.ProjectRootPath()
	if err != nil {
		return "", fmt.Errorf("failed to get project root path: %w", err)
	}

	return filepath.Join(appPath, localPath), nil

}
