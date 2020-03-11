package memory

import (
	"sync"
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
func (m *TransactionManager) Transaction(txFunc func() error) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return txFunc()
}
