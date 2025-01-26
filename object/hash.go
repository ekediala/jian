package object

import (
	"fmt"
	"strings"
)

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

type Hashable interface {
	HashKey() HashKey
}

func (h *Hash) Type() ObjectType {
	return HASH
}

func (h *Hash) Inspect() string {
	var s strings.Builder

	pairs := make([]string, 0, len(h.Pairs))

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	s.WriteByte('{')
	s.WriteString(strings.Join(pairs, ","))
	s.WriteByte('}')

	return s.String()
}
