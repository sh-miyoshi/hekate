package memory

import (
	"sync"

	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// TransactionManager ...
type TransactionManager struct {
	mu sync.Mutex
}

// NewTransactionManager ...
func NewTransactionManager() *TransactionManager {
	return &TransactionManager{}
}

// Transaction ...
func (m *TransactionManager) Transaction(txFunc func() *errors.Error) *errors.Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return txFunc()
}
