package exprs

import "github.com/simonwater/gopression/values"

// UnaryExpr 一元表达式
type UnaryExpr struct {
	Operator *values.Token
	Right    Expr
}

func NewUnaryExpr(operator *values.Token, right Expr) *UnaryExpr {
	return &UnaryExpr{
		Operator: operator,
		Right:    right,
	}
}
