package land

import (
	"database/sql/driver"
	"errors"
)

type Land interface {
	EntityManager() EntityManager
	Ping() error
}

type land struct {
	db            *db
	entityManager *entityManager
	config        Config
}

func New(config Config, connector driver.Connector) Land {
	l := &land{
		config: config,
	}
	l.entityManager = createEntityManager(l)
	if connector != nil {
		l.db = createConnection(config, connector)
	}
	return l
}

func (l *land) EntityManager() EntityManager {
	return l.entityManager
}

func (l *land) Ping() error {
	if l.db != nil {
		return errors.New("land orm isn't connected to the database")
	}
	return l.db.connection.Ping()
}
