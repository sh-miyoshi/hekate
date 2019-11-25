package db

import (
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
)

// ProjectInfoHandler ...
type ProjectInfoHandler interface {
	Add(ent *model.ProjectInfo) error
	Delete(id string) error
	GetList() []string
	Get(id string) (*model.ProjectInfo, error)
	Update(ent *model.ProjectInfo) error
}

// UserInfoHandler ...
type UserInfoHandler interface {
	Add(ent *model.UserInfo) error
	Delete(id string) error
	GetList() []string
	Get(id string) (*model.UserInfo, error)
	Update(ent *model.UserInfo) error
}
