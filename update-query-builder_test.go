package land

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dd"
)

func TestUpdate(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance().EntityManager()).
		Update().
		SetData(testModel{
			Name:     "Dominik",
			Lastname: "Linduska",
		}).
		SetTSVectors("Dominik", "Linduska")
	q.Return(Id, "name", "lastname")
	dd.Print(q.GetSQL())
	test.Equal(``, ``)
}
