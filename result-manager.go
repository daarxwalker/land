package land

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type resultManager struct {
	entity    *entity
	context   context.Context
	queryType string
	query     string
	dest      any
	rows      [][]sql.RawBytes
	errors    []error
	duration  time.Duration
}

func createResultManager(entity *entity, context context.Context) *resultManager {
	return &resultManager{
		entity:  entity,
		context: context,
		rows:    make([][]sql.RawBytes, 0),
	}
}

func (m *resultManager) setQuery(query string) *resultManager {
	m.query = query
	return m
}

func (m *resultManager) setQueryType(queryType string) *resultManager {
	m.queryType = queryType
	return m
}

func (m *resultManager) setDest(dest any) *resultManager {
	m.dest = dest
	return m
}

func (m *resultManager) getResult() {
	m.scan()
	m.log()
}

func (m *resultManager) scan() {
	start := time.Now()
	rows, err := m.connection().QueryContext(m.context, m.query)
	m.duration = time.Now().Sub(start)
	defer func() {
		if err := rows.Close(); err != nil {
			panic(Error{error: err, query: m.query})
		}
	}()
	if err != nil {
		panic(Error{error: err, query: m.query})
	}
	columnsTypes, err := rows.ColumnTypes()
	if err != nil {
		panic(Error{error: err, query: m.query})
	}
	for rows.Next() {
		row := make([]any, len(columnsTypes))
		for i := range columnsTypes {
			row[i] = new([]sql.RawBytes)
		}
		err = rows.Scan(row...)
		m.addRow(row)
	}
}

func (m *resultManager) addRow(row []any) {
	result := make([]sql.RawBytes, 0)
	for _, c := range row {
		result = append(result, *c.(*sql.RawBytes))
	}
	m.rows = append(m.rows, result)
}

func (m *resultManager) exec() {
	_, err := m.connection().Exec(m.query)
	if err != nil {
		panic(Error{error: err, query: m.query})
	}
}

func (m *resultManager) connection() *sql.DB {
	return m.entity.entityManager.land.db.connection
}

func (m *resultManager) log() {
	if !m.entity.entityManager.land.config.Log {
		return
	}
	fmt.Println(fmt.Sprintf("%s in %v: %s\n", strings.ToUpper(m.queryType), m.duration, m.query))
}
