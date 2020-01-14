package memory

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
)

// AuthCodeHandler implement db.AuthCodeHandler
type AuthCodeHandler struct {
	authCodeList map[string]*model.AuthCode
}

// NewAuthCodeHandler ...
func NewAuthCodeHandler() (*AuthCodeHandler, error) {
	res := &AuthCodeHandler{
		authCodeList: make(map[string]*model.AuthCode),
	}
	return res, nil
}

// New ...
func (h *AuthCodeHandler) New(code *model.AuthCode) error {
	h.authCodeList[code.CodeID] = code
	return nil
}

// Get ...
func (h *AuthCodeHandler) Get(codeID string) (*model.AuthCode, error) {
	if _, exists := h.authCodeList[codeID]; !exists {
		return nil, errors.Cause(model.ErrNoSuchCode)
	}
	return h.authCodeList[codeID], nil
}

// Delete ...
func (h *AuthCodeHandler) Delete(codeID string) error {
	if _, exists := h.authCodeList[codeID]; exists {
		delete(h.authCodeList, codeID)
		return nil
	}

	return errors.Cause(model.ErrNoSuchCode)
}
