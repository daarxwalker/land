package land

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumn(t *testing.T) {
	test := assert.New(t)
	e := testEntity(testCreatePostgresInstance())
	c := createColumnQuery(e.getPtr(), "name")
	test.Equal(`"t"."name"`, c.getQueryString())
}

func TestColumnCount(t *testing.T) {
	test := assert.New(t)
	e := testEntity(testCreatePostgresInstance())
	c := createColumnQuery(e.getPtr(), "name")
	c.Count().Alias("count_name")
	test.Equal(`COUNT("t"."name") AS "count_name"`, c.getQueryString())
}

func TestColumnSum(t *testing.T) {
	test := assert.New(t)
	e := testEntity(testCreatePostgresInstance())
	c := createColumnQuery(e.getPtr(), "name")
	c.Sum().Alias("sum_name")
	test.Equal(`SUM("t"."name") AS "sum_name"`, c.getQueryString())
}

func TestColumnMin(t *testing.T) {
	test := assert.New(t)
	e := testEntity(testCreatePostgresInstance())
	c := createColumnQuery(e.getPtr(), "posts")
	c.Min().Alias("min_posts")
	test.Equal(`MIN("t"."posts") AS "min_posts"`, c.getQueryString())
}

func TestColumnMax(t *testing.T) {
	test := assert.New(t)
	e := testEntity(testCreatePostgresInstance())
	c := createColumnQuery(e.getPtr(), "posts")
	c.Max().Alias("max_posts")
	test.Equal(`MAX("t"."posts") AS "max_posts"`, c.getQueryString())
}

func TestColumnArrayAgg(t *testing.T) {
	test := assert.New(t)
	e := testEntity(testCreatePostgresInstance())
	c := createColumnQuery(e.getPtr(), "posts")
	c.ArrayAgg().Alias("posts_array")
	test.Equal(`ARRAY_AGG("t"."posts") AS "posts_array"`, c.getQueryString())
}

func TestColumnStringAgg(t *testing.T) {
	test := assert.New(t)
	e := testEntity(testCreatePostgresInstance())
	c1 := createColumnQuery(e.getPtr(), "name")
	empty := createColumnQuery(nil, "").Empty().Separator(" ")
	c2 := createColumnQuery(e.getPtr(), "lastname")
	c1.StringAgg(empty, c2).Separator(",").Alias("value")
	test.Equal(`STRING_AGG("t"."name" || ' ' || "t"."lastname", ',') AS "value"`, c1.getQueryString())
}
