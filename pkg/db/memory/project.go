package memory

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
)

// ProjectInfoHandler implement db.ProjectInfoHandler
type ProjectInfoHandler struct {
	projectList map[string]model.ProjectInfo
}

// NewProjectHandler ...
func NewProjectHandler() (*ProjectInfoHandler, error) {
	res := &ProjectInfoHandler{
		projectList: make(map[string]model.ProjectInfo),
	}
	return res, nil
}

// Add ...
func (h *ProjectInfoHandler) Add(ent *model.ProjectInfo) error {
	if ent.Name == "" {
		return errors.Cause(fmt.Errorf("name of entry is empty"))
	}

	if _, exists := h.projectList[ent.Name]; exists {
		return errors.Cause(model.ErrProjectAlreadyExists)
	}

	for _, project := range h.projectList {
		if project.Name == ent.Name {
			return errors.Cause(model.ErrProjectAlreadyExists)
		}
	}

	h.projectList[ent.Name] = *ent
	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(name string) error {
	if name == "" {
		return errors.Cause(fmt.Errorf("name of entry is empty"))
	}

	if name == "master" {
		return errors.Cause(fmt.Errorf("master project can not delete"))
	}

	if _, exists := h.projectList[name]; exists {
		delete(h.projectList, name)
		return nil
	}
	return errors.Cause(model.ErrNoSuchProject)
}

// GetList ...
func (h *ProjectInfoHandler) GetList() ([]string, error) {
	res := []string{}
	for key := range h.projectList {
		res = append(res, key)
	}
	return res, nil
}

// Get ...
func (h *ProjectInfoHandler) Get(name string) (*model.ProjectInfo, error) {
	res, ok := h.projectList[name]
	if ok {
		return &res, nil
	}
	return nil, errors.Cause(model.ErrNoSuchProject)
}

// Update ...
func (h *ProjectInfoHandler) Update(ent *model.ProjectInfo) error {
	if _, ok := h.projectList[ent.Name]; ok {
		h.projectList[ent.Name] = *ent
		return nil
	}
	return errors.Cause(model.ErrNoSuchProject)
}
