package gop_test

import (
	"testing"

	"github.com/simonwater/gopression/env"
	"github.com/simonwater/gopression/gop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试基础数值表达式
func TestBasicNumericalExpressions(t *testing.T) {
	runner := gop.NewGopRunner()

	tests := []struct {
		expr   string
		expect interface{}
	}{
		{"1 + 2 - 3", 0},
		{"1 + 2 * 3", 7},
		{"3 * (2 + 1)", 9},
		{"1 + 2 * 3 ** 2 ** 1", 19.0},
		{"3 * (2 + 1.0)", 9.0},
		{"3 * (2 + 1.0) > 7", true},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := runner.Execute(tt.expr, nil)
			require.NoError(t, err, "执行表达式出错")
			assert.Equal(t, tt.expect, result, "表达式结果不匹配")
		})
	}
}

// 测试带变量的赋值表达式
func TestAssignmentWithVariables(t *testing.T) {
	ev := env.NewDefaultEnvironment()
	runner := gop.NewGopRunner()

	// 设置初始变量
	ev.PutInt("a", 1)
	ev.PutInt("b", 2)
	ev.PutInt("c", 3)

	// 执行赋值表达式
	result, err := runner.Execute("x = y = a + b * c", ev)
	require.NoError(t, err, "执行表达式出错")
	assert.Equal(t, 7, result, "返回值不匹配")

	// 验证变量值
	x := ev.Get("x").GetValue()
	assert.Equal(t, 7, x, "变量 x 值不匹配")

	y := ev.Get("y").GetValue()
	assert.Equal(t, 7, y, "变量 y 值不匹配")
}

// 测试带变量的求值表达式
func TestEvaluationWithVariables(t *testing.T) {
	ev := env.NewDefaultEnvironment()
	runner := gop.NewGopRunner()

	// 设置初始变量
	ev.PutInt("a", 1)
	ev.PutInt("b", 2)
	ev.PutInt("c", 3)

	tests := []struct {
		expr   string
		expect interface{}
	}{
		{"a + b * c - 100 / 5 ** 2 ** 1", 3.0},
		{"a + b * c >= 6", true},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := runner.Execute(tt.expr, ev)
			require.NoError(t, err, "执行表达式出错")
			assert.Equal(t, tt.expect, result, "表达式结果不匹配")
		})
	}
}

// 测试批量执行表达式
func TestBatchEvaluation(t *testing.T) {
	ev := env.NewDefaultEnvironment()
	runner := gop.NewGopRunner()

	// 设置初始变量
	ev.PutInt("a", 1)
	ev.PutInt("b", 2)
	ev.PutInt("c", 3)

	// 批量表达式
	expressions := []string{
		"a + b * c - 100 / 5 ** 2 ** 1",
		"a + b * c >= 6",
		"1 + 2 - 3",
		"3 * (2 + 1)",
		"a + (b - c)",
		"a * 2 + (b - c)",
		"x = y = a + b * c",
	}
	// 验证结果
	expectedResults := []interface{}{
		3.0,
		true,
		0,
		9,
		0,
		1,
		7,
	}

	// 执行批量操作
	results, err := runner.ExecuteBatch(expressions, ev)
	require.NoError(t, err, "批量执行出错")
	require.Len(t, results, len(expressions), "结果数量不匹配")

	for i, expr := range expressions {
		assert.Equal(t, expectedResults[i], results[i],
			"表达式 '%s' 结果不匹配", expr)
	}

	// 验证变量赋值
	x := ev.Get("x").GetValue()
	assert.Equal(t, 7, x, "变量 x 值不匹配")

	y := ev.Get("y").GetValue()
	assert.Equal(t, 7, y, "变量 y 值不匹配")
}

// 测试批量赋值表达式
func TestBatchAssignment(t *testing.T) {
	ev := env.NewDefaultEnvironment()
	runner := gop.NewGopRunner()

	// 设置初始变量
	ev.PutInt("m", 2)
	ev.PutInt("n", 4)
	ev.PutInt("w", 6)

	// 批量赋值表达式
	assignments := []string{
		"x = a + b * c",
		"a = m + n",
		"b = a * 2",
		"c = n + w",
	}
	// 验证返回结果
	expectedResults := []interface{}{
		126, // x = a + b * c = 6 + 12*10 = 6+120=126
		6,   // a = m + n = 2+4=6
		12,  // b = a * 2 = 6*2=12
		10,  // c = n + w = 4+6=10
	}

	// 执行批量操作
	results, err := runner.ExecuteBatch(assignments, ev)
	require.NoError(t, err, "批量赋值出错")
	require.Len(t, results, len(assignments), "结果数量不匹配")

	for i, expr := range assignments {
		assert.Equal(t, expectedResults[i], results[i],
			"表达式 '%s' 返回值不匹配", expr)
	}

	// 验证环境变量
	testCases := []struct {
		name  string
		value interface{}
	}{
		{"x", 126},
		{"a", 6},
		{"b", 12},
		{"c", 10},
	}

	for _, tc := range testCases {
		val := ev.Get(tc.name).GetValue()
		assert.Equal(t, tc.value, val, "变量 %s 值不匹配", tc.name)
	}
}
