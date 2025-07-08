package util

import (
	"fmt"
)

type NodeSet[T any] struct {
	nodesMap map[string]*Node[T]
	nodes    []*Node[T]
	cnt      int
}

func NewNodeSet[T any]() *NodeSet[T] {
	return &NodeSet[T]{
		nodesMap: make(map[string]*Node[T]),
		nodes:    make([]*Node[T], 0),
	}
}

func (ns *NodeSet[T]) AddNode(name string) *Node[T] {
	if node, ok := ns.nodesMap[name]; ok {
		return node
	}
	node := NewNode[T](name, ns.cnt)
	ns.cnt++
	ns.nodesMap[name] = node
	ns.nodes = append(ns.nodes, node)
	return node
}

func (ns *NodeSet[T]) GetNodeByName(name string) *Node[T] {
	return ns.nodesMap[name]
}

func (ns *NodeSet[T]) GetNode(index int) *Node[T] {
	ns.validateIndex(index)
	return ns.nodes[index]
}

func (ns *NodeSet[T]) validateIndex(i int) {
	if i < 0 || i >= ns.cnt {
		panic(fmt.Sprintf("index %d is not between 0 and %d", i, ns.cnt-1))
	}
}

func (ns *NodeSet[T]) Size() int {
	return ns.cnt
}

func (ns *NodeSet[T]) String() string {
	result := ""
	for i := 0; i < ns.cnt; i++ {
		node := ns.GetNode(i)
		result += fmt.Sprintf("%d: %s(%d)\n", i, node.Name, node.Index)
	}
	return result
}
