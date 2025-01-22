package object

import "fmt"

type Error struct {
	Message string
}

func (e *Error) Inspect() string {
	return e.Message
}

func (e *Error) Type() ObjectType {
	return ERROR
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(format string, args ...interface{}) *Error {
	error := Error{
		Message: fmt.Sprintf(format, args...),
	}
	return &error
}
