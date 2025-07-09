package values

type ClazzMethod struct{}

func (m *ClazzMethod) Arity() int {
	// TODO: 实现参数个数逻辑
	return 0
}

func (m *ClazzMethod) Call(arguments []Value) (Value, error) {
	// TODO: 实现方法调用逻辑
	return Value{}, nil
}
