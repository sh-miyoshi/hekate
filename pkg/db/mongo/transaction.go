package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
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
func (m *TransactionManager) Transaction(txFunc func() error) error {
	if err := m.beginTx(); err != nil {
		return errors.Wrap(err, "Begin transaction failed")
	}

	if err := txFunc(); err != nil {
		m.abortTx()
		return err
	}
	return m.commitTx()
}

func (m *TransactionManager) beginTx() error {
	var err error
	m.session, err = m.dbClient.StartSession()
	if err != nil {
		return err
	}
	err = m.session.StartTransaction()
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
		defer cancel()
		m.session.EndSession(ctx)
		return err
	}
	return nil
}

func (m *TransactionManager) abortTx() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := m.session.AbortTransaction(ctx)
	m.session.EndSession(ctx)
	return err
}

func (m *TransactionManager) commitTx() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := m.session.CommitTransaction(ctx)
	m.session.EndSession(ctx)
	return err
}
