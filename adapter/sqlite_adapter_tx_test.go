package adapter

import (
	"casbin-transaction-demo/adapter/sqlite"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustSQLiteAdapterTx() *SQLiteAdapterTx {
	dsn := "file:policy.db?cache=shared&mode=memory"
	tableName := "policy"

	db, err := sqlite.NewDB(dsn, tableName)
	if err != nil {
		panic(err)
	}

	tx, err := NewSQLiteAdapterTx(db, tableName)
	if err != nil {
		panic(err)
	}

	return tx
}

func TestNewSqliteAdapterTx(t *testing.T) {
	dsn := "file:policy.db?cache=shared&mode=memory"
	tableName := "policy"

	db, err := sqlite.NewDB(dsn, tableName)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	tx, err := NewSQLiteAdapterTx(db, tableName)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
}

func TestSqliteAdapterTx(t *testing.T) {
	tx := mustSQLiteAdapterTx()

	err := tx.AddPolicy("p", "p", []string{"root", "/res1", "*"})
	assert.NoError(t, err)

	err = tx.AddPolicy("p", "p", []string{"root", "/res2", "*"})
	assert.NoError(t, err)

	err = tx.RemovePolicy("p", "p", []string{"root", "/res1", "*"})
	assert.NoError(t, err)

	err = tx.CleanPolicy()
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)
}
