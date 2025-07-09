package functions

import "github.com/simonwater/gopression/values"

type Callable interface {
	Arity() int
	Call(arguments []values.Value) (values.Value, error)
}

// 需要由具体函数实现
type CallableFunction interface {
	Callable
	GetName() string
	GetTitle() string
	GetGroup() string
}
