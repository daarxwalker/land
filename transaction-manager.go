package land

import (
	"database/sql"
)

type Transaction interface {
	ErrorManager
	Begin()
	Rollback()
	Commit()
}

type transactionManager struct {
	*errorManager
	land         *land
	errorHandler *errorHandler
}

func createTransactionManager(land *land) *transactionManager {
	return &transactionManager{
		errorManager: createErrorManager(),
		land:         land,
		errorHandler: createErrorHandler(land),
	}
}

func (m *transactionManager) Begin() {
	defer m.errorHandler.recover()
	query := "BEGIN;"
	_, err := m.connection().Exec(query)
	m.check(err, query)
}

func (m *transactionManager) Rollback() {
	defer m.errorHandler.recover()
	query := "ROLLBACK;"
	_, err := m.connection().Exec(query)
	m.check(err, query)
}

func (m *transactionManager) Commit() {
	defer m.errorHandler.recover()
	query := "COMMIT;"
	_, err := m.connection().Exec(query)
	m.check(err, query)
}

func (m *transactionManager) connection() *sql.DB {
	return m.land.db.connection
}
