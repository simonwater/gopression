package parselet

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

type LogicParselet struct {
	Precedence int
	IsRight    bool // 是否右结合
}

// NewLogicParselet 创建逻辑表达式解析器（默认左结合）
func NewLogicParselet(precedence int) *LogicParselet {
	return &LogicParselet{
		Precedence: precedence,
		IsRight:    false,
	}
}

// 创建右结合逻辑表达式解析器
func NewRightCombineLogicParselet(precedence int) *LogicParselet {
	return &LogicParselet{
		Precedence: precedence,
		IsRight:    true,
	}
}

// Parse 解析逻辑表达式
func (bp *LogicParselet) Parse(p IParser, lhs exprs.Expr, token values.Token) exprs.Expr {
	nextPrecedence := bp.Precedence
	if bp.IsRight {
		nextPrecedence--
	}
	rhs := p.ExpressionPrec(nextPrecedence)
	return exprs.NewLogicExpr(lhs, &token, rhs)
}

// GetPrecedence 获取优先级
func (bp *LogicParselet) GetPrecedence() int {
	return bp.Precedence
}
