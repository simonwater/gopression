package ir

import (
	"fmt"
	"strings"

	"github.com/simonwater/gopression/util"
)

type ExecuteContext struct {
	exprInfos []*ExprInfo
	nodeSet   *util.NodeSet[*ExprInfo]
	graph     *util.Digraph
	global    *GopContext
}

func NewExecuteContext(global *GopContext) *ExecuteContext {
	return &ExecuteContext{
		global: global,
	}
}

func (ec *ExecuteContext) GetExprInfos() []*ExprInfo {
	return ec.exprInfos
}

func (ec *ExecuteContext) GetNodeSet() *util.NodeSet[*ExprInfo] {
	return ec.nodeSet
}

func (ec *ExecuteContext) GetGraph() *util.Digraph {
	return ec.graph
}

func (ec *ExecuteContext) HasAssign() bool {
	return ec.graph != nil && ec.graph.V > 0
}

func (ec *ExecuteContext) PreExecute(exprInfos []*ExprInfo) {
	ec.nodeSet = util.NewNodeSet[*ExprInfo]()
	ec.exprInfos = exprInfos
	ec.initNodes()
	ec.initGraph()
}

func (ec *ExecuteContext) initNodes() {
	tracer := ec.global.GetTracer()
	tracer.StartTimer()

	for _, exprInfo := range ec.exprInfos {
		if !exprInfo.IsAssign() { // 只对赋值表达式构造有向图
			continue
		}

		// 添加前置节点
		for name := range exprInfo.GetPrecursors() {
			ec.nodeSet.AddNode(name)
		}

		// 添加后继节点并关联表达式信息
		first := true
		for name := range exprInfo.GetSuccessors() {
			node := ec.nodeSet.AddNode(name)
			if first {
				node.Info = exprInfo
				first = false
			}
		}
	}

	tracer.EndTimer("完成图节点初始化。")
}

func (ec *ExecuteContext) initGraph() {
	tracer := ec.global.GetTracer()
	tracer.StartTimer()

	if ec.nodeSet.Size() == 0 {
		ec.graph = util.NewDigraph(0)
		tracer.EndTimer("空图无需构造。")
		return
	}

	ec.graph = util.NewDigraph(ec.nodeSet.Size())

	for _, info := range ec.exprInfos {
		if !info.IsAssign() { // 只对赋值表达式构造有向图
			continue
		}

		for prec := range info.GetPrecursors() {
			preNode := ec.nodeSet.GetNodeByName(prec)
			u := preNode.Index

			for succ := range info.GetSuccessors() {
				succNode := ec.nodeSet.GetNodeByName(succ)
				v := succNode.Index
				ec.graph.AddEdge(u, v)
			}
		}
	}

	tracer.EndTimer("完成图的构造。")
}

func (ec *ExecuteContext) PrintGraph() string {
	if ec.graph == nil {
		return "图未初始化\n"
	}

	var builder strings.Builder
	vertices := ec.graph.V
	edges := ec.graph.E

	builder.WriteString(fmt.Sprintf("%d vertices, %d edges\n", vertices, edges))

	for u := 0; u < vertices; u++ {
		pre := ec.nodeSet.GetNodeByIndex(u)
		indegree := ec.graph.Indegree(u)

		builder.WriteString(fmt.Sprintf("%d(%s-%d): ", u, pre.Name, indegree))

		for _, v := range ec.graph.Adj(u) {
			succ := ec.nodeSet.GetNodeByIndex(v)
			builder.WriteString(fmt.Sprintf("%d(%s) ", v, succ.Name))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}
