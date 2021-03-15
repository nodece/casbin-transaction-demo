package adapter

type Tx interface {
	Commit() error
	Rollback() error
	AddPolicy(sec string, ptype string, rule []string) error
	RemovePolicy(sec string, ptype string, rule []string) error
	CleanPolicy() error
	SavePolicy(policy map[string]map[string][]string) error
}

type Adapter interface {
	Update(fn func(Tx) error) error
	LoadPolicy(model *interface{}) error
}
