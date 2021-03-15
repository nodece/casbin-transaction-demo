package model

import "sync"

type Model struct {
	policy *Policy
	rwLock *sync.RWMutex
}

func (m *Model) AddPolicy (rule []string) {
	m.policy.Add(rule)
}

func (m *Model) RemovePolicy (rule []string) {
	m.policy.Remove(rule)
}

func (m *Model) CopyPolicy(dst *Policy) {
	*dst = *m.policy
}

func (m *Model) SetPolicy(policy *Policy) {
	*m.policy = *policy
}

func (m *Model) RLockPolicy() *Policy {
	m.rwLock.RLock()
	return m.policy
}

func (m *Model) RUnLockPolicy() {
	m.rwLock.RUnlock()
}
