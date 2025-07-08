package util

import (
	"fmt"
	"strings"
)

type Digraph struct {
	V        int     // number of vertices
	E        int     // number of edges
	adj      [][]int // adjacency list for each vertex
	indegree []int   // indegree of each vertex
}

func NewDigraph(V int) *Digraph {
	if V < 0 {
		panic("Number of vertices in a Digraph must be non-negative")
	}
	adj := make([][]int, V)
	for i := 0; i < V; i++ {
		adj[i] = make([]int, 0)
	}
	return &Digraph{
		V:        V,
		E:        0,
		adj:      adj,
		indegree: make([]int, V),
	}
}

func (g *Digraph) AddEdge(v, w int) {
	g.validateVertex(v)
	g.validateVertex(w)
	g.adj[v] = append(g.adj[v], w)
	g.indegree[w]++
	g.E++
}

func (g *Digraph) Adj(v int) []int {
	g.validateVertex(v)
	return g.adj[v]
}

func (g *Digraph) Outdegree(v int) int {
	g.validateVertex(v)
	return len(g.adj[v])
}

func (g *Digraph) Indegree(v int) int {
	g.validateVertex(v)
	return g.indegree[v]
}

func (g *Digraph) Reverse() *Digraph {
	reverse := NewDigraph(g.V)
	for v := 0; v < g.V; v++ {
		for _, w := range g.adj[v] {
			reverse.AddEdge(w, v)
		}
	}
	return reverse
}

func (g *Digraph) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d vertices, %d edges \n", g.V, g.E))
	for v := 0; v < g.V; v++ {
		sb.WriteString(fmt.Sprintf("%d: ", v))
		for _, w := range g.adj[v] {
			sb.WriteString(fmt.Sprintf("%d ", w))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (g *Digraph) validateVertex(v int) {
	if v < 0 || v >= g.V {
		panic(fmt.Sprintf("vertex %d is not between 0 and %d", v, g.V-1))
	}
}
