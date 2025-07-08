package exprs

// Expr 表达式接口
type Expr interface {
	Accept(visitor Visitor) any
}

// Visitor 访问者接口
type Visitor interface {
	VisitBinary(expr *BinaryExpr) any
	VisitLogic(expr *LogicExpr) any
	VisitLiteral(expr *LiteralExpr) any
	VisitUnary(expr *UnaryExpr) any
	VisitId(expr *IdExpr) any
	VisitAssign(expr *AssignExpr) any
	VisitCall(expr *CallExpr) any
	VisitIf(expr *IfExpr) any
	VisitGet(expr *GetExpr) any
	VisitSet(expr *SetExpr) any
}
