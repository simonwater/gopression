package parselet

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

// BinaryParselet 实现二元表达式解析
type BinaryParselet struct {
	Precedence int
	IsRight    bool // 是否右结合
}

// NewBinaryParselet 创建二元表达式解析器（默认左结合）
func NewBinaryParselet(precedence int) *BinaryParselet {
	return &BinaryParselet{
		Precedence: precedence,
		IsRight:    false,
	}
}

// NewRightAssociativeBinaryParselet 创建右结合二元表达式解析器
func NewRightCombineBinaryParselet(precedence int) *BinaryParselet {
	return &BinaryParselet{
		Precedence: precedence,
		IsRight:    true,
	}
}

// Parse 解析二元表达式
func (bp *BinaryParselet) Parse(p IParser, lhs exprs.Expr, token values.Token) exprs.Expr {
	// 根据结合性确定优先级
	nextPrecedence := bp.Precedence
	if bp.IsRight {
		nextPrecedence--
	}

	// 解析右侧表达式
	rhs := p.ExpressionPrec(nextPrecedence)

	return exprs.NewBinaryExpr(lhs, &token, rhs)
}

// GetPrecedence 获取优先级
func (bp *BinaryParselet) GetPrecedence() int {
	return bp.Precedence
}
