package land

import (
	"fmt"
	"log"
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
	migrationsEntityName string = "land_migrations"
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
	Migrations.Add(%s).
		Up(func(orm land.ORM) {
		
		}).
		Down(func(orm land.ORM) {
		
		})
}
`
)

func createMigrator(land *land, migrationsManager *migrationsManager) Migrator {
	return &migrator{
		land:              land,
		migrationsManager: migrationsManager,
	}
}

func (m *migrator) Init() {
	fmt.Println("### Initializing...")
	m.createMigrationsEntity().CreateTable().IfNotExists().Exec()
	fmt.Println("### INIT SUCCESS!")
}

func (m *migrator) New() {
	m.createMigration()
}

func (m *migrator) Up() {
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
		lm.Begin()
		fileMigration.up(m.land)
		lm.Commit()
		lm.Insert().Values(landMigration{Name: fileMigration.id}).Exec()
		fmt.Println("### MIGRATION SUCCESS: " + fileMigration.id)
	}
}

func (m *migrator) Down() {
	var lastMigration landMigration
	lm := m.createMigrationsEntity()
	{
		q := lm.Select()
		q.Columns(Id, Name)
		q.Order().Desc(CreatedAt)
		q.Single().GetResult(&lastMigration)
	}
	migration := m.getMigrationWithId(lastMigration.Name)
	if migration == nil {
		return
	}
	fmt.Println("### Rollbacking: " + migration.id)
	lm.Begin()
	migration.down(m.land)
	lm.Commit()
	{
		q := lm.Delete()
		q.Where().Column(Name).Equal(migration.id)
		q.Exec()
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
	return root + "/" + migrationsFolder
}

func (m *migrator) getMigrationsMainFileDir() string {
	dir := m.getDir()
	if len(dir) == 0 {
		return ""
	}
	return dir + "/" + migrationsMainFile
}

func (m *migrator) check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func (m *migrator) createMigrationsEntity() Entity {
	return m.land.CreateEntity(migrationsEntityName).
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
		m.check(err)
		_, err = file.WriteString(fmt.Sprintf(newMigrationFileContent, id))
		m.check(err)
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
			log.Fatal(err)
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
		m.check(err)
		_, err = file.WriteString(mainMigrationsFileContent)
		m.check(err)
	}
}
