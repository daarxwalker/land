
# Land - simple & elegant Go ORM
**ALPHA version**\
Land was created for my personal purpose.\
Primarly using Go workspaces.\
It's specialized to do DB operations in very simple & elegant way.\
Right now, only works with Postgres.\
This project was rewritten of first version, so bugs are expected!\
All features will be added over time.\
Issues are very welcome.


## Get started
Use *land.Connector* with func *.Connect()* to connect to your database.
```go
l := land.New(  
    land.Config{  
        Production: false,  
        Log: true,  
        DatabaseType: land.Postgres,  
    },  
    land.Connect().  
        Postgres().  
        Host("localhost").  
        Port(5432).  
        User("land").  
        Password("land").  
        Dbname("land"),  
)
```
## Entity
Entities are blocks, which hold table structure.\
You need Land instance to create and use entity.\
Entities have access to queries.

### Entity syntax example
```go
package user_entity  

import "land"  

const (  
    EntityName 	= "users"  
    EntityAlias = "u"  
    Active 		= "active"  
    Name 		= "name"  
    Lastname 	= "lastname"  
)  

func User(l land.Land) land.Entity {  
    return l.CreateEntity(EntityName).  
        SetAlias(EntityAlias).  
        SetColumn(Active, land.Varchar, land.ColOpts{NotNull: true, Default: false}).  
        SetColumn(Name, land.Varchar, land.ColOpts{Limit: 255, NotNull: true}).  
        SetColumn(Lastname, land.Varchar, land.ColOpts{Limit: 255, NotNull: true}).  
        SetFulltext().  
        SetCreatedAt().  
        SetUpdatedAt()  
}
```
## Queries examples
```go
package user_repository

import (
    "land"
    u "project/entity/user_entity"
    "project/model/user_model"
)

type UserRepository interface {
    GetAll() user_model.User
    CreateOne(data user_model.User) user_model.User
    CreateOne(data user_model.User) user_model.User
    RemoveOne(id int)
}

type userRepository struct {
    land land.Land
}

func (r *userRepository) GetAll() user_model.User {  
    var result user_model.User
    q := u.User(r.land).Select()  
    q.Columns(land.Id, u.Active, u.Name, u.Lastname)  
    q.GetResult(&result)  
    return result
}

func (r *userRepository) CreateOne(data user_model.User) user_model.User {  
    q := u.User(r.land).Insert()
    q.SetData(data)  
    q.SetVectors(data.Name, data.Lastname)
    q.Return(land.Id)  
    q.GetResult(&data)  
    return data
}

func (r *userRepository) UpdateOne(data user_model.User) user_model.User {  
    var result user_model.User
    q := u.User(r.land).Update()
    q.SetData(data)  
    q.SetVectors(data.Name, data.Lastname)
    q.Return(land.Id)  
    q.GetResult(&result)  
    return result
}

func (r *userRepository) RemoveOne(id int) {  
    q := u.User(r.land).Delete()
    q.Where().Column(land.Id).Equal(id)  
    q.Exec() 
}
```
## Migrations
Migrations folder have to be in project root!

### Migrations folder structure example
```
- []migrations
    - 1687967915906263000_G037Ozocz4NrkgwH.go
    - go.mod
    - main.go
```
### Migrations main.go
```go
package main

import (
    "flag"
    "land"
    "project/infrastructure/postgres"
)

var Migrations = land.Migrations()  
  
func main() {  
    l := postgres.New()  
    initMigrations := flag.Bool("init", false, "Initialize migrations")  
    newMigration := flag.Bool("new", false, "New migration")  
    upMigrations := flag.Bool("up", false, "Up migrations")  
    downMigration := flag.Bool("down", false, "Down migration")  
    flag.Parse()  
    if *initMigrations {  
        l.Migrator(Migrations).Init()  
        return  
    }  
    if *newMigration {  
        l.Migrator(Migrations).New()  
        return  
    }  
    if *upMigrations {  
        l.Migrator(Migrations).Up()  
        return  
    }  
    if *downMigration {  
        l.Migrator(Migrations).Down()  
        return  
    }  
}
```
### Migrations commands example
Commands are based on config mentioned above.
```
- Init migrations: go run ./migrations/*.go --init
- New migration: go run ./migrations/*.go --new
- Up migrations: go run ./migrations/*.go --up
- Down migration: go run ./migrations/*.go --down
```