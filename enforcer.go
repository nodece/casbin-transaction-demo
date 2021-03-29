package main

import (
	"casbin-transaction-demo/adapter"
	"casbin-transaction-demo/model"
	"context"
	"sync"
)

type Enforcer struct {
	txLock  *sync.RWMutex
	adapter adapter.Adapter
	model   *model.Model
}

func (e *Enforcer) LoadPolicy() error {
	m := model.NewModel()
	e.model = m
	var p interface{}
	p = m
	return e.adapter.LoadPolicy(context.Background(), &p)
}

func (e *Enforcer) GetPolicy() [][]string {
	data := e.model.GetPolicyWithRLock()
	defer e.model.RUnlockPolicy()
	return data
}

func (e *Enforcer) begin() (*Tx, error) {
	e.txLock.Lock()
	adapterTx, err := e.adapter.Begin()
	if err != nil {
		return nil, err
	}

	tx := &Tx{}
	tx.adapterTx = adapterTx
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
