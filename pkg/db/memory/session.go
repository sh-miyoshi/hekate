package memory

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// SessionHandler implement db.SessionHandler
type SessionHandler struct {
	// sessionList[sessionID] = Session
	sessionList []*model.Session
}

// NewSessionHandler ...
func NewSessionHandler() *SessionHandler {
	res := &SessionHandler{}
	return res
}

// Add ...
func (h *SessionHandler) Add(projectName string, ent *model.Session) *errors.Error {
	h.sessionList = append(h.sessionList, ent)
	return nil
}

// Delete ...
func (h *SessionHandler) Delete(projectName string, filter *model.SessionFilter) *errors.Error {
	newList := missMatchFilterSessionList(h.sessionList, projectName, filter)

	h.sessionList = newList
	return nil
}

// DeleteAll ...
func (h *SessionHandler) DeleteAll(projectName string) *errors.Error {
	newList := []*model.Session{}
	for _, s := range h.sessionList {
		if s.ProjectName != projectName {
			newList = append(newList, s)
		}
	}
	h.sessionList = newList
	return nil
}

// GetList ...
func (h *SessionHandler) GetList(projectName string, filter *model.SessionFilter) ([]*model.Session, *errors.Error) {
	res := []*model.Session{}

	for _, s := range h.sessionList {
		if s.ProjectName == projectName {
			res = append(res, s)
		}
	}

	if filter != nil {
		res = matchFilterSessionList(res, projectName, filter)
	}

	return res, nil
}

// Cleanup ...
func (h *SessionHandler) Cleanup(now time.Time) *errors.Error {
	newList := []*model.Session{}
	for _, s := range h.sessionList {
		expire := s.CreatedAt.Add(time.Second * time.Duration(s.ExpiresIn))
		if now.Before(expire) {
			newList = append(newList, s)
		}
	}

	h.sessionList = newList
	return nil
}

func matchFilterSessionList(data []*model.Session, projectName string, filter *model.SessionFilter) []*model.Session {
	if filter == nil {
		return data
	}
	res := []*model.Session{}

	for _, s := range data {
		if projectName == s.ProjectName {
			if filter.SessionID != "" && s.SessionID != filter.SessionID {
				// missmatch session id
				continue
			}
			if filter.UserID != "" && s.UserID != filter.UserID {
				// missmatch user id
				continue
			}
		}
		res = append(res, s)
	}

	return res
}

func missMatchFilterSessionList(data []*model.Session, projectName string, filter *model.SessionFilter) []*model.Session {
	if filter == nil {
		return data
	}
	res := []*model.Session{}

	for _, s := range data {
		if projectName == s.ProjectName {
			if filter.SessionID != "" && s.SessionID == filter.SessionID {
				// matched to session id
				continue
			}
			if filter.UserID != "" && s.UserID == filter.UserID {
				// matched to user id
				continue
			}
		}
		res = append(res, s)
	}

	return res
}
