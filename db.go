package land

import (
	"database/sql"
	"database/sql/driver"

	_ "github.com/lib/pq"
)

type db struct {
	config     Config
	connector  driver.Connector
	connection *sql.DB
}

func createConnection(config Config, connector driver.Connector) *db {
	d := &db{config: config, connector: connector}
	d.connection = sql.OpenDB(d.connector)
	return d
}
