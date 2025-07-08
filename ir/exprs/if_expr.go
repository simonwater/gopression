package exprs

// IfExpr 条件表达式
type IfExpr struct {
	Condition  Expr
	ThenBranch Expr
	ElseBranch Expr
}

func NewIfExpr(condition Expr, thenBranch Expr, elseBranch Expr) *IfExpr {
	return &IfExpr{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (e *IfExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitIf(e)
}
