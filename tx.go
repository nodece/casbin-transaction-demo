package main

import (
	"casbin-transaction-demo/adapter"
	"casbin-transaction-demo/model"
)

type Tx struct {
	policy       *model.Model
	root         model.PolicyContainer
	adapterTx    adapter.Tx
}

func (t *Tx) Rollback() error {
	err := t.adapterTx.Rollback()
	if err != nil {
		return err
	}

	t.policy = nil
	return nil
}

func (t *Tx) Commit() error {
	if t.policy == nil {
		return nil
	}
	err := t.adapterTx.Commit()
	if err != nil {
		return err
	}

	t.policy.SetPolicy(&t.root)
	return nil
}

func (t *Tx) Add(sec string, ptype string, rule []string) error {
	err := t.adapterTx.RemovePolicy(sec, ptype, rule)
	if err != nil {
		return err
	}

	t.root.Add([][]string{rule})
	return nil
}

func (t *Tx) Remove(sec string, ptype string, rule []string) error {
	err := t.adapterTx.RemovePolicy(sec, ptype, rule)
	if err != nil {
		return err
	}

	t.root.Remove([][]string{rule})
	return nil
}
