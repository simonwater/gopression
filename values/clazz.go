package values

type Clazz struct {
	Name       string
	SuperClass *Clazz
	Methods    map[string]*ClazzMethod
}

func NewClazz(name string, superClass *Clazz, methods map[string]*ClazzMethod) *Clazz {
	return &Clazz{
		Name:       name,
		SuperClass: superClass,
		Methods:    methods,
	}
}
