package model

import "github.com/sh-miyoshi/hekate/pkg/errors"

// TransactionManager ...
type TransactionManager interface {
	Transaction(txFunc func() *errors.Error) *errors.Error
}
