package land

type MigrationsManager interface {
}

type migrationsManager struct {
}

const (
	migrationsFolder = "migrations"
)

const (
	baseMigrationsFileContent = `package main
	
import (
	"land/land"
)

var Migrations = land.NewMigrations()
		`
	newMigrationFileContent = `package main

import (
	"land"
)

func init() {
	Migrations.Add().
		Up(func(orm land.ORM) {
		
		}).
		Down(func(orm land.ORM) {
		
		})
}
`
)

func Migrations() MigrationsManager {
	return &migrationsManager{}
}
