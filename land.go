package land

import (
	"errors"
)

type Land interface {
	CreateEntity(name string) Entity
	Migrator(migrationsManager MigrationsManager) Migrator
	Ping() error
	Begin() error
	Commit() error
	Rollback() error
}

type land struct {
	db        *db
	entities  []*entity
	config    Config
	migration bool
}

func New(config Config, connector Connector) Land {
	l := &land{
		config: config,
	}
	if connector != nil {
		l.db = createConnection(config, connector.getPtr())
	}
	return l
}

func (l *land) CreateEntity(name string) Entity {
	e := createEntity(l, name)
	l.entities = append(l.entities, e)
	return e
}

func (l *land) Begin() error {
	_, err := l.db.connection.Exec("BEGIN;")
	return err
}

func (l *land) Commit() error {
	_, err := l.db.connection.Exec("COMMIT;")
	return err
}

func (l *land) Rollback() error {
	_, err := l.db.connection.Exec("ROLLBACK;")
	return err
}

func (l *land) Migrator(migrationsManager MigrationsManager) Migrator {
	l.migration = true
	return createMigrator(l, migrationsManager.getPtr())
}

func (l *land) Ping() error {
	if l.db == nil {
		return errors.New("land failed connect to the database")
	}
	return l.db.connection.Ping()
}

func (l *land) Transaction() Transaction {
	return createTransactionManager(l)
}
