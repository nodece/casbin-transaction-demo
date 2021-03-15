package adapter

import (
	"casbin-transaction-demo/adapter/sqlite"
	"casbin-transaction-demo/model"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteAdapter struct {
	lock      *sync.RWMutex
	db        *sql.DB
	tableName string
}

var _ Adapter = &sqliteAdapter{}

// file:test.db?cache=shared&mode=memory
func NewSqliteAdapter(dsn string, tableName string) (Adapter, error) {
	s := &sqliteAdapter{
		lock:      &sync.RWMutex{},
		tableName: "policy",
	}

	if tableName != "" {
		s.tableName = tableName
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(sqlite.MigrateSchema(s.tableName))
	if err != nil {
		return nil, err
	}

	s.db = db

	return s, nil
}

func (s *sqliteAdapter) LoadPolicy(m *interface{}) error {
	root := (*m).(model.Model)

	s.lock.RLock()

	rows, err := s.db.Query(sqlite.QueryPolicySQL(s.tableName))
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			ptype string
			v0    string
			v1    string
			v2    string
			v3    string
			v4    string
			v5    string
		)

		err := rows.Scan(&ptype, &v1, &v2, &v3, &v4, &v5)
		if err != nil {
			return err
		}

		var rule []string
		if ptype != "" {
			rule = append(rule, ptype)
		}
		if v0 != "" {
			rule = append(rule, v0)
		}
		if v1 != "" {
			rule = append(rule, v1)
		}
		if v2 != "" {
			rule = append(rule, v2)
		}
		if v3 != "" {
			rule = append(rule, v3)
		}
		if v4 != "" {
			rule = append(rule, v4)
		}
		if v5 != "" {
			rule = append(rule, v5)
		}
		root.AddPolicy(rule)
	}

	err = rows.Err()
	return err
}

func (s *sqliteAdapter) begin() (Tx, error) {
	s.lock.Lock()
	t, err := NewSqliteAdapterTx(s.db, s.tableName)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *sqliteAdapter) Update(fn func(Tx) error) error {
	t, err := s.begin()
	defer s.lock.Unlock()
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
