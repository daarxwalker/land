package land

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlterTable(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance().EntityManager()).
		AlterTable().IfExists().
		AddColumn("middle_name", Varchar, ColOpts{Limit: 255, NotNull: true, Unique: true}).
		RenameColumn("name", "custom_name").
		DropColumn("custom_name")
	test.Equal(`ALTER TABLE IF EXISTS "tests" ADD COLUMN "middle_name" VARCHAR(255) NOT NULL UNIQUE,RENAME "name" TO "custom_name",DROP COLUMN "custom_name";`, q.GetSQL())
}
