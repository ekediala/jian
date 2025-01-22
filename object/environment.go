package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: map[string]Object{},
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{
		store: map[string]Object{},
		outer: outer,
	}
}

func (e *Environment) Get(key string) (Object, bool) {
	obj, ok := e.store[key]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(key)
	}

	return obj, ok
}

func (e *Environment) Set(key string, value Object) Object {
	e.store[key] = value
	return value
}
