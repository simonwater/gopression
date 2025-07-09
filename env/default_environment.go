package env

import (
	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/values"
)

type DefaultEnvironment struct {
	*Environment
	data map[string]values.Value
}

func NewDefaultEnvironment() *DefaultEnvironment {
	env := &DefaultEnvironment{
		data: make(map[string]values.Value),
	}
	env.Environment = NewEnvironment(env)
	return env
}
func (de *DefaultEnvironment) BeforeExecute(vars []*util.Field) bool {
	return true
}
func (de *DefaultEnvironment) Get(id string) values.Value {
	if v, ok := de.data[id]; ok {
		return v
	}
	return values.NewNullValue()
}
func (de *DefaultEnvironment) GetOrDefault(id string, defValue values.Value) values.Value {
	if v, ok := de.data[id]; ok {
		return v
	}
	return defValue
}
func (de *DefaultEnvironment) Put(id string, value values.Value) {
	de.data[id] = value
}
func (de *DefaultEnvironment) Size() int {
	return len(de.data)
}
