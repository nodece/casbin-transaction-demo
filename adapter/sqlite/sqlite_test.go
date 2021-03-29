package sqlite

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDB(t *testing.T) {
	dsn := "file:policy.db?cache=shared&mode=memory"

	db, err := NewDB(dsn, "policy")
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestNewDBWithInvalidParams(t *testing.T) {
	db, err := NewDB("", "")
	assert.EqualError(t, err,"dsn is not provided")
	assert.Nil(t, db)

	db, err = NewDB("file:policy.db?cache=shared&mode=memory", "")
	assert.EqualError(t, err,"tableName is not provided")
	assert.Nil(t, db)
}
