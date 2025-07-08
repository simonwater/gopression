package parselet

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

type IdParselet struct{}

// NewIdParselet 创建标识符解析器
func NewIdParselet() *IdParselet {
	return &IdParselet{}
}

// Parse 解析标识符表达式
func (ip *IdParselet) Parse(p IParser, token values.Token) exprs.Expr {
	// 创建标识符表达式
	return exprs.NewIdExpr(token.Lexeme)
}
