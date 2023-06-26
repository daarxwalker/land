package land

type MigrationsManager interface {
}

type migrationsManager struct {
}

func Migrations() MigrationsManager {
	return &migrationsManager{}
}
