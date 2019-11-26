package local

import (
	"encoding/csv"
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"os"
	"path/filepath"
)

// UserInfoHandler implement db.UserInfoHandler
type UserInfoHandler struct {
	fileDir string
}

// NewUserHandler ...
func NewUserHandler(dbDir string) (*UserInfoHandler, error) {
	fileInfo, err := os.Stat(dbDir)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not directory", dbDir)
	}

	res := &UserInfoHandler{
		fileDir: dbDir,
	}
	return res, nil
}

// Add ...
func (h *UserInfoHandler) Add(ent *model.UserInfo) error {
	if err := ent.Validate(); err!=nil {
		return err
	}

	fp, err := os.OpenFile(h.getFilePath(ent.ProjectID), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	for {
		user, err := reader.Read()
		if err != nil {
			// TODO(consider not EOF err)
			break
		}
		if user[0] == ent.ID || user[2] == ent.Name {
			return model.ErrUserAlreadyExists
		}
	}

	// append new project
	fmt.Fprintf(fp, "%s,%s,%s,%t,%s,%s", ent.ID, ent.ProjectID, ent.Name, ent.Enabled, ent.CreatedAt, ent.PasswordHash)
	for _, role := range ent.Roles {
		fmt.Fprintf(fp, ",%s", role)
	}
	fmt.Fprintf(fp, "\n")

	return nil
}

// Delete ...
func (h *UserInfoHandler) Delete(id string) error {
	// TODO(not implemented yet)
	return nil
}

// GetList ...
func (h *UserInfoHandler) GetList() ([]string, error) {
	// TODO(not implemented yet)
	return []string{}, nil
}

// Get ...
func (h *UserInfoHandler) Get(id string) (*model.UserInfo, error) {
	// TODO(not implemented yet)
	return nil, nil
}

// Update ...
func (h *UserInfoHandler) Update(ent *model.UserInfo) error {
	// TODO(not implemented yet)
	return nil
}

func (h *UserInfoHandler) getFilePath(projectID string) string {
	return filepath.Join(h.fileDir, projectID + "_users.csv")
}