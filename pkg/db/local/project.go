package local

import (
	"encoding/csv"
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"os"
	"path/filepath"
	"strings"
	"bufio"
)

// ProjectInfoHandler implement db.ProjectInfoHandler
type ProjectInfoHandler struct {
	filePath string
}

// NewProjectHandler ...
func NewProjectHandler(dbDir string) (*ProjectInfoHandler, error) {
	fileInfo, err := os.Stat(dbDir)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not directory", dbDir)
	}

	res := &ProjectInfoHandler{
		filePath: filepath.Join(dbDir, "projects.csv"),
	}
	return res, nil
}

// Add ...
func (h *ProjectInfoHandler) Add(ent *model.ProjectInfo) error {
	if ent.ID == "" {
		return fmt.Errorf("id of entry is empty")
	}

	fp, err := os.OpenFile(h.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	for {
		project, err := reader.Read()
		if err != nil {
			// TODO(consider not EOF err)
			break
		}
		if project[0] == ent.ID {
			return model.ErrProjectAlreadyExists
		}
	}

	// append new project
	fmt.Fprintf(fp, "%s,%s\n", ent.ID, ent.Name)

	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(id string) error {
	if id == "" {
        return fmt.Errorf("id of entry is empty")
    }

	if id == "master" {
		return fmt.Errorf("master project can not delete")
	}

    fp, err := os.OpenFile(h.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    defer fp.Close()

	deleted := false
	scanner := bufio.NewScanner(fp)
	result := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		reader := csv.NewReader(strings.NewReader(line))
		project, err := reader.Read()
		if err != nil {
			return err
		}
		if project[0] == id {
            deleted = true
        } else {
			result = append(result, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if !deleted {
		return fmt.Errorf("No such project")
	}

	// TODO(write result)

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
