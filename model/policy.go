package model

type Policy [][]string

func (p *Policy) Add(rule []string) {
	*p = append(*p, rule)
}

func (p *Policy) Remove(rule []string) {
	for i := range *p {
		if equal((*p)[i], rule) {
			*p = append((*p)[:i], (*p)[i+1:]...)
			break
		}
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

