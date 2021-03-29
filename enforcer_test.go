package main

import (
	"casbin-transaction-demo/adapter"
	"casbin-transaction-demo/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestEnforcer(t *testing.T) {
	a, err := adapter.NewSqliteAdapter("file:policy.db?cache=shared&mode=memory", "")
	assert.NoError(t, err)

	e := Enforcer{
		txLock:  &sync.RWMutex{},
		adapter: a,
		model:   model.NewModel(),
	}
	err = e.LoadPolicy()
	assert.NoError(t, err)

	assert.Equal(t, [][]string{}, e.GetPolicy())

	err = e.Update(func(tx *Tx) error {
		err := tx.Add("p", "p", []string{"root", "/", "GET"})
		if err != nil {
			return err
		}
		err = tx.Add("p", "p", []string{"root", "/", "POST"})
		if err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err)

	assert.Equal(t, [][]string{
		{"root", "/", "GET"},
		{"root", "/", "POST"},
	}, e.GetPolicy())

	throwErr := errors.New("throw an exception to test the rollback")
	err = e.Update(func(tx *Tx) error {
		err := tx.Add("p", "p", []string{"root", "/", "DELETE"})
		if err != nil {
			return err
		}
		err = tx.Add("p", "p", []string{"root", "/", "PUT"})
		if err != nil {
			return err
		}
		return throwErr
	})
	assert.EqualError(t, err, throwErr.Error())
	assert.Equal(t, [][]string{
		{"root", "/", "GET"},
		{"root", "/", "POST"},
	}, e.GetPolicy())
}
