package parselet

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

type GroupParselet struct{}

// NewGroupParselet 创建分组解析器
func NewGroupParselet() *GroupParselet {
	return &GroupParselet{}
}

// Parse 解析分组表达式
func (gp *GroupParselet) Parse(p IParser, token values.Token) exprs.Expr {
	// 解析括号内的表达式
	expr := p.ExpressionPrec(0)

	// 确保右括号存在
	p.Consume(values.RIGHT_PAREN, "期望在表达式后出现 ')'")

	return expr
}
