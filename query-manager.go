package land

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type queryManager struct {
	entity    *entity
	context   context.Context
	queryType string
	query     string
	dest      any
	rows      [][]sql.RawBytes
	errors    []error
	duration  time.Duration
}

func createQueryManager(entity *entity, context context.Context) *queryManager {
	return &queryManager{
		entity:  entity,
		context: context,
		rows:    make([][]sql.RawBytes, 0),
	}
}

func (m *queryManager) setQuery(query string) *queryManager {
	m.query = query
	return m
}

func (m *queryManager) setQueryType(queryType string) *queryManager {
	m.queryType = queryType
	return m
}

func (m *queryManager) setDest(dest any) *queryManager {
	m.dest = dest
	return m
}

func (m *queryManager) getResult() {
	m.scan()
	m.log()
}

func (m *queryManager) scan() {
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

func (m *queryManager) addRow(row []any) {
	result := make([]sql.RawBytes, 0)
	for _, c := range row {
		result = append(result, *c.(*sql.RawBytes))
	}
	m.rows = append(m.rows, result)
}

func (m *queryManager) exec() {
	_, err := m.connection().Exec(m.query)
	if err != nil {
		panic(Error{error: err, query: m.query})
	}
}

func (m *queryManager) connection() *sql.DB {
	return m.entity.land.db.connection
}

func (m *queryManager) log() {
	if !m.entity.land.config.Log {
		return
	}
	fmt.Println(fmt.Sprintf("%s in %v: %s\n", strings.ToUpper(m.queryType), m.duration, m.query))
}
