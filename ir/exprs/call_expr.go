package exprs

import "github.com/simonwater/gopression/values"

// CallExpr 函数调用表达式
type CallExpr struct {
	Callee Expr
	Args   []Expr
	RParen *values.Token
}

func NewCallExpr(callee Expr, args []Expr, rParen *values.Token) *CallExpr {
	return &CallExpr{
		Callee: callee,
		Args:   args,
		RParen: rParen,
	}
}
