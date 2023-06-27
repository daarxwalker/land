package land

type MigrationsManager interface {
	Add(id string) Migration

	getPtr() *migrationsManager
}

type migrationsManager struct {
	migrations []*migration
}

func Migrations() MigrationsManager {
	return &migrationsManager{}
}

func (m *migrationsManager) Add(id string) Migration {
	migration := &migration{id: id}
	m.migrations = append(m.migrations, migration)
	return migration
}

func (m *migrationsManager) getPtr() *migrationsManager {
	return m
}
