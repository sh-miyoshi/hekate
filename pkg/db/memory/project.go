package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
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
func (h *ProjectInfoHandler) Add(ent *model.ProjectInfo) error {
	h.projectList[ent.Name] = ent
	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(name string) error {
	if _, exists := h.projectList[name]; exists {
		delete(h.projectList, name)
		return nil
	}
	return model.ErrNoSuchProject
}

// GetList ...
func (h *ProjectInfoHandler) GetList() ([]*model.ProjectInfo, error) {
	res := []*model.ProjectInfo{}
	for _, prj := range h.projectList {
		res = append(res, prj)
	}
	return res, nil
}

// Get ...
func (h *ProjectInfoHandler) Get(name string) (*model.ProjectInfo, error) {
	res, ok := h.projectList[name]
	if ok {
		return res, nil
	}
	return nil, model.ErrNoSuchProject
}

// Update ...
func (h *ProjectInfoHandler) Update(ent *model.ProjectInfo) error {
	if _, ok := h.projectList[ent.Name]; ok {
		h.projectList[ent.Name] = ent
		return nil
	}
	return model.ErrNoSuchProject
}
