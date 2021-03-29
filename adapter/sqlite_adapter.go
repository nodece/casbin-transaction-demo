package adapter

import (
	"casbin-transaction-demo/adapter/sqlite"
	"casbin-transaction-demo/model"
	"context"
	"database/sql"
	"fmt"
	"sync"
)

type SQLiteAdapter struct {
	lock      *sync.RWMutex
	db        *sql.DB
	tableName string
}

var _ Adapter = &SQLiteAdapter{}

// file:test.db?cache=shared&mode=memory
func NewSqliteAdapter(dsn string, tableName string) (*SQLiteAdapter, error) {
	s := &SQLiteAdapter{
		lock:      &sync.RWMutex{},
		tableName: "policy",
	}

	if tableName != "" {
		s.tableName = tableName
	}

	db, err := sqlite.NewDB(dsn, s.tableName)
	if err != nil {
		return nil, err
	}
	s.db = db

	return s, nil
}

func (s *SQLiteAdapter) LoadPolicy(_ context.Context, m *interface{}) error {
	root := (*m).(*model.Model)

	s.lock.RLock()
	defer s.lock.RUnlock()

	rows, err := s.db.Query(sqlite.QueryPolicySQL(s.tableName))
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			ptype sql.NullString
			v0    sql.NullString
			v1    sql.NullString
			v2    sql.NullString
			v3    sql.NullString
			v4    sql.NullString
			v5    sql.NullString
		)

		err := rows.Scan(&ptype, &v0, &v1, &v2, &v3, &v4, &v5)
		if err != nil {
			return err
		}

		var rule []string
		if ptype.Valid {
			rule = append(rule, ptype.String)
		}
		if v0.Valid {
			rule = append(rule, v0.String)
		}
		if v1.Valid {
			rule = append(rule, v1.String)
		}
		if v2.Valid {
			rule = append(rule, v2.String)
		}
		if v3.Valid {
			rule = append(rule, v3.String)
		}
		if v4.Valid {
			rule = append(rule, v4.String)
		}
		if v5.Valid {
			rule = append(rule, v5.String)
		}
		root.AddPolicy(rule)
	}

	err = rows.Err()
	return err
}

func (s *SQLiteAdapter) begin() (Tx, error) {
	s.lock.Lock()
	t, err := NewSQLiteAdapterTx(s.lock, s.db, s.tableName)
	if err != nil {
		s.lock.Unlock()
		return nil, err
	}
	return t, nil
}

func (s *SQLiteAdapter) Update(fn func(Tx) error) error {
	t, err := s.begin()
	if err != nil {
		return err
	}

	err = fn(t)
	if err != nil {
		rollbackErr := t.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("%v; %v", err, rollbackErr)
		}
		return err
	}

	return t.Commit()
}

func (s *SQLiteAdapter) Begin() (Tx, error) {
	return s.begin()
}

func (s *SQLiteAdapter) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}
