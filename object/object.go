package object

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
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
