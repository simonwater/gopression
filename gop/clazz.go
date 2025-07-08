package gopression

import "github.com/simonwater/gopression/functions"

type Clazz struct {
	Name       string
	SuperClass *Clazz
	Methods    map[string]*functions.ClazzMethod
}

func NewClazz(name string, superClass *Clazz, methods map[string]*functions.ClazzMethod) *Clazz {
	return &Clazz{
		Name:       name,
		SuperClass: superClass,
		Methods:    methods,
	}
}
