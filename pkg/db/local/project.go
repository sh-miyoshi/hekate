package local

import (
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"os"
	"path/filepath"
)

// ProjectInfoHandler implement db.ProjectInfoHandler
type ProjectInfoHandler struct {
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

// Add ...
func (h *ProjectInfoHandler) Add(ent *model.ProjectInfo) error {
	// TODO(not implemented yet)
	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(id string) error {
	// TODO(not implemented yet)
	return nil
}

// GetList ...
func (h *ProjectInfoHandler) GetList() []string {
	// TODO(not implemented yet)
	return []string{}
}

// Get ...
func (h *ProjectInfoHandler) Get(id string) (*model.ProjectInfo, error) {
	// TODO(not implemented yet)
	return nil, nil
}

// Update ...
func (h *ProjectInfoHandler) Update(ent *model.ProjectInfo) error {
	// TODO(not implemented yet)
	return nil
}
