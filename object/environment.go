package object

type Environment struct {
	store    map[string]Object
	outerEnv *Environment
}

func NewEnv() *Environment {
	return &Environment{store: make(map[string]Object), outerEnv: nil}
}

func NewWrappedEnv(outerEnv *Environment) *Environment {
	env := NewEnv()
	env.outerEnv = outerEnv
	return env
}

func (e *Environment) Get(key string) (Object, bool) {
	obj, ok := e.store[key]
	if !ok && e.outerEnv != nil {
		obj, ok = e.outerEnv.Get(key)
	}
	return obj, ok
}

func (e *Environment) Set(key string, obj Object) {
	e.store[key] = obj
}
