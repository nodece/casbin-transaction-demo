package adapter

import (
	"context"
	"sync"
)

type Tx interface {
	Commit() error
	Rollback() error

	AddPolicy(sec string, ptype string, rule []string) error
	AddPolicyContext(ctx context.Context, sec string, ptype string, rule []string) error
	RemovePolicy(sec string, ptype string, rule []string) error
	RemovePolicyContext(ctx context.Context, sec string, ptype string, rule []string) error
	CleanPolicy() error
	CleanPolicyContext(ctx context.Context) error
}

type Adapter interface {
	Update(fn func(Tx) error) error
	LoadPolicy(ctx context.Context, model *interface{}) error
	Begin() (Tx,error)
}

func x() {
	l := sync.RWMutex{}
	l.Unlock()
	l.Lock()
	l.RLock()
	l.RUnlock()
}
