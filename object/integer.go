package object

import "fmt"

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: INTEGER, Value: uint64(i.Value)}
}
