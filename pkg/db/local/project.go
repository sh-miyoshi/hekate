package local

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"os"
	"path/filepath"
	"strconv"
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
		return nil, errors.Wrap(err, "No such dir")
	}
	if !fileInfo.IsDir() {
		return nil, errors.Cause(fmt.Errorf("%s is not directory", dbDir))
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
		return errors.Cause(fmt.Errorf("id of entry is empty"))
	}

	fp, err := os.OpenFile(h.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to open file")
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
			return errors.Cause(model.ErrProjectAlreadyExists)
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
		return errors.Cause(fmt.Errorf("id of entry is empty"))
	}

	if id == "master" {
		return errors.Cause(fmt.Errorf("master project can not delete"))
	}

	fp, err := os.OpenFile(h.filePath, os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to open file")
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
			return errors.Wrap(err, "Failed to read line")
		}
		if project[0] == id {
			deleted = true
		} else {
			results = append(results, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "Failed to read csv")
	}

	if !deleted {
		return errors.Cause(model.ErrNoSuchProject)
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
func (h *ProjectInfoHandler) GetList() ([]string, error) {
	fp, err := os.Open(h.filePath)
	if err != nil {
		return []string{}, errors.Wrap(err, "Failed to open project file")
	}
	defer fp.Close()

	results := []string{}

	reader := csv.NewReader(fp)
	for {
		project, err := reader.Read()
		if err != nil {
			// TODO(consider not EOF err)
			break
		}
		results = append(results, project[0])
	}

	return results, nil
}

// Get ...
func (h *ProjectInfoHandler) Get(id string) (*model.ProjectInfo, error) {
	fp, err := os.Open(h.filePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open project file")
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	for {
		project, err := reader.Read()
		if err != nil {
			// TODO(consider not EOF err)
			break
		}
		enabled, _ := strconv.ParseBool(project[2])
		accessTokenLifeSpan, _ := strconv.Atoi(project[4])
		refreshTokenLifeSpan, _ := strconv.Atoi(project[5])

		if project[0] == id {
			return &model.ProjectInfo{
				ID:        project[0],
				Name:      project[1],
				Enabled:   enabled,
				CreatedAt: project[3],
				TokenConfig: &model.TokenConfig{
					AccessTokenLifeSpan:  accessTokenLifeSpan,
					RefreshTokenLifeSpan: refreshTokenLifeSpan,
				},
			}, nil
		}
	}
	return nil, errors.Cause(model.ErrNoSuchProject)
}

// Update ...
func (h *ProjectInfoHandler) Update(ent *model.ProjectInfo) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if ent.ID == "" {
		return errors.Cause(fmt.Errorf("id of entry is empty"))
	}

	fp, err := os.OpenFile(h.filePath, os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to open file")
	}
	defer fp.Close()

	updated := false
	scanner := bufio.NewScanner(fp)
	results := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		reader := csv.NewReader(strings.NewReader(line))
		project, err := reader.Read()
		if err != nil {
			return errors.Wrap(err, "Failed to read line")
		}
		if project[0] == ent.ID {
			updated = true
			newData := fmt.Sprintf("%s,%s,%t,%s,%d,%d\n", ent.ID, ent.Name, ent.Enabled, ent.CreatedAt, ent.TokenConfig.AccessTokenLifeSpan, ent.TokenConfig.RefreshTokenLifeSpan)
			results = append(results, newData)
		} else {
			results = append(results, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "Failed to read csv")
	}

	if !updated {
		return errors.Cause(model.ErrNoSuchProject)
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
