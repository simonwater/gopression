package util

import "fmt"

type TopologicalSort struct {
	g        Digraph
	indegree []int
	order    []int
}

func NewTopologicalSort(g Digraph) *TopologicalSort {
	return &TopologicalSort{
		g:        g,
		indegree: make([]int, g.V),
		order:    make([]int, g.V),
	}
}

func (ts *TopologicalSort) Sort() bool {
	for v := 0; v < ts.g.V; v++ {
		ts.indegree[v] = ts.g.Indegree(v)
	}

	queue := make([]int, 0)
	for v := 0; v < ts.g.V; v++ {
		if ts.indegree[v] == 0 {
			queue = append(queue, v)
		}
	}

	count := 0
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		ts.order[count] = u
		count++
		for _, v := range ts.g.Adj(u) {
			ts.indegree[v]--
			if ts.indegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	if count != ts.g.V {
		ts.order = nil
		return false
	}
	return true
}

func (ts *TopologicalSort) HasOrder() bool {
	return ts.order != nil
}

func (ts *TopologicalSort) GetOrders() []int {
	return ts.order
}

func (ts *TopologicalSort) String() string {
	if ts.order == nil {
		return "null"
	}
	s := "["
	for i, v := range ts.order {
		s += fmt.Sprintf("%d", v)
		if i < len(ts.order)-1 {
			s += ","
		}
	}
	s += "]"
	return s
}
