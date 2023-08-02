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
		`UPDATE "tests" AS "t" SET "name" = 'Dominik',"lastname" = 'Linduska',"active" = false,"vectors" = to_tsvector('dominik linduska'),"updated_at" = CURRENT_TIMESTAMP RETURNING "id","name","lastname";`,
		q.GetSQL(),
	)
}

func TestUpdateMapValue(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).
		Update().
		SetValues(
			map[string]any{"name": "Dominik", "lastname": "Linduska"},
		)
	q.Return(Id, "name", "lastname")
	test.Equal(
		`UPDATE "tests" AS "t" SET "name" = 'Dominik',"lastname" = 'Linduska',"updated_at" = CURRENT_TIMESTAMP RETURNING "id","name","lastname";`,
		q.GetSQL(),
	)
}
