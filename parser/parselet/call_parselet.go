package parselet

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

// CallParselet 实现函数调用表达式解析
type CallParselet struct {
	Precedence int
}

// NewCallParselet 创建函数调用解析器
func NewCallParselet(precedence int) *CallParselet {
	return &CallParselet{Precedence: precedence}
}

// Parse 解析函数调用表达式
func (cp *CallParselet) Parse(p IParser, lhs exprs.Expr, token values.Token) exprs.Expr {
	args := make([]exprs.Expr, 0)
	argCount := 0

	// 检查是否有参数（下一个token不是右括号）
	if !p.Check(values.RIGHT_PAREN) {
		for {
			// 检查参数数量限制
			if argCount >= 255 {
				panic(NewLoxParseError(p.Peek(), "参数数量不能超过255个"))
			}

			// 解析参数表达式
			arg := p.ExpressionPrec(0)
			args = append(args, arg)
			argCount++

			// 如果没有逗号，结束参数解析
			if !p.Match(values.COMMA) {
				break
			}
		}
	}

	// 确保右括号存在
	rParen := p.Consume(values.RIGHT_PAREN, "期望在参数后出现 ')'")
	return exprs.NewCallExpr(lhs, args, &rParen)
}

// GetPrecedence 获取优先级
func (cp *CallParselet) GetPrecedence() int {
	return cp.Precedence
}
