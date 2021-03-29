package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// NewDB returns a DB that will automatically migrate the schema.
func NewDB(dsn string, tableName string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("dsn is not provided")
	}
	if tableName == "" {
		return nil, errors.New("tableName is not provided")
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(MigrateSchema(tableName))
	if err != nil {
		return nil, err
	}

	return db, err
}

func MigrateSchema(tableName string) string {
	return fmt.Sprintf(`
CREATE table IF NOT EXISTS %s
(
    ptype TEXT,
    v0    TEXT,
    v1    TEXT,
    v2    TEXT,
    v3    TEXT,
    v4    TEXT,
    v5    TEXT,
    UNIQUE (ptype, v0, v1, v2, v3, v4, v5)
)
`, tableName)
}

func QueryPolicySQL(tableName string) string {
	return fmt.Sprintf(`SELECT ptype, v0, v1, v2, v3, v4, v5 from %s`, tableName)
}

func AddPolicySQL(tableName string) string {
	return fmt.Sprintf("insert into %s values (?,?,?,?,?,?,?)", tableName)
}

func RemovePolicySQL(tableName string) string {
	return fmt.Sprintf("delete from %s where ptype=? and v0=? and v1=? and v2=? and v3=? and v4=? and v5=?", tableName)
}

func CleanPolicySQL(tableName string) string {
	return fmt.Sprintf("delete from %s", tableName)
}
