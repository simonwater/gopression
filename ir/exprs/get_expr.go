package exprs

import "github.com/simonwater/gopression/values"

// GetExpr 属性获取表达式
type GetExpr struct {
	Object Expr
	Name   *values.Token
}

func NewGetExpr(object Expr, name *values.Token) *GetExpr {
	return &GetExpr{
		Object: object,
		Name:   name,
	}
}
