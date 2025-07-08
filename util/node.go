package util

type Node[T any] struct {
	Name  string
	Info  T
	Index int
}

func NewNode[T any](name string, index int) *Node[T] {
	return &Node[T]{
		Name:  name,
		Index: index,
	}
}
