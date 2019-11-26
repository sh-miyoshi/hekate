package db

import (
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
)

// ProjectInfoHandler ...
type ProjectInfoHandler interface {
	Add(ent *model.ProjectInfo) error
	Delete(id string) error
	GetList() ([]string, error)
	Get(id string) (*model.ProjectInfo, error)
	Update(ent *model.ProjectInfo) error
}

// UserInfoHandler ...
type UserInfoHandler interface {
	Add(ent *model.UserInfo) error
	Delete(projectID string, userID string) error
	GetList(projectID string) ([]string, error)
	Get(projectID string, userID string) (*model.UserInfo, error)
	Update(ent *model.UserInfo) error
	GetIDByName(projectID string, userName string) (string, error)
}
