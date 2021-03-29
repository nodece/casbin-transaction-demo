package adapter

import (
	"casbin-transaction-demo/adapter/sqlite"
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

type SQLiteAdapterTx struct {
	tableName string

	tx         *sql.Tx
	db         *sql.DB
	rootLock   *sync.RWMutex
	addStmt    *sql.Stmt
	removeStmt *sql.Stmt
	cleanStmt  *sql.Stmt
}

func NewSQLiteAdapterTx(rootLock *sync.RWMutex, db *sql.DB, tableName string) (*SQLiteAdapterTx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	s := &SQLiteAdapterTx{rootLock: rootLock, tx: tx, db: db, tableName: tableName}

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

func (s *SQLiteAdapterTx) Commit() error {
	s.rootLock.Unlock()
	return s.tx.Commit()
}

func (s *SQLiteAdapterTx) Rollback() error {
	s.rootLock.Unlock()
	return s.tx.Rollback()
}

func (s *SQLiteAdapterTx) AddPolicy(sec string, ptype string, rule []string) error {
	return s.AddPolicyContext(context.Background(), sec, ptype, rule)
}

func (s *SQLiteAdapterTx) AddPolicyContext(ctx context.Context, _ string, ptype string, rule []string) error {
	params := make([]interface{}, 7)
	params[0] = ptype
	for i, v := range rule {
		params[i+1] = v
	}

	_, err := s.addStmt.ExecContext(ctx, params...)
	return err
}

func (s *SQLiteAdapterTx) RemovePolicy(sec string, ptype string, rule []string) error {
	return s.RemovePolicyContext(context.Background(), sec, ptype, rule)
}

func (s *SQLiteAdapterTx) RemovePolicyContext(ctx context.Context, _ string, ptype string, rule []string) error {
	params := make([]interface{}, 7)
	params[0] = ptype
	for i, v := range rule {
		params[i+1] = v
	}

	_, err := s.removeStmt.ExecContext(ctx, params...)
	return err
}

func (s *SQLiteAdapterTx) CleanPolicy() error {
	return s.CleanPolicyContext(context.Background())
}

func (s *SQLiteAdapterTx) CleanPolicyContext(ctx context.Context) error {
	_, err := s.cleanStmt.ExecContext(ctx)
	return err
}
