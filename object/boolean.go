package object

import "fmt"

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	var value uint64 = 0
	if b.Value {
		value = 1
	}

	return HashKey{Type: BOOLEAN, Value: value}
}
