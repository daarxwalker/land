package land

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dd"
)

func TestDelete(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).
		Delete()
	q.Where().Column(Id).Equal(1)
	q.Return(Id)
	dd.Print(q.GetSQL())
	test.Equal(``, ``)
}
