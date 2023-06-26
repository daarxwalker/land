package land

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance().EntityManager()).
		Insert().
		SetData(testModel{
			Name:     "Dominik",
			Lastname: "Linduska",
		}).
		SetTSVectors("Dominik", "Linduska")
	q.Return(Id, "name", "lastname")
	test.Equal(``, ``)
}
