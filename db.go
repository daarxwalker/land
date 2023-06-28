package land

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type db struct {
	config     Config
	connector  *connector
	connection *sql.DB
}

func createConnection(config Config, connector *connector) *db {
	d := &db{config: config, connector: connector}
	connection, err := sql.Open(connector.dbtype, connector.createConnectionString())
	if err != nil {
		log.Fatalln(err)
	}
	d.connection = connection
	return d
}
