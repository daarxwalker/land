package land

import (
	"fmt"
	"os"
	"strconv"
	"time"
	
	"github.com/dchest/uniuri"
)

type Migrator interface {
	Init()
	New()
	Up()
	Down()
}

type migrator struct {
	land              *land
	migrationsManager *migrationsManager
	errorHandler      *errorHandler
}

type landMigration struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

const (
	migrationsFolder   string = "migrations"
	migrationsMainFile        = "migrations.go"
)

const (
	migrationsEntityName  string = "land_migrations"
	migrationsEntityAlias string = "lmig"
)

const (
	mainMigrationsFileContent = `package main
	
import (
	"land/land"
)

var Migrations = land.Migrations()
		`
	newMigrationFileContent = `package main

import (
	"land"
)

func init() {
	Migrations.Add("%s").
		Up(func(l land.Land) {
		
		}).
		Down(func(l land.Land) {
		
		})
}
`
)

func createMigrator(land *land, migrationsManager *migrationsManager) Migrator {
	return &migrator{
		land:              land,
		migrationsManager: migrationsManager,
		errorHandler:      createErrorHandler(land),
	}
}

func (m *migrator) Init() {
	defer m.errorHandler.recover()
	fmt.Println("### Initializing...")
	m.createMigrationsEntity().CreateTable().IfNotExists().Exec()
	fmt.Println("### INIT SUCCESS!")
}

func (m *migrator) New() {
	defer m.errorHandler.recover()
	m.createMigration()
}

func (m *migrator) Up() {
	defer m.errorHandler.recover()
	dbMigrations := make([]landMigration, 0)
	lm := m.createMigrationsEntity()
	{
		q := lm.Select()
		q.Columns(Id, Name)
		q.Order().Asc(CreatedAt)
		q.GetResult(&dbMigrations)
	}
	for _, fileMigration := range m.migrationsManager.migrations {
		var exist bool
		for _, dbMigration := range dbMigrations {
			if fileMigration.id == dbMigration.Name {
				exist = true
			}
		}
		if exist {
			continue
		}
		fmt.Println("### Migrating: " + fileMigration.id)
		tx := m.land.Transaction()
		tx.Begin()
		fileMigration.up(m.land)
		tx.Commit()
		lm.Insert().SetValues(landMigration{Name: fileMigration.id}).Exec()
		fmt.Println("### MIGRATION SUCCESS: " + fileMigration.id)
	}
}

func (m *migrator) Down() {
	defer m.errorHandler.recover()
	var lastMigration landMigration
	{
		lm := m.createMigrationsEntity()
		q := lm.Select()
		q.Columns(Id, Name)
		q.Order().Desc(CreatedAt)
		q.Single().GetResult(&lastMigration)
		if lm.IsError() {
			err := lm.Error()
			m.errorHandler.createErrorMessage(err.Error, err.Query, err.Message)
			return
		}
	}
	migration := m.getMigrationWithId(lastMigration.Name)
	if migration == nil {
		return
	}
	fmt.Println("### Rollbacking: " + migration.id)
	tx := m.land.Transaction()
	tx.Begin()
	migration.down(m.land)
	tx.Commit()
	{
		lm := m.createMigrationsEntity()
		q := lm.Delete()
		q.Where().Column(Name).Equal(migration.id)
		q.Exec()
		if lm.IsError() {
			err := lm.Error()
			m.errorHandler.createErrorMessage(err.Error, err.Query, err.Message)
			return
		}
	}
	fmt.Println("### ROLLBACK SUCCESS: " + migration.id)
}

func (m *migrator) getRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func (m *migrator) getDir() string {
	root := m.getRoot()
	if len(root) == 0 {
		return ""
	}
	result := root + "/" + migrationsFolder
	if len(m.migrationsManager.dbname) > 0 {
		result += "/" + m.migrationsManager.dbname
	}
	return result
}

func (m *migrator) getMigrationsMainFileDir() string {
	dir := m.getDir()
	if len(dir) == 0 {
		return ""
	}
	return dir + "/" + migrationsMainFile
}

func (m *migrator) createMigrationsEntity() Entity {
	return m.land.CreateEntity(migrationsEntityName).
		SetAlias(migrationsEntityAlias).
		SetColumn(Name, Text, ColOpts{NotNull: true}).
		SetCreatedAt().
		SetUpdatedAt()
}

func (m *migrator) getMigrationWithId(id string) *migration {
	for _, item := range m.migrationsManager.migrations {
		if item.id == id {
			return item
		}
	}
	return nil
}

func (m *migrator) createMigration() Migrator {
	dir := m.getDir()
	if len(dir) == 0 {
		return m
	}
	id := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + uniuri.New()
	filedir := dir + "/" + id + ".go"
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		file, err := os.Create(filedir)
		if err != nil {
			m.errorHandler.createErrorMessage(err, "create new migration file failed", "")
		}
		_, err = file.WriteString(fmt.Sprintf(newMigrationFileContent, id))
		if err != nil {
			m.errorHandler.createErrorMessage(err, "write init content to new migration failed", "")
		}
	}
	return m
}

func (m *migrator) verifyMigrationsDir() {
	dir := m.getDir()
	if len(dir) == 0 {
		return
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			m.errorHandler.createErrorMessage(err, "create migrations folder failed", "")
		}
	}
}

func (m *migrator) verifyMainMigrationsFile() {
	dir := m.getMigrationsMainFileDir()
	if len(dir) == 0 {
		return
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		file, err := os.Create(dir)
		if err != nil {
			m.errorHandler.createErrorMessage(err, "write init content to new migration failed", "")
		}
		_, err = file.WriteString(mainMigrationsFileContent)
		if err != nil {
			m.errorHandler.createErrorMessage(err, "write init content to new migration failed", "")
		}
	}
}
