package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"sync"
)

// AuthCodeHandler implement db.AuthCodeHandler
type AuthCodeHandler struct {
	authCodeList map[string]*model.AuthCode
	mu           sync.Mutex
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

// BeginTx ...
func (h *AuthCodeHandler) BeginTx() error {
	h.mu.Lock()
	return nil
}

// CommitTx ...
func (h *AuthCodeHandler) CommitTx() error {
	h.mu.Unlock()
	return nil
}

// AbortTx ...
func (h *AuthCodeHandler) AbortTx() error {
	h.mu.Unlock()
	return nil
}
