package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// ProjectInfoHandler implement db.ProjectInfoHandler
type ProjectInfoHandler struct {
	projectList map[string]*model.ProjectInfo
}

// NewProjectHandler ...
func NewProjectHandler() *ProjectInfoHandler {
	res := &ProjectInfoHandler{
		projectList: make(map[string]*model.ProjectInfo),
	}
	return res
}

// Add ...
func (h *ProjectInfoHandler) Add(ent *model.ProjectInfo) *errors.Error {
	h.projectList[ent.Name] = ent
	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(name string) *errors.Error {
	if _, exists := h.projectList[name]; exists {
		delete(h.projectList, name)
		return nil
	}
	return model.ErrNoSuchProject
}

// GetList ...
func (h *ProjectInfoHandler) GetList() ([]*model.ProjectInfo, *errors.Error) {
	res := []*model.ProjectInfo{}
	for _, prj := range h.projectList {
		res = append(res, prj)
	}
	return res, nil
}

// Get ...
func (h *ProjectInfoHandler) Get(name string) (*model.ProjectInfo, *errors.Error) {
	res, ok := h.projectList[name]
	if ok {
		return res, nil
	}
	return nil, model.ErrNoSuchProject
}

// Update ...
func (h *ProjectInfoHandler) Update(ent *model.ProjectInfo) *errors.Error {
	if _, ok := h.projectList[ent.Name]; ok {
		h.projectList[ent.Name] = ent
		return nil
	}
	return model.ErrNoSuchProject
}
