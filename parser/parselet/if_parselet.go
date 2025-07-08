package parselet

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

type IfParselet struct{}

func NewIfParselet() *IfParselet {
	return &IfParselet{}
}

// Parse 解析 if 表达式，函数形式的if
func (ip *IfParselet) Parse(p IParser, token values.Token) exprs.Expr {
	p.Consume(values.LEFT_PAREN, "if 后期望 '('")
	condition := p.ExpressionPrec(0)
	p.Consume(values.COMMA, "条件表达式后期望 ','")
	thenExpr := p.ExpressionPrec(0)
	var elseExpr exprs.Expr = nil
	if p.Match(values.COMMA) {
		elseExpr = p.ExpressionPrec(0)
	}
	p.Consume(values.RIGHT_PAREN, "if 表达式结尾期望 ')'")

	return exprs.NewIfExpr(condition, thenExpr, elseExpr)
}
