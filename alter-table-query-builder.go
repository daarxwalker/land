package land

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type AlterTableQuery interface {
	AddColumn(name, dataType string, options ...ColOpts) AlterTableQuery
	RenameColumn(currentName, newName string) AlterTableQuery
	DropColumn(name string) AlterTableQuery
	GetSQL() string
	Exec()
	IfExists() AlterTableQuery
}

type alterTableQueryBuilder struct {
	*queryBuilder
	entity   *entity
	context  context.Context
	ifExists bool
	add      []*column
	rename   map[string]string
	drop     []string
}

func createAlterTableQuery(entity *entity) *alterTableQueryBuilder {
	return &alterTableQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(AlterTable),
		context:      context.Background(),
		entity:       entity,
		add:          make([]*column, 0),
		rename:       make(map[string]string),
		drop:         make([]string, 0),
	}
}

func (q *alterTableQueryBuilder) AddColumn(name, dataType string, options ...ColOpts) AlterTableQuery {
	opts := ColOpts{}
	if len(options) > 0 {
		opts = options[0]
	}
	c := createColumn(name, dataType, opts)
	q.add = append(q.add, c)
	return q
}

func (q *alterTableQueryBuilder) RenameColumn(currentName, newName string) AlterTableQuery {
	q.rename[currentName] = newName
	return q
}

func (q *alterTableQueryBuilder) DropColumn(name string) AlterTableQuery {
	q.drop = append(q.drop, name)
	return q
}

func (q *alterTableQueryBuilder) Exec() {
	createQueryManager(q.entity, q.context).setQuery(q.GetSQL()).setQueryType(AlterTable).exec()
}

func (q *alterTableQueryBuilder) GetSQL() string {
	return q.createQueryString()
}

func (q *alterTableQueryBuilder) IfExists() AlterTableQuery {
	q.ifExists = true
	return q
}

func (q *alterTableQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	result = append(result, "ALTER TABLE")
	if q.ifExists {
		result = append(result, "IF EXISTS")
	}
	result = append(result, q.escape(q.entity.name))
	alters := make([]string, 0)
	alters = append(alters, q.createAddParts()...)
	alters = append(alters, q.createRenameParts()...)
	alters = append(alters, q.createDropParts()...)
	result = append(result, strings.Join(alters, ","))
	return strings.Join(result, " ") + q.getQueryDivider()
}

func (q *alterTableQueryBuilder) createAddParts() []string {
	result := make([]string, 0)
	for _, c := range q.add {
		colSql := make([]string, 0)
		colSql = append(colSql, "ADD COLUMN", q.escape(c.name), q.createDataType(c))
		if c.options.NotNull {
			colSql = append(colSql, "NOT NULL")
		}
		if c.options.Unique {
			colSql = append(colSql, "UNIQUE")
		}
		if c.options.Default != nil {
			colSql = append(colSql, "DEFAULT", q.createValue(c, reflect.ValueOf(c.options.Default)))
		}
		if (c.options.Reference.Self && len(c.options.Reference.Column) > 0) || (c.options.Reference.Entity != nil && len(c.options.Reference.Column) > 0) {
			colSql = append(
				colSql, "REFERENCES", q.getReferenceName(c)+fmt.Sprintf("(%s)", q.escape(c.options.Reference.Column)),
			)
		}
		result = append(result, strings.Join(colSql, " "))
	}
	return result
}

func (q *alterTableQueryBuilder) createDropParts() []string {
	result := make([]string, 0)
	for _, c := range q.drop {
		result = append(result, fmt.Sprintf("DROP COLUMN %s", q.escape(c)))
	}
	return result
}

func (q *alterTableQueryBuilder) createRenameParts() []string {
	result := make([]string, 0)
	for currentName, newName := range q.rename {
		result = append(result, fmt.Sprintf("RENAME %s TO %s", q.escape(currentName), q.escape(newName)))
	}
	return result
}

func (q *alterTableQueryBuilder) getReferenceName(c *column) string {
	if c.options.Reference.Self && c.options.Reference.Entity == nil {
		return q.escape(q.entity.name)
	}
	return q.escape(c.options.Reference.Entity.getPtr().name)
}
