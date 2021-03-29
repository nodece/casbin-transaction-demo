package adapter

import (
	"bytes"
	"casbin-transaction-demo/model"
	"context"
	"encoding/csv"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"testing"
)

func equalPolicy(t *testing.T, a *SQLiteAdapter, expected [][]string) {
	var m interface{}
	originalM := model.NewModel()
	m = originalM
	err := a.LoadPolicy(context.Background(), &m)
	assert.NoError(t, err)

	assert.Equal(t, originalM.GetPolicyWithRLock(), expected)
	defer originalM.RUnlockPolicy()
}

func TestNewSqliteAdapter(t *testing.T) {
	a, err := NewSqliteAdapter("file:policy.db?cache=shared&mode=memory", "")
	assert.NoError(t, err)
	assert.NotNil(t, a)

	err = a.Close()
	assert.NoError(t, err)
}

func TestSqliteAdapter_LoadPolicy(t *testing.T) {
	a, err := NewSqliteAdapter(path.Join("testdata", "policy.sqlite"), "policy")
	assert.NoError(t, err)
	defer a.Close()

	var m interface{}
	originalM := model.NewModel()
	m = originalM
	err = a.LoadPolicy(context.Background(), &m)
	assert.NoError(t, err)

	b, err := ioutil.ReadFile(path.Join("testdata", "main_policy.csv"))
	assert.NoError(t, err)

	reader := csv.NewReader(bytes.NewBuffer(b))
	reader.Comma = ','
	reader.Comment = '#'
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	for rootIndex, record := range records {
		for i, s := range record {
			if s == "" {
				records[rootIndex] = record[:i]
				break
			}
		}
	}
	assert.Equal(t, originalM.GetPolicyWithRLock(), records)
	defer originalM.RUnlockPolicy()
}

func TestSqliteAdapter_Close(t *testing.T) {
	a, err := NewSqliteAdapter("file:policy.db?cache=shared&mode=memory", "")
	assert.NoError(t, err)
	err = a.Close()
	assert.NoError(t, err)
}

func TestNewSqliteAdapter_Update(t *testing.T) {
	a, err := NewSqliteAdapter("file:policy.db?cache=shared&mode=memory", "")
	assert.NoError(t, err)
	defer a.Close()

	err = a.Update(func(tx Tx) error {
		err := tx.AddPolicy("p", "p", []string{"root", "/res1", "GET"})
		if err != nil {
			return err
		}
		err = tx.AddPolicy("p", "p", []string{"root", "/res2", "GET"})
		if err != nil {
			return err
		}
		err = tx.RemovePolicy("p", "p", []string{"root", "/res3", "GET"})
		if err != nil {
			return err
		}

		return nil
	})
	assert.NoError(t, err)
	equalPolicy(t, a, [][]string{
		{"p", "root", "/res1", "GET"},
		{"p", "root", "/res2", "GET"},
	})

	throwErr := errors.New("throw an exception to test the rollback")
	err = a.Update(func(tx Tx) error {
		err := tx.CleanPolicy()
		if err != nil {
			return err
		}
		err = tx.AddPolicy("p", "p", []string{"root", "/res1", "GET"})
		if err != nil {
			return err
		}
		err = tx.AddPolicy("p", "p", []string{"root", "/res2", "GET"})
		if err != nil {
			return err
		}
		return throwErr
	})
	assert.EqualError(t, err, throwErr.Error())
	equalPolicy(t, a, [][]string{
		{"p", "root", "/res1", "GET"},
		{"p", "root", "/res2", "GET"},
	})
}
