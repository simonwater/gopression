package parselet

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

// AssignParselet 实现赋值表达式解析
type AssignParselet struct {
	Precedence int
}

// NewAssignParselet 创建赋值解析器
func NewAssignParselet(precedence int) *AssignParselet {
	return &AssignParselet{Precedence: precedence}
}

// Parse 解析赋值表达式
func (ap *AssignParselet) Parse(p IParser, lhs exprs.Expr, token values.Token) exprs.Expr {
	// 右结合，优先级降低一位，支持连续赋值
	rhs := p.ExpressionPrec(ap.Precedence - 1)

	// 检查是否是属性赋值 (object.property = value)
	if getExpr, ok := lhs.(*exprs.GetExpr); ok {
		return exprs.NewSetExpr(getExpr.Object, getExpr.Name, rhs)
	}
	return exprs.NewAssignExpr(lhs, &token, rhs)
}

// GetPrecedence 获取优先级
func (ap *AssignParselet) GetPrecedence() int {
	return ap.Precedence
}
