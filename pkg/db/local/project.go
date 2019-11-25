package local

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ProjectInfoHandler implement db.ProjectInfoHandler
type ProjectInfoHandler struct {
	filePath string
	mu       sync.Mutex
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
	h.mu.Lock()
	defer h.mu.Unlock()

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
		if project[0] == ent.ID || project[1] == ent.Name {
			return model.ErrProjectAlreadyExists
		}
	}

	// append new project
	fmt.Fprintf(fp, "%s,%s,%t,%s,%d,%d\n", ent.ID, ent.Name, ent.Enabled, ent.CreatedAt, ent.TokenConfig.AccessTokenLifeSpan, ent.TokenConfig.RefreshTokenLifeSpan)

	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if id == "" {
		return fmt.Errorf("id of entry is empty")
	}

	if id == "master" {
		return fmt.Errorf("master project can not delete")
	}

	fp, err := os.OpenFile(h.filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	deleted := false
	scanner := bufio.NewScanner(fp)
	results := []string{}

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
			results = append(results, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if !deleted {
		return model.ErrNoSuchProject
	}

	// Remove all data at first
	fp.Truncate(0)
	fp.Seek(0, 0)

	// Write results
	for _, line := range results {
		fmt.Fprintln(fp, line)
	}

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
