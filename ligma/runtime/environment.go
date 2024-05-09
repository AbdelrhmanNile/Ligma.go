package runtime

type Environment struct {
	store map[string]LigmaObject
	parent *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]LigmaObject)

	return &Environment{store: s, parent: nil}
}

func (e *Environment) Get(name string) (LigmaObject, bool) {
	obj, ok := e.store[name]
	if !ok && e.parent != nil {
		obj, ok = e.parent.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val LigmaObject) LigmaObject {
	e.store[name] = val
	return val
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.parent
	}
	return env
}

func (e *Environment) GetAt(distance int, name string) (LigmaObject, bool) {
	env := e.ancestor(distance)
	obj, ok := env.store[name]
	return obj, ok
}

func (e *Environment) SetAt(distance int, name string, val LigmaObject) LigmaObject {
	env := e.ancestor(distance)
	env.store[name] = val
	return val
}


func NewEnclosedEnvironment(parent *Environment) *Environment {
	env := NewEnvironment()
	env.parent = parent
	return env
}