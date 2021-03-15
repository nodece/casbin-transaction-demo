package main

import (
	"casbin-transaction-demo/model"
	"sync"
)

type Enforcer struct {
	txLock *sync.RWMutex

	model *model.Model
}

func (e *Enforcer) begin() (*Tx, error) {
	e.txLock.Lock()

	tx := &Tx{}
	e.model.CopyPolicy(&tx.root)

	tx.policy = e.model

	return tx, nil
}

func (e *Enforcer) Update(fn func(*Tx) error) error {
	tx, err := e.begin()
	if err != nil {
		return err
	}
	defer e.txLock.Unlock()

	err = fn(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func main() {
}
