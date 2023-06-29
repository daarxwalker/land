
# Land - simple & elegant Go ORM
**ALPHA version**
> **Bugs are expected!** **Opened issues are very welcome!** \
> Mainly created for my personal purpose. Used with go workspaces. \
> Right now, only works with **Postgres**.
> All features will be added over time!

## Why another Go ORM?
>I just don't like Gorm and Bun way. Sometimes too much simple for complicated queries, sometimes these ORMs lack of features, which I need, so I had to create wrapper and that's not good, so I've created a new way of writing db queries in Go.

## Get started
- [Connection](#connection)
- [Entity](#entity)
- [Migrations](#migrations)
- [Transactions](#transactions)

## Query
- [Select](#select-query)
- [Insert](#insert-query)
- [Update](#update-query)
- [Delete](#delete-query)
- [Create table](#create-table-query)
- [Alter table](#alter-table-query)
- [Drop table](#drop-table-query)
- [Join](#join-query)
- [Where](#where-query)
- [Group](#group-query)
- [Order](#order-query)

## Get started

### Connection
Use *land.Connect()* to create connection with your database.\
You can verify it with *.Ping()*.
```go
l := land.New(  
    land.Config{  
        Production: false,  
        Log: true,   
    },  
    land.Connect().  
        Postgres().  
        Host("localhost").  
        Port(5432).  
        User("land").  
        Password("land").  
        Dbname("land"),  
)
fmt.Println(l.Ping())
```
### Entity
Entities are logic blocks, which hold table structure.\
You need Land instance to create and use entity.\
Queries are created with entities.

### Entity syntax example
```go
package user_entity  

import (
	"land"
	r "project/entity/role_entity"
)

const (  
    EntityName    = "users"  
    EntityAlias   = "u"  
    RoleId        = "role_id"  
    Active        = "active"  
    Name          = "name"  
    Lastname      = "lastname"  
)  

var (
	Columns = []string{land.Id, Active, Name, Lastname}
)

func User(l land.Land) land.Entity {  
    return l.CreateEntity(EntityName).  
        SetAlias(EntityAlias).
        SetColumn(
            RoleId, 
            land.Int,
            land.ColOpts{
                Reference: land.EntityReference{ Entity: r.Role(l), Column: land.Id },
            },
        ).
        SetColumn(Active, land.Varchar, land.ColOpts{NotNull: true, Default: false}).  
        SetColumn(Name, land.Varchar, land.ColOpts{Limit: 255, NotNull: true}).  
        SetColumn(Lastname, land.Varchar, land.ColOpts{Limit: 255, NotNull: true}).  
        SetFulltext().  
        SetCreatedAt().  
        SetUpdatedAt()  
}
```

### Migrations
Migrations folder has to be in project root!

#### Migrations folder structure example
```
- []migrations
    - 1687967915906263000_G037Ozocz4NrkgwH.go
    - go.mod
    - main.go
```
#### Migrations main.go
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
#### Migrations commands example
>Commands are based on config mentioned above.
```
- Init migrations: go run ./migrations/*.go --init
- New migration: go run ./migrations/*.go --new
- Up migrations: go run ./migrations/*.go --up
- Down migration: go run ./migrations/*.go --down
```

### Transactions
```go
l := postgres.New()
tx := l.Transaction()
tx.Begin()
if someError != nil {
  tx.Rollback()
}
tx.Commit()
```


## Queries
For our queries, I prefer to use import alias:
```go
import (
    "land"
    u "project/entity/user_entity"
    r "project/entity/role_entity"
)
```

### Select query
```go
func GetAll(l land.Land) []user_model.User {  
    result := make([]user_model.User, 0)
    q := u.User(l).Select()  
    q.Columns(u.Columns...)  
    q.GetResult(&result)  
    return result
}
```

### Insert query
```go
func CreateOne(l land.Land, data user_model.User) user_model.User {  
    q := u.User(l).Insert()
    q.SetData(data)  
    q.SetVectors(data.Name, data.Lastname)
    q.Return(land.Id)  
    q.GetResult(&data)  
    return data
}
```

### Update query
```go
func UpdateOne(l land.Land, data user_model.User) user_model.User {  
    var result user_model.User
    q := u.User(l).Update()
    q.SetData(data)  
    q.SetVectors(data.Name, data.Lastname)
    q.Return(land.Id)  
    q.GetResult(&result)  
    return result
}
```

### Delete query
```go
func RemoveOne(l land.Land, id int) {  
    q := u.User(l).Delete()
    q.Where().Column(land.Id).Equal(id)  
    q.Exec() 
}
```

### Create table query
```go
func CreateTable(l land.Land) {
  u.User(l).CreateTable().Exec()
}
```

### Alter table query
```go
func AlterTable(l land.Land) {
  u.User(l).AlterTable().IfExists().
    AddColumn("middle_name", Varchar, ColOpts{Limit: 255, NotNull: true, Unique: true}).
    RenameColumn("name", "custom_name").
    DropColumn("custom_name").
    Exec()
}
```

### Drop table query
```go
func DropTable(l land.Land) {
  u.User(l).DropTable().Exec()
}
```

### Join query
```go
func GetAllWithJoinedRole(l land.Land, id int) user_model.User {
    result := make([]user_model.User, 0)
    q := u.User(l).Select()  
    q.Columns(u.Columns...)  
    q.Join().Column(u.RoleId).On(r.Role(l))
    q.GetResult(&result)  
    return result
}
```

### Where query
```go
func GetOne(l land.Land, id int) user_model.User {  
    var result user_model.User
    q := u.User(l).Select()  
    q.Columns(u.Columns...)  
    q.Where().Column(land.Id).Equal(id)
    q.Single()
    q.GetResult(&result)  
    return result
}
```

### Group query
```go
func GetAllLastnamesCount(l land.Land, id int) user_model.User {
    result := make([]user_model.UserWithLastnameCount, 0)  
    q := u.User(l).Select()  
    q.Column(u.Lastname)
    q.Column(u.Lastname).Count().Alias("lastname_count")
    q.Group().Columns(u.Lastname)
    q.GetResult(&result)  
    return result
}
```

### Order query
```go
func GetAllOrderDescById(l land.Land) []user_model.User {
    result := make([]user_model.User, 0)
    q := u.User(l).Select()  
    q.Columns(u.Columns...)  
    q.Order().Desc(land.Id)
    q.GetResult(&result)  
    return result
}
```
