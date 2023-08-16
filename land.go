package land

import (
	"errors"
	"fmt"
	
	"github.com/iancoleman/strcase"
)

type Land interface {
	CreateEntity(name string) Entity
	Migrator(migrationsManager MigrationsManager) Migrator
	Ping() error
	Begin() error
	Commit() error
	Rollback() error
	Query(query string, args ...any) ([]map[string]any, error)
	FixSequence(table string) error
	Reset(table string) error
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

func (l *land) Query(query string, args ...any) ([]map[string]any, error) {
	result := make([]map[string]any, 0)
	rows, err := l.db.connection.Query(query, args...)
	if err != nil {
		return result, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return result, err
	}
	for rows.Next() {
		rowCols := make([]any, len(cols))
		rowColsPtrs := make([]any, len(cols))
		for i, _ := range rowCols {
			rowColsPtrs[i] = &rowCols[i]
		}
		err = rows.Scan(rowColsPtrs...)
		if err != nil {
			return result, err
		}
		rowResult := make(map[string]any)
		for i, columnName := range cols {
			rowResult[columnName] = *(rowColsPtrs[i].(*any))
		}
		result = append(result, rowResult)
	}
	return result, rows.Close()
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

func (l *land) FixSequence(table string) error {
	table = strcase.ToSnake(table)
	_, err := l.db.connection.Exec(fmt.Sprintf("SELECT setval('%[1]s_id_seq', (SELECT MAX(id) FROM %[1]s));", table))
	return err
}

func (l *land) Reset(table string) error {
	_, err := l.db.connection.Exec(
		fmt.Sprintf(
			`TRUNCATE TABLE %[1]s RESTART IDENTITY CASCADE; ALTER SEQUENCE %[1]s_id_seq RESTART WITH 1
;`, table,
		),
	)
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
