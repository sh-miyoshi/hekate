package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// ProjectInfoHandler implement db.ProjectInfoHandler
type ProjectInfoHandler struct {
	projectList []*model.ProjectInfo
}

// NewProjectHandler ...
func NewProjectHandler() *ProjectInfoHandler {
	return &ProjectInfoHandler{}
}

// Add ...
func (h *ProjectInfoHandler) Add(ent *model.ProjectInfo) *errors.Error {
	h.projectList = append(h.projectList, ent)
	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(name string) *errors.Error {
	newList := []*model.ProjectInfo{}
	found := false
	for _, p := range h.projectList {
		if p.Name == name {
			found = true
		} else {
			newList = append(newList, p)
		}
	}

	if found {
		h.projectList = newList
		return nil
	}
	return errors.New("Internal Error", "No such project %s", name)
}

// GetList ...
func (h *ProjectInfoHandler) GetList(filter *model.ProjectFilter) ([]*model.ProjectInfo, *errors.Error) {
	res := []*model.ProjectInfo{}
	for _, prj := range h.projectList {
		res = append(res, prj)
	}

	if filter != nil {
		res = filterProjectList(res, filter)
	}

	return res, nil
}

// Update ...
func (h *ProjectInfoHandler) Update(ent *model.ProjectInfo) *errors.Error {
	for i, p := range h.projectList {
		if p.Name == ent.Name {
			h.projectList[i] = ent
			return nil
		}
	}
	return errors.New("Internal Error", "No such project %s", ent.Name)
}

func filterProjectList(data []*model.ProjectInfo, filter *model.ProjectFilter) []*model.ProjectInfo {
	if filter == nil {
		return data
	}
	res := []*model.ProjectInfo{}

	for _, prj := range data {
		if filter.Name != "" && prj.Name != filter.Name {
			// missmatch name
			continue
		}
		res = append(res, prj)
	}

	return res
}
