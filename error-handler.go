package land

import (
	"log"
	"reflect"
	"strings"
)

type errorHandler struct {
	land *land
}

func createErrorHandler(land *land) *errorHandler {
	return &errorHandler{
		land: land,
	}
}

func Recover(land Land) {
	createErrorHandler(land.getPtr()).recover()
}

func (e *errorHandler) createErrorMessage(err error, msg string, query string) {
	formatSlice := make([]string, 0)
	if err != nil {
		formatSlice = append(formatSlice, "[ERROR]: "+err.Error())
	}
	if len(msg) > 0 {
		formatSlice = append(formatSlice, "[MESSAGE]: "+msg)
	}
	if len(query) > 0 {
		formatSlice = append(formatSlice, "[QUERY]: "+query)
	}
	formatSlice = append(formatSlice, "----------")
	result := "\n" + strings.Join(formatSlice, "\n")
	if e.land.migration {
		log.Fatalln(result)
		return
	}
	log.Println(result)
}

func (e *errorHandler) recover() {
	if recovered := recover(); recovered != nil {
		if reflect.TypeOf(recovered) != reflect.TypeOf(Error{}) {
			return
		}
		err := recovered.(Error)
		e.createErrorMessage(err.Error, err.Message, err.Query)
	}
}
