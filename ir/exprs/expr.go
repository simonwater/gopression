package exprs

// Expr 表达式接口
type Expr interface {
}

// Visitor 访问者接口
type Visitor[T any] interface {
	VisitBinary(expr *BinaryExpr) T
	VisitLogic(expr *LogicExpr) T
	VisitLiteral(expr *LiteralExpr) T
	VisitUnary(expr *UnaryExpr) T
	VisitId(expr *IdExpr) T
	VisitAssign(expr *AssignExpr) T
	VisitCall(expr *CallExpr) T
	VisitIf(expr *IfExpr) T
	VisitGet(expr *GetExpr) T
	VisitSet(expr *SetExpr) T
}

func VisitExpr[T any](expr Expr, visitor Visitor[T]) T {
	switch t := expr.(type) {
	case *BinaryExpr:
		return visitor.VisitBinary(t)
	case *LogicExpr:
		return visitor.VisitLogic(t)
	case *LiteralExpr:
		return visitor.VisitLiteral(t)
	case *UnaryExpr:
		return visitor.VisitUnary(t)
	case *IdExpr:
		return visitor.VisitId(t)
	case *AssignExpr:
		return visitor.VisitAssign(t)
	case *CallExpr:
		return visitor.VisitCall(t)
	case *IfExpr:
		return visitor.VisitIf(t)
	case *GetExpr:
		return visitor.VisitGet(t)
	case *SetExpr:
		return visitor.VisitSet(t)
	default:
		panic("类型尚未支持！")
	}
}
