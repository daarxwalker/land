package land

import (
	"errors"

	"util/dd"
)

type ErrorManager interface {
	IsError() bool
	Errors() []Error
	Error() Error
}

type errorManager struct {
	errors []Error
}

func createErrorManager() *errorManager {
	return &errorManager{
		errors: make([]Error, 0),
	}
}

func (m *errorManager) IsError() bool {
	return len(m.errors) > 0
}

func (m *errorManager) Errors() []Error {
	return m.errors
}

func (m *errorManager) Error() Error {
	if m.IsError() {
		return m.errors[len(m.errors)-1]
	}
	return Error{}
}

func (m *errorManager) check(err error, query string) {
	if err == nil {
		return
	}
	e := Error{Error: err, Query: query}
	dd.Print(e)
	m.errors = append(m.errors, e)
	panic(e)
}

func (m *errorManager) throw(message, query string) {
	e := Error{Error: errors.New(message), Query: query}
	m.errors = append(m.errors, e)
	panic(e)
}
