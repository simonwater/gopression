package exprs

import (
	"github.com/simonwater/gopression/values"
)

// AssignExpr 赋值表达式
type AssignExpr struct {
	Left     Expr
	Operator *values.Token
	Right    Expr
}

func NewAssignExpr(left Expr, operator *values.Token, right Expr) *AssignExpr {
	return &AssignExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}
