package adapter

import (
	"casbin-transaction-demo/adapter/sqlite"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteAdapterTx struct {
	tableName string

	tx *sql.Tx

	addStmt    *sql.Stmt
	removeStmt *sql.Stmt
	cleanStmt  *sql.Stmt
}

func NewSqliteAdapterTx(db *sql.DB, tableName string) (*sqliteAdapterTx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	s := &sqliteAdapterTx{tx: tx, tableName: tableName}

	addStmt, err := tx.Prepare(sqlite.AddPolicySQL(tableName))
	if err != nil {
		return nil, err
	}
	s.addStmt = addStmt

	removeStmt, err := tx.Prepare(sqlite.RemovePolicySQL(tableName))
	if err != nil {
		return nil, err
	}
	s.removeStmt = removeStmt

	cleanStmt, err := tx.Prepare(sqlite.CleanPolicySQL(tableName))
	if err != nil {
		return nil, err
	}
	s.cleanStmt = cleanStmt

	return s, nil
}

func (s *sqliteAdapterTx) Commit() error {
	return s.tx.Commit()
}

func (s *sqliteAdapterTx) Rollback() error {
	return s.tx.Rollback()
}

func (s *sqliteAdapterTx) AddPolicy(_ string, ptype string, rule []string) error {
	params := make([]interface{}, 7)
	params[0] = ptype
	for i, v := range rule {
		params[i+1] = v
	}

	_, err := s.addStmt.Exec(params...)
	return err
}

func (s *sqliteAdapterTx) RemovePolicy(_ string, ptype string, rule []string) error {
	params := make([]interface{}, 7)
	params[0] = ptype
	for i, v := range rule {
		params[i+1] = v
	}

	_, err := s.removeStmt.Exec(params...)
	return err
}

func (s *sqliteAdapterTx) CleanPolicy() error {
	_, err := s.cleanStmt.Exec()
	return err
}

func (s *sqliteAdapterTx) SavePolicy(policy map[string]map[string][]string) error {
	err := s.CleanPolicy()
	if err != nil {
		return err
	}

	for _, value := range policy {
		for ptype, rule := range value {
			err = s.AddPolicy("", ptype, rule)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
