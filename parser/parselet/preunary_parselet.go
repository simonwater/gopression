package parselet

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

type PreUnaryParselet struct {
	Precedence int
}

func NewPreUnaryParselet(precedence int) *PreUnaryParselet {
	return &PreUnaryParselet{
		Precedence: precedence,
	}
}

// Parse 解析前缀一元表达式
func (pup *PreUnaryParselet) Parse(parser IParser, token values.Token) exprs.Expr {
	rhs := parser.ExpressionPrec(pup.Precedence)
	return exprs.NewUnaryExpr(&token, rhs)
}
