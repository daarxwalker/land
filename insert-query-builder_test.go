package land

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).
		Insert().
		SetData(testModel{
			Name:     "Dominik",
			Lastname: "Linduska",
		}).
		SetVectors("Dominik", "Linduska")
	q.Return(Id, "name", "lastname")
	test.Equal(`INSERT INTO "tests" ("name","lastname","active","vectors","created_at","updated_at") VALUES ('Dominik','Linduska',false,to_tsvector('dominik linduska'),CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) RETURNING "id","name","lastname";`, q.GetSQL())
}
