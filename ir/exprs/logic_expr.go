package exprs

import "github.com/simonwater/gopression/values"

// LogicExpr 逻辑表达式
type LogicExpr struct {
	Left     Expr
	Operator *values.Token
	Right    Expr
}

func NewLogicExpr(left Expr, operator *values.Token, right Expr) *LogicExpr {
	return &LogicExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (e *LogicExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitLogic(e)
}
