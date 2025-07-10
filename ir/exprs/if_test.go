package exprs_test

import (
	"testing"

	"github.com/simonwater/gopression/env"
	"github.com/simonwater/gopression/gop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIfExpression 测试 if 表达式的各种情况
func TestIfExpression(t *testing.T) {

	t.Run("基本 if 表达式", func(t *testing.T) {
		env := env.NewDefaultEnvironment()
		runner := gop.NewGopRunner()
		// 设置环境变量
		env.PutInt("a", 1)
		env.PutInt("b", 2)
		env.PutInt("c", 3)

		// 测试条件为真
		result, err := runner.Execute("if(a + b * c >= 6, 6 ** 2, -6 * 2)", env)
		require.NoError(t, err)
		assert.Equal(t, 36.0, result)

		// 测试条件为假
		result, err = runner.Execute("if(a + b * c < 6, 6 ** 2, -6 * 2)", env)
		require.NoError(t, err)
		assert.Equal(t, -12, result)

		// 测试只有真分支
		result, err = runner.Execute("if(a + b * c < 6, 6 ** 2)", env)
		require.NoError(t, err)
		assert.Nil(t, result)

		// 测试只有真分支且条件为真
		result, err = runner.Execute("if(a + b * c >= 6, 6 ** 2)", env)
		require.NoError(t, err)
		assert.Equal(t, 36.0, result)
	})

	t.Run("嵌套 if 表达式", func(t *testing.T) {
		env := env.NewDefaultEnvironment()
		runner := gop.NewGopRunner()
		// 定义嵌套表达式
		str1 := "if(score >= 85, \"A\", if(score >= 70, \"B\", if(score >= 60, \"C\", \"D\")))"
		str2 := "if(score >= 70, if(score < 85, \"B\",\"A\"), if(score >= 60, \"C\", \"D\"))"

		// 测试分数 90
		env.PutInt("score", 90)
		result, err := runner.Execute(str1, env)
		require.NoError(t, err)
		assert.Equal(t, "A", result)

		result, err = runner.Execute(str2, env)
		require.NoError(t, err)
		assert.Equal(t, "A", result)

		// 测试分数 65
		env.PutInt("score", 65)
		result, err = runner.Execute(str1, env)
		require.NoError(t, err)
		assert.Equal(t, "C", result)

		result, err = runner.Execute(str2, env)
		require.NoError(t, err)
		assert.Equal(t, "C", result)

		// 测试分数 50
		env.PutInt("score", 50)
		result, err = runner.Execute(str1, env)
		require.NoError(t, err)
		assert.Equal(t, "D", result)

		result, err = runner.Execute(str2, env)
		require.NoError(t, err)
		assert.Equal(t, "D", result)
	})

	t.Run("单分支 if 表达式", func(t *testing.T) {
		env := env.NewDefaultEnvironment()
		runner := gop.NewGopRunner()
		// 设置环境变量
		env.PutInt("x1", 0)
		env.PutInt("y1", 0)
		env.PutInt("x2", 0)
		env.PutInt("y2", 0)

		// 测试条件为真
		result, err := runner.Execute("if(1 == 1, x1 = 1, y1 = 2)", env)
		require.NoError(t, err)
		assert.Equal(t, 1, result)

		// 验证环境变量
		x1 := env.Get("x1")
		assert.Equal(t, 1, x1.GetValue())

		y1 := env.Get("y1")
		assert.Equal(t, 0, y1.GetValue())

		// 测试条件为假
		result, err = runner.Execute("if(1 != 1, x2 = 1, y2 = 2)", env)
		require.NoError(t, err)
		assert.Equal(t, 2, result)

		// 验证环境变量
		x2 := env.Get("x2")
		assert.Equal(t, 0, x2.GetValue())

		y2 := env.Get("y2")
		assert.Equal(t, 2, y2.GetValue())
	})

	t.Run("错误处理", func(t *testing.T) {
		env := env.NewDefaultEnvironment()
		runner := gop.NewGopRunner()
		// 设置环境变量
		env.PutInt("a", 1)
		env.PutInt("b", 2)
		env.PutInt("c", 3)

		// 测试各种错误情况
		tests := []struct {
			name      string
			expr      string
			expectErr string
		}{
			{"空 if 表达式", "if()", "if 表达式需要至少 2 个参数"},
			{"缺少分支", "if(a + b * c >= 6)", "if 表达式需要至少 2 个参数"},
			{"逗号错误", "if(a + b * c >= 6,)", "解析错误"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := runner.Execute(tt.expr, env)
				require.NotNil(t, err)
			})
		}
	})
}
