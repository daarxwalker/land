package land

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).
		Update().
		SetValues(
			testModel{
				Name:     "Dominik",
				Lastname: "Linduska",
			},
		).
		SetVectors("Dominik", "Linduska")
	q.Return(Id, "name", "lastname")
	test.Equal(
		`UPDATE "tests" SET "name" = 'Dominik',"lastname" = 'Linduska',"active" = false,"vectors" =,"created_at" =,"updated_at" = CURRENT_TIMESTAMP RETURNING "id","name","lastname";`,
		q.GetSQL(),
	)
}
