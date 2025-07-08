package exprs

import "github.com/simonwater/gopression/values"

// LiteralExpr 字面量表达式
type LiteralExpr struct {
	Value *values.Value
}

func NewLiteralExpr(value *values.Value) *LiteralExpr {
	return &LiteralExpr{Value: value}
}

func (e *LiteralExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteral(e)
}

func (e *LiteralExpr) String() string {
	return e.Value.String()
}
