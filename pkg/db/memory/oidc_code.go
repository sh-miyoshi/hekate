package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
)

// AuthCodeHandler implement db.AuthCodeHandler
type AuthCodeHandler struct {
	authCodeList map[string]*model.AuthCode
}

// NewAuthCodeHandler ...
func NewAuthCodeHandler() *AuthCodeHandler {
	res := &AuthCodeHandler{
		authCodeList: make(map[string]*model.AuthCode),
	}
	return res
}

// Add ...
func (h *AuthCodeHandler) Add(code *model.AuthCode) error {
	h.authCodeList[code.CodeID] = code
	return nil
}

// Get ...
func (h *AuthCodeHandler) Get(codeID string) (*model.AuthCode, error) {
	if _, exists := h.authCodeList[codeID]; !exists {
		return nil, model.ErrNoSuchCode
	}
	return h.authCodeList[codeID], nil
}

// Delete ...
func (h *AuthCodeHandler) Delete(codeID string) error {
	if _, exists := h.authCodeList[codeID]; exists {
		delete(h.authCodeList, codeID)
		return nil
	}

	return model.ErrNoSuchCode
}

// DeleteAll ...
func (h *AuthCodeHandler) DeleteAll(userID string) error {
	newList := make(map[string]*model.AuthCode)
	for _, a := range h.authCodeList {
		if a.UserID != userID {
			newList[a.CodeID] = a
		}
	}
	h.authCodeList = newList
	return nil
}
