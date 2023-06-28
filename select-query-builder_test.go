package land

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectBase(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).
		Select().
		All().
		GetSQL()
	test.Equal(`SELECT * FROM "tests" AS "t";`, q)
}

func TestSelectJoin(t *testing.T) {
	test := assert.New(t)
	ent2 := testSecondEntity(testCreatePostgresInstance())
	q := testEntity(testCreatePostgresInstance()).Select()
	q.Join().On(ent2, testLastname)
	q.All()
	test.Equal(`SELECT * FROM "tests" AS "t" LEFT JOIN tests AS t2 ON "t"."id" = "t2"."lastname";`, q.GetSQL())
}

func TestSelectColumns(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).Select()
	q.Columns("daar", "walker")
	q.Column("test").Alias("random")
	q.All()
	test.Equal(`SELECT "t"."daar","t"."walker","t"."test" AS "test" FROM "tests" AS "t";`, q.GetSQL())
}

func TestSelectWhere(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).Select()
	q.Column("name")
	q.Where().Column("name").Equal("daar")
	q.Where().Column("name").Equal("dominik").
		Or(q.Where().Column("name").Equal("linduska"))
	q.All()
	test.Equal(`SELECT "t"."name" FROM "tests" AS "t" WHERE "t"."name" = 'daar' AND ("t"."name" = 'dominik' OR "t"."name" = 'linduska');`, q.GetSQL())
}

func TestSelectGroup(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).Select()
	q.Group().Columns("test1", "test2")
	q.All()
	test.Equal(`SELECT * FROM "tests" AS "t" GROUP BY "t"."test1","t"."test2";`, q.GetSQL())
}

func TestSelectOrder(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).Select()
	q.Column("test")
	q.Order().Slice([]OrderParam{
		{Key: "test", Direction: "asc"},
	})
	q.Order().Asc("test_one").Desc("test_two")
	q.All()
	test.Equal(`SELECT "t"."test" FROM "tests" AS "t" ORDER BY "t"."test" ASC,"t"."test_one" ASC,"t"."test_two" DESC;`, q.GetSQL())
}

func TestSelectFulltext(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).Select()
	q.Column("test")
	q.Fulltext("test")
	q.All()
	test.Equal(`SELECT "t"."test" FROM "tests" AS "t" WHERE "t"."vectors" @@ to_tsquery('test:*');`, q.GetSQL())
}
