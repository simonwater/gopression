package ir

import (
	"errors"

	"github.com/simonwater/gopression/util"
)

type ExprSorter struct {
	nodeSet *util.NodeSet[*ExprInfo]
	graph   *util.Digraph
	context *GopContext
}

func NewExprSorter(context *GopContext) *ExprSorter {
	execContext := context.GetExecContext()
	return &ExprSorter{
		context: context,
		nodeSet: execContext.GetNodeSet(),
		graph:   execContext.GetGraph(),
	}
}

func (es *ExprSorter) Sort() ([]*ExprInfo, error) {
	if es.graph == nil || es.graph.V == 0 {
		return nil, nil
	}

	tracer := es.context.GetTracer()
	tracer.StartTimer()

	topSorter := util.NewTopologicalSort(es.graph)

	if !topSorter.Sort() {
		return nil, errors.New("公式列表存在循环引用！")
	}

	nodeOrders := topSorter.GetOrders()
	result := make([]*ExprInfo, 0, len(nodeOrders))

	for _, nodeIndex := range nodeOrders {
		node := es.nodeSet.GetNodeByIndex(nodeIndex)
		if node.Info != nil {
			result = append(result, node.Info)
		}
	}

	origInfos := es.context.GetExecContext().GetExprInfos()
	for _, expr := range origInfos {
		if !expr.IsAssign() {
			result = append(result, expr)
		}
	}

	tracer.EndTimer("完成拓扑排序。")
	return result, nil
}

// 可选：打印循环依赖的方法
func (es *ExprSorter) PrintCircle() {
	// 实现循环依赖检测和打印逻辑
	// 取决于 TopologicalSort 是否提供相关功能
}
