package ir_test

import (
	"testing"

	"github.com/simonwater/gopression/env"
	"github.com/simonwater/gopression/gop"
	"github.com/simonwater/gopression/ir"
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExprSorter_ShouldSortFormulas(t *testing.T) {
	srcs := []string{
		"x = y = a + b * c",
		"a = m + n",
		"b = a * 2",
		"c = n + w + b",
	}

	context := ir.NewGopContext()
	exprs := parse(srcs, context)
	sortedExprInfos, err := analyze(exprs, context)

	require.NoError(t, err, "排序不应出错")
	require.NotNil(t, sortedExprInfos, "表达式信息不应为nil")
	require.Len(t, sortedExprInfos, 4, "应有4个表达式")

	assert.Equal(t, "a = m + n", srcs[sortedExprInfos[0].GetIndex()])
	assert.Equal(t, "b = a * 2", srcs[sortedExprInfos[1].GetIndex()])
	assert.Equal(t, "c = n + w + b", srcs[sortedExprInfos[2].GetIndex()])
	assert.Equal(t, "x = y = a + b * c", srcs[sortedExprInfos[3].GetIndex()])

	runner := gop.NewGopRunner()
	environment := env.NewDefaultEnvironment()
	environment.PutInt("m", 2)
	environment.PutInt("n", 4)
	environment.PutInt("w", 6)
	runner.ExecuteBatch(srcs, environment)
	assert.Equal(t, 270, environment.Get("x").GetValue())
	assert.Equal(t, 270, environment.Get("y").GetValue())
	assert.Equal(t, 6, environment.Get("a").GetValue())
	assert.Equal(t, 12, environment.Get("b").GetValue())
	assert.Equal(t, 22, environment.Get("c").GetValue())

}

func TestExprSorter_ShouldSortMixedFormulas(t *testing.T) {
	srcs := []string{
		"b * 2 + 1",
		"a * b + c",
		"x = y = a + b * c",
		"a = m + n",
		"b = a * 2",
		"c = n + w + b",
	}

	context := ir.NewGopContext()
	exprs := parse(srcs, context)
	sortedExprInfos, err := analyze(exprs, context)
	require.NoError(t, err, "排序不应出错")
	require.NotNil(t, sortedExprInfos, "表达式信息不应为nil")
	require.Len(t, sortedExprInfos, 6, "应有6个表达式")

	assert.Equal(t, "a = m + n", srcs[sortedExprInfos[0].GetIndex()])
	assert.Equal(t, "b = a * 2", srcs[sortedExprInfos[1].GetIndex()])
	assert.Equal(t, "c = n + w + b", srcs[sortedExprInfos[2].GetIndex()])
	assert.Equal(t, "x = y = a + b * c", srcs[sortedExprInfos[3].GetIndex()])
	assert.Equal(t, "b * 2 + 1", srcs[sortedExprInfos[4].GetIndex()])
	assert.Equal(t, "a * b + c", srcs[sortedExprInfos[5].GetIndex()])

	runner := gop.NewGopRunner()
	environment := env.NewDefaultEnvironment()
	environment.PutInt("m", 2)
	environment.PutInt("n", 4)
	environment.PutInt("w", 6)
	result := runner.ExecuteBatch(srcs, environment)
	assert.NotNil(t, result)
	assert.Equal(t, 6, len(result))
	assert.Equal(t, 270, environment.Get("x").GetValue())
	assert.Equal(t, 270, environment.Get("y").GetValue())
	assert.Equal(t, 6, environment.Get("a").GetValue())
	assert.Equal(t, 12, environment.Get("b").GetValue())
	assert.Equal(t, 22, environment.Get("c").GetValue())

	assert.Equal(t, 12*2+1, result[0])
	assert.Equal(t, 6*12+22, result[1])
	assert.Equal(t, 270, result[2])
	assert.Equal(t, 6, result[3])
	assert.Equal(t, 12, result[4])
	assert.Equal(t, 22, result[5])
}

func parse(srcs []string, context *ir.GopContext) []exprs.Expr {
	tracer := context.GetTracer()
	tracer.StartTimerWithMsg("解析")
	result := make([]exprs.Expr, 0, len(srcs))
	for _, src := range srcs {
		parser := parser.NewParser(src)
		expr := parser.Parse()
		result = append(result, expr)
	}
	tracer.EndTimer("完成表达式解析。")
	return result
}

func analyze(exprs []exprs.Expr, context *ir.GopContext) ([]*ir.ExprInfo, error) {
	tracer := context.GetTracer()
	tracer.StartTimerWithMsg("分析")

	exprInfos := make([]*ir.ExprInfo, len(exprs))
	for i, expr := range exprs {
		exprInfos[i] = ir.NewExprInfo(expr, i)
	}

	context.PrepareExecute(exprInfos)
	sorter := ir.NewExprSorter(context)
	sortedInfos, err := sorter.Sort()

	tracer.EndTimer("完成表达式分析。")
	return sortedInfos, err
}
