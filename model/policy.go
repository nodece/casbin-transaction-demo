package model

import (
	"strings"
)

type PolicyContainer struct {
	policy      [][]string
	policyIndex map[string]int
}

func NewPolicyContainer() *PolicyContainer {
	p := &PolicyContainer{
		policy:      make([][]string, 0),
		policyIndex: make(map[string]int),
	}
	return p
}

const DefaultSep = ","

func (p *PolicyContainer) Add(rules [][]string) {
	for _, rule := range rules {
		hashKey := strings.Join(rule, DefaultSep)
		_, ok := p.policyIndex[hashKey]
		if ok {
			continue
		}
		p.policy = append(p.policy, rule)
		p.policyIndex[hashKey] = len(p.policy) - 1
	}
}

func (p *PolicyContainer) Remove(rules [][]string) {
	for _, rule := range rules {
		hashKey := strings.Join(rule, DefaultSep)
		index, ok := p.policyIndex[hashKey]
		if !ok {
			continue
		}

		p.policy = append(p.policy[:index], p.policy[index+1:]...)
		delete(p.policyIndex, hashKey)
		for i := index; i < len(p.policy); i++ {
			p.policyIndex[strings.Join(p.policy[i], DefaultSep)] = i
		}
	}
}
