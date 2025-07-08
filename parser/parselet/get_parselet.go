package parselet

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

// GetParselet 实现属性访问表达式解析
type GetParselet struct {
	Precedence int
}

// NewGetParselet 创建属性访问解析器
func NewGetParselet(precedence int) *GetParselet {
	return &GetParselet{Precedence: precedence}
}

// Parse 解析属性访问表达式
func (gp *GetParselet) Parse(p IParser, lhs exprs.Expr, token values.Token) exprs.Expr {
	// 确保点操作符后是标识符
	name := p.Consume(values.IDENTIFIER, "点操作符后必须是属性名")
	return exprs.NewGetExpr(lhs, &name)
}

// GetPrecedence 获取优先级
func (gp *GetParselet) GetPrecedence() int {
	return gp.Precedence
}
