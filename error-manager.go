package land

type errorManager struct {
	errors []Error
}

func createErrorManager() *errorManager {
	return &errorManager{}
}

func (m *errorManager) check(err error, query string) {

}
