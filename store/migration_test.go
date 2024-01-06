package store

import (
	"testing"
)

func TestMigration(t *testing.T) {
	_, teardown := SetupTest(t)
	defer teardown(t)
}
