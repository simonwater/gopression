package exprs

import "github.com/simonwater/gopression/values"

// SetExpr 属性设置表达式
type SetExpr struct {
	Object Expr
	Name   *values.Token
	Value  Expr
}

func NewSetExpr(object Expr, name *values.Token, value Expr) *SetExpr {
	return &SetExpr{
		Object: object,
		Name:   name,
		Value:  value,
	}
}

func (e *SetExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitSet(e)
}
