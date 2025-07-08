package gopression

import (
	"github.com/simonwater/gopression/values"
)

type Instance struct {
	Clazz  *Clazz
	Fields map[string]values.Value
}

func NewInstance() *Instance {
	return &Instance{
		Fields: make(map[string]values.Value),
	}
}

func NewInstanceWithClazz(clazz *Clazz) *Instance {
	return &Instance{
		Clazz:  clazz,
		Fields: make(map[string]values.Value),
	}
}

func (inst *Instance) Get(name string) (values.Value, bool) {
	if v, ok := inst.Fields[name]; ok {
		return v, true
	}
	return values.Value{}, false
}

func (inst *Instance) Set(name string, value values.Value) {
	inst.Fields[name] = value
}

func (inst *Instance) String() string {
	if inst.Clazz != nil {
		return inst.Clazz.Name + " instance"
	}
	return "instance"
}
