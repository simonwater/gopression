package exprs_test

import (
	"testing"

	"github.com/simonwater/gopression/env"
	"github.com/simonwater/gopression/gop"
)

func TestLogicExpressions(t *testing.T) {
	// 创建环境并设置变量
	ev := env.NewDefaultEnvironment()
	ev.PutInt("a", 1)
	ev.PutInt("b", 2)
	ev.PutInt("c", 3)

	runner := gop.NewGopRunner()

	// 测试用例表
	tests := []struct {
		expr   string
		expect bool
	}{
		{"a == 1 || b == 0 || c == 0", true},
		{"a == 0 || b == 2 || c == 0", true},
		{"a == 0 || b == 0 || c == 3", true},
		{"a == 0 || b == 0 || c == 0", false},
		{"a == 1 && b == 2 && c == 3", true},
		{"a == 0 && b == 2 && c == 3", false},
		{"a == 1 && b == 0 && c == 3", false},
		{"a == 1 && b == 2 && c == 0", false},
	}

	// 遍历执行所有测试用例
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := runner.Execute(tt.expr, ev)
			if err != nil {
				t.Fatalf("执行表达式 '%s' 出错: %v", tt.expr, err)
			}

			if result != tt.expect {
				t.Errorf("表达式 '%s' 错误\n期望: %t\n实际: %t",
					tt.expr, tt.expect, result)
			}
		})
	}
}
