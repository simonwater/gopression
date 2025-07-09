package values

type Instance struct {
	Clazz  *Clazz
	Fields map[string]Value
}

func NewInstance() *Instance {
	return &Instance{
		Fields: make(map[string]Value),
	}
}

func NewInstanceWithClazz(clazz *Clazz) *Instance {
	return &Instance{
		Clazz:  clazz,
		Fields: make(map[string]Value),
	}
}

func (inst *Instance) Get(name string) (Value, bool) {
	if v, ok := inst.Fields[name]; ok {
		return v, true
	}
	return Value{}, false
}

func (inst *Instance) Set(name string, value Value) {
	inst.Fields[name] = value
}

func (inst *Instance) String() string {
	if inst.Clazz != nil {
		return inst.Clazz.Name + " instance"
	}
	return "instance"
}
