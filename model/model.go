package model

import "sync"

type Model struct {
	policy *PolicyContainer
	rwLock *sync.RWMutex
}

func NewModel() *Model {
	return &Model{
		policy: NewPolicyContainer(),
		rwLock: &sync.RWMutex{},
	}
}

func (m *Model) AddPolicy(rule []string) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	m.policy.Add([][]string{rule})
}

func (m *Model) AddPolicies(rule [][]string) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	m.policy.Add(rule)
}

func (m *Model) RemovePolicy(rule []string) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	m.policy.Remove([][]string{rule})
}

func (m *Model) RemovePolicies(rule [][]string) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	m.policy.Remove(rule)
}

func (m *Model) CopyPolicy(dst *PolicyContainer) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	*dst = *m.policy
}

func (m *Model) SetPolicy(policy *PolicyContainer) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	*m.policy = *policy
}

func (m *Model) GetPolicyWithRLock() [][]string {
	m.rwLock.RLock()
	return m.policy.policy
}

func (m *Model) RUnlockPolicy() {
	m.rwLock.RUnlock()
}
