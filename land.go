package land

import (
	"database/sql/driver"
	"errors"
)

type Land interface {
	CreateEntity(name string) Entity
	Migrator(migrationsManager MigrationsManager) Migrator
	Ping() error
}

type land struct {
	db       *db
	entities []*entity
	config   Config
}

func New(config Config, connector driver.Connector) Land {
	l := &land{
		config: config,
	}
	if connector != nil {
		l.db = createConnection(config, connector)
	}
	return l
}

func (l *land) CreateEntity(name string) Entity {
	e := createEntity(l, name)
	l.entities = append(l.entities, e)
	return e
}

func (l *land) Migrator(migrationsManager MigrationsManager) Migrator {
	return createMigrator(l, migrationsManager.getPtr())
}

func (l *land) Ping() error {
	if l.db != nil {
		return errors.New("land orm isn't connected to the database")
	}
	return l.db.connection.Ping()
}
