package land

import (
	"errors"
)

type Land interface {
	CreateEntity(name string) Entity
	Migrator(migrationsManager MigrationsManager) Migrator
	Ping() error
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
