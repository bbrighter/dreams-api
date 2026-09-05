package migrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMigration(t *testing.T) {
	logger := zap.NewExample()
	db := NewDatabase(":memory:", logger)

	assert.NotPanics(t, func() {
		Migration(db, logger)
	})
}
