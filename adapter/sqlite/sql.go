package sqlite

import (
	"fmt"
)

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