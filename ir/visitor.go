package ir

import "github.com/simonwater/gopression/ir/exprs"

// Visitor 访问者接口
type Visitor[T any] interface {
	VisitBinary(expr *exprs.BinaryExpr) T
	VisitLogic(expr *exprs.LogicExpr) T
	VisitLiteral(expr *exprs.LiteralExpr) T
	VisitUnary(expr *exprs.UnaryExpr) T
	VisitId(expr *exprs.IdExpr) T
	VisitAssign(expr *exprs.AssignExpr) T
	VisitCall(expr *exprs.CallExpr) T
	VisitIf(expr *exprs.IfExpr) T
	VisitGet(expr *exprs.GetExpr) T
	VisitSet(expr *exprs.SetExpr) T
}

type BaseVisitor[T any] struct {
	Visitor[T]
}

func NewBaseVisitor[T any](visitor Visitor[T]) *BaseVisitor[T] {
	return &BaseVisitor[T]{
		Visitor: visitor,
	}
}

func (bv *BaseVisitor[T]) Accept(expr exprs.Expr) T {
	switch t := expr.(type) {
	case *exprs.BinaryExpr:
		return bv.VisitBinary(t)
	case *exprs.LogicExpr:
		return bv.VisitLogic(t)
	case *exprs.LiteralExpr:
		return bv.VisitLiteral(t)
	case *exprs.UnaryExpr:
		return bv.VisitUnary(t)
	case *exprs.IdExpr:
		return bv.VisitId(t)
	case *exprs.AssignExpr:
		return bv.VisitAssign(t)
	case *exprs.CallExpr:
		return bv.VisitCall(t)
	case *exprs.IfExpr:
		return bv.VisitIf(t)
	case *exprs.GetExpr:
		return bv.VisitGet(t)
	case *exprs.SetExpr:
		return bv.VisitSet(t)
	default:
		panic("类型尚未支持！")
	}
}
