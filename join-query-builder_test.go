package land

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeftJoin(t *testing.T) {
	test := assert.New(t)
	e1 := testEntity(testCreatePostgresInstance().EntityManager())
	e2 := testSecondEntity(testCreatePostgresInstance().EntityManager())
	join := createJoinQuery(e1.getPtr())
	join.On(e2)
	test.Equal(`LEFT JOIN tests AS t2 ON "t"."id" = "t2"."id"`, strings.Join(join.createQueryString(), " "))
}
