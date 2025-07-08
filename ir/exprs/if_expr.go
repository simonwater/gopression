package exprs

// IfExpr 条件表达式
type IfExpr struct {
	Condition Expr
	Then      Expr
	Else      Expr
}

func NewIfExpr(condition Expr, thenBranch Expr, elseBranch Expr) *IfExpr {
	return &IfExpr{
		Condition: condition,
		Then:      thenBranch,
		Else:      elseBranch,
	}
}

func (e *IfExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitIf(e)
}
