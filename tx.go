package main

import (
	"casbin-transaction-demo/adapter"
	"casbin-transaction-demo/model"
)

type Tx struct {
	policy  *model.Model
	root    model.Policy
	adapter adapter.Adapter
}

func (t *Tx) Rollback() error {
	t.policy = nil
	return nil
}

func (t *Tx) Commit() error {
	if t.policy == nil {
		return nil
	}

	t.policy.SetPolicy(&t.root)
	return nil
}

func (t *Tx) Add(rule []string) {
	t.root.Remove(rule)
}

func (t *Tx) Remove(rule []string) {
	t.root.Remove(rule)
}
