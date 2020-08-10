package mongo

import (
	"context"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// TransactionManager ...
type TransactionManager struct {
	session  mongo.Session
	dbClient *mongo.Client
}

// NewTransactionManager ...
func NewTransactionManager(dbClient *mongo.Client) *TransactionManager {
	return &TransactionManager{
		dbClient: dbClient,
	}
}

// Transaction ...
func (m *TransactionManager) Transaction(txFunc func() *errors.Error) *errors.Error {
	if err := m.beginTx(); err != nil {
		return errors.New("DB failed", "Begin transaction failed: %v", err)
	}

	if err := txFunc(); err != nil {
		m.abortTx()
		return err
	}
	return m.commitTx()
}

func (m *TransactionManager) beginTx() *errors.Error {
	var err error
	m.session, err = m.dbClient.StartSession()
	if err != nil {
		return errors.New("DB failed", "Failed to start mongo session: %v", err)
	}
	err = m.session.StartTransaction()
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
		defer cancel()
		m.session.EndSession(ctx)
		return errors.New("DB failed", "Failed to start mongo transaction: %v", err)
	}
	return nil
}

func (m *TransactionManager) abortTx() *errors.Error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := m.session.AbortTransaction(ctx)
	m.session.EndSession(ctx)
	if err != nil {
		return errors.New("DB failed", "Failed abort transaction: %v", err)
	}
	return nil
}

func (m *TransactionManager) commitTx() *errors.Error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := m.session.CommitTransaction(ctx)
	m.session.EndSession(ctx)
	if err != nil {
		return errors.New("DB failed", "Failed commit transaction: %v", err)
	}
	return nil
}
