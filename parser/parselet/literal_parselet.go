package parselet

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

type LiteralParselet struct{}

// NewLiteralParselet 创建字面量解析器
func NewLiteralParselet() *LiteralParselet {
	return &LiteralParselet{}
}

// Parse 解析字面量表达式
func (lp *LiteralParselet) Parse(p IParser, token values.Token) exprs.Expr {
	return exprs.NewLiteralExpr(&token.Literal)
}
