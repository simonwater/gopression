package exprs

import "github.com/simonwater/gopression/values"

// BinaryExpr 二元表达式
type BinaryExpr struct {
	Left     Expr
	Operator *values.Token
	Right    Expr
}

func NewBinaryExpr(left Expr, operator *values.Token, right Expr) *BinaryExpr {
	return &BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}
