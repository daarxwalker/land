package land

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTable(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).
		CreateTable().IfNotExists()
	test.Equal(`CREATE TABLE IF NOT EXISTS "tests" ("name" VARCHAR(255) NOT NULL,"lastname" VARCHAR(255) NOT NULL,"active" BOOLEAN NOT NULL DEFAULT false,"vectors" TSVECTOR NOT NULL DEFAULT to_tsquery(''),"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);`, q.GetSQL())
}
