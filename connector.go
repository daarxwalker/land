package land

import (
	"fmt"
	"strings"
)

type Connector interface {
	Postgres() Connector
	Host(host string) Connector
	Port(port int) Connector
	User(user string) Connector
	Dbname(dbname string) Connector
	Password(password string) Connector
	SSL(sslmode string) Connector

	getPtr() *connector
}

type connector struct {
	dbtype   string
	host     string
	port     int
	user     string
	dbname   string
	password string
	sslmode  string
}

func Connect() Connector {
	return &connector{
		sslmode: "disable",
	}
}

const (
	Postgres string = "postgres"
)

func (c *connector) Postgres() Connector {
	c.dbtype = Postgres
	return c
}

func (c *connector) Host(host string) Connector {
	c.host = host
	return c
}

func (c *connector) Port(port int) Connector {
	c.port = port
	return c
}

func (c *connector) User(user string) Connector {
	c.user = user
	return c
}

func (c *connector) Dbname(dbname string) Connector {
	c.dbname = dbname
	return c
}

func (c *connector) Password(password string) Connector {
	c.password = password
	return c
}

func (c *connector) SSL(sslmode string) Connector {
	c.sslmode = sslmode
	return c
}

func (c *connector) createConnectionString() string {
	result := make([]string, 0)
	if len(c.host) > 0 {
		result = append(result, fmt.Sprintf("host=%s", c.host))
	}
	if c.port > 0 {
		result = append(result, fmt.Sprintf("port=%d", c.port))
	}
	if len(c.user) > 0 {
		result = append(result, fmt.Sprintf("user=%s", c.user))
	}
	if len(c.password) > 0 {
		result = append(result, fmt.Sprintf("password=%s", c.password))
	}
	if len(c.dbname) > 0 {
		result = append(result, fmt.Sprintf("dbname=%s", c.dbname))
	}
	if len(c.sslmode) > 0 {
		result = append(result, fmt.Sprintf("sslmode=%s", c.sslmode))
	}
	return strings.Join(result, " ")
}

func (c *connector) getPtr() *connector {
	return c
}
