package land

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).
		Update().
		SetData(testModel{
			Name:     "Dominik",
			Lastname: "Linduska",
		}).
		SetTSVectors("Dominik", "Linduska")
	q.Return(Id, "name", "lastname")
	test.Equal(`UPDATE "tests" SET "name" = 'Dominik',"lastname" = 'Linduska',"active" = false RETURNING "id","name","lastname";`, q.GetSQL())
}
