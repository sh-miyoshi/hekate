package local

import (
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"os"
	"path/filepath"
)

// ProjectInfoHandler implement db.ProjectInfoHandler
type ProjectInfoHandler struct {
	db.ProjectInfoHandler

	filePath string
}

// NewHandler ...
func NewHandler(dbDir string) (*ProjectInfoHandler, error) {
	fileInfo, err := os.Stat(dbDir)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not directory", dbDir)
	}

	res := &ProjectInfoHandler{
		filePath: filepath.Join(dbDir, "projects.txt"),
	}
	return res, nil
}
