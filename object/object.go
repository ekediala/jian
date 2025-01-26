package object

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type ObjectType string

func (o ObjectType) String() string {
	return string(o)
}

const (
	INTEGER      ObjectType = "INTEGER"
	BOOLEAN      ObjectType = "BOOLEAN"
	NULL         ObjectType = "NULL"
	RETURN_VALUE ObjectType = "RETURN_VALUE"
	ERROR        ObjectType = "ERROR"
	FUNCTION     ObjectType = "FUNCTION"
	STRING       ObjectType = "STRING"
	BUILTIN      ObjectType = "BUILTIN"
	ARRAY        ObjectType = "ARRAY"
	HASH         ObjectType = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
