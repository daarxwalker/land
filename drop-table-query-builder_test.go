package land

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDropTable(t *testing.T) {
	test := assert.New(t)
	q := testEntity(testCreatePostgresInstance()).
		DropTable().IfExists()
	test.Equal(`DROP TABLE IF EXISTS "tests";`, q.GetSQL())
}
