package local

import (
	"encoding/csv"
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"os"
	"path/filepath"
	"strconv"
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
	if err := ent.Validate(); err != nil {
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
func (h *UserInfoHandler) Delete(projectID string, userID string) error {
	// TODO(not implemented yet)
	return nil
}

// GetList ...
func (h *UserInfoHandler) GetList(projectID string) ([]string, error) {
	// TODO(not implemented yet)
	return []string{}, nil
}

// Get ...
func (h *UserInfoHandler) Get(projectID string, userID string) (*model.UserInfo, error) {
	fp, err := os.Open(h.getFilePath(projectID))
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	for {
		user, err := reader.Read()
		if err != nil {
			// TODO(consider not EOF err)
			break
		}
		if user[0] == userID {
			enabled, _ := strconv.ParseBool(user[3])
			res := &model.UserInfo{
				ID:           user[0],
				ProjectID:    user[1],
				Name:         user[2],
				Enabled:      enabled,
				CreatedAt:    user[4],
				PasswordHash: user[5],
			}
			for i := 6; i < len(user); i++ {
				res.Roles = append(res.Roles, user[i])
			}
			return res, nil
		}
	}

	return nil, model.ErrNoSuchUser
}

// Update ...
func (h *UserInfoHandler) Update(ent *model.UserInfo) error {
	// TODO(not implemented yet)
	return nil
}

// GetIDByName ...
func (h *UserInfoHandler) GetIDByName(projectID string, userName string) (string, error) {
	fp, err := os.Open(h.getFilePath(projectID))
	if err != nil {
		return "", err
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	for {
		user, err := reader.Read()
		if err != nil {
			// TODO(consider not EOF err)
			break
		}
		if user[2] == userName {
			return user[0], nil
		}
	}

	return "", model.ErrNoSuchUser
}

func (h *UserInfoHandler) getFilePath(projectID string) string {
	return filepath.Join(h.fileDir, projectID+"_users.csv")
}
