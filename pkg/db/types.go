package db

import (
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
)

// ProjectInfoHandler ...
type ProjectInfoHandler interface {
	Add(ent *model.ProjectInfo) error
	Delete(name string) error
	GetList() ([]string, error)
	Get(name string) (*model.ProjectInfo, error)
	Update(ent *model.ProjectInfo) error
}

// UserInfoHandler ...
type UserInfoHandler interface {
	Add(ent *model.UserInfo) error
	Delete(projectName string, userID string) error
	GetList(projectName string) ([]string, error)
	Get(projectName string, userID string) (*model.UserInfo, error)
	Update(ent *model.UserInfo) error
	GetIDByName(projectName string, userName string) (string, error)
	DeleteProjectDefine(projectName string) error
}
