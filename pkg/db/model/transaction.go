package model

// TransactionManager ...
type TransactionManager interface {
	Transaction(txFunc func() error) error
}
