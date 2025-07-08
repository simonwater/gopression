package ir

import "github.com/simonwater/gopression/ir/exprs"

type ExprInfo struct {
	precursors map[string]bool // 依赖的变量 read
	successors map[string]bool // 被赋值的变量 write
	expr       exprs.Expr
	index      int
}

func NewExprInfo(e exprs.Expr, idx int) *ExprInfo {
	ei := &ExprInfo{
		precursors: make(map[string]bool),
		successors: make(map[string]bool),
		expr:       e,
		index:      idx,
	}
	ei.initVariables()
	return ei
}

func (ei *ExprInfo) initVariables() {
	varQuery := NewVarsQuery()
	varSet := varQuery.Execute(ei.expr)
	if varSet != nil {
		ei.precursors = varSet.GetDepends()
		ei.successors = varSet.GetAssigns()
	}
}

func (ei *ExprInfo) IsAssign() bool {
	_, isAssign := ei.expr.(*exprs.AssignExpr)
	_, isSet := ei.expr.(*exprs.SetExpr)
	return isAssign || isSet
}

// Getters and Setters
func (ei *ExprInfo) GetPrecursors() map[string]bool {
	return ei.precursors
}

func (ei *ExprInfo) SetPrecursors(precursors map[string]bool) {
	ei.precursors = precursors
}

func (ei *ExprInfo) GetSuccessors() map[string]bool {
	return ei.successors
}

func (ei *ExprInfo) SetSuccessors(successors map[string]bool) {
	ei.successors = successors
}

func (ei *ExprInfo) GetExpr() exprs.Expr {
	return ei.expr
}

func (ei *ExprInfo) SetExpr(e exprs.Expr) {
	ei.expr = e
}

func (ei *ExprInfo) GetIndex() int {
	return ei.index
}

func (ei *ExprInfo) SetIndex(idx int) {
	ei.index = idx
}
