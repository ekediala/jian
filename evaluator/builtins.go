package evaluator

import (
	"fmt"

	"github.com/ekediala/jian/object"
)

var builtins = map[string]*object.Builtin{
	"len":   {Fn: length},
	"first": {Fn: first},
	"last":  {Fn: last},
	"rest":  {Fn: rest},
	"push":  {Fn: push},
	"puts":  {Fn: puts},
}

func puts(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NULL
}

func push(args ...object.Object) object.Object {
	if got := len(args); got != 2 {
		return object.NewError("wrong number of arguments. got=%d, want=2", got)
	}

	if arr, ok := args[0].(*object.Array); ok {
		elements := make([]object.Object, len(arr.Elements)+1)
		copy(elements, arr.Elements)
		elements[len(arr.Elements)] = args[1]
		return &object.Array{Elements: elements}
	}

	return object.NewError("argument to `push` must be ARRAY, got %s", args[0].Type())
}

func rest(args ...object.Object) object.Object {
	if got := len(args); got != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", got)
	}

	if arg, ok := args[0].(*object.Array); ok {
		if len(arg.Elements) > 0 {
			return &object.Array{Elements: arg.Elements[1:]}
		}

		return NULL
	}

	return object.NewError("argument to `rest` must be ARRAY, got %s", args[0].Type())
}

func last(args ...object.Object) object.Object {
	if got := len(args); got != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", got)
	}

	if arg, ok := args[0].(*object.Array); ok {
		if len(arg.Elements) > 0 {
			return arg.Elements[len(arg.Elements)-1]
		}

		return NULL
	}

	return object.NewError("argument to `last` must be ARRAY, got %s", args[0].Type())
}

func first(args ...object.Object) object.Object {
	if got := len(args); got != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", got)
	}

	if arg, ok := args[0].(*object.Array); ok {
		if len(arg.Elements) > 0 {
			return arg.Elements[0]
		}

		return NULL
	}

	return object.NewError("argument to `first` must be ARRAY, got %s", args[0].Type())
}

func length(args ...object.Object) object.Object {
	if got := len(args); got != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", got)
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return object.NewError("argument to `len` not supported, got %v", arg.Type())
	}
}
