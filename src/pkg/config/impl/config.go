package impl

type Configer interface {
	MergeFromFile(f string) error
	Unmarshal(interface{}) error
}
