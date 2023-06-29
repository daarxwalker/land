package land

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).
		Delete()
	q.Where().Column(Id).Equal(1)
	q.Return(Id)
	test.Equal(`DELETE FROM "tests" WHERE "t"."id" = 1 RETURNING "id";`, q.GetSQL())
}
