package ir

import (
	"strings"

	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/parser"
)

type VarsQuery struct{}

func NewVarsQuery() *VarsQuery {
	return &VarsQuery{}
}

func (vq *VarsQuery) ExecuteAll(exprs []exprs.Expr) *VariableSet {
	if len(exprs) == 0 {
		return nil
	}

	result := NewVariableSet()
	for _, expr := range exprs {
		result.Combine(vq.Execute(expr))
	}
	return result
}

func (vq *VarsQuery) ExecuteSrc(src string) *VariableSet {
	if strings.TrimSpace(src) == "" {
		return nil
	}

	p := parser.NewParser(src)
	expr := p.Parse()
	return vq.Execute(expr)
}

func (vq *VarsQuery) Execute(expr exprs.Expr) *VariableSet {
	if expr == nil {
		return nil
	}
	result := expr.Accept(vq)
	if result == nil {
		return nil
	}
	return result.(*VariableSet)
}

func (vq *VarsQuery) VisitBinary(expr *exprs.BinaryExpr) interface{} {
	result := vq.Execute(expr.Left)
	rhs := vq.Execute(expr.Right)

	if result == nil {
		return rhs
	}
	if rhs != nil {
		result.Combine(rhs)
	}
	return result
}

func (vq *VarsQuery) VisitLogic(expr *exprs.LogicExpr) interface{} {
	result := vq.Execute(expr.Left)
	rhs := vq.Execute(expr.Right)

	if result == nil {
		return rhs
	}
	if rhs != nil {
		result.Combine(rhs)
	}
	return result
}

func (vq *VarsQuery) VisitLiteral(expr *exprs.LiteralExpr) interface{} {
	return nil
}

func (vq *VarsQuery) VisitUnary(expr *exprs.UnaryExpr) interface{} {
	return vq.Execute(expr.Right)
}

func (vq *VarsQuery) VisitId(expr *exprs.IdExpr) interface{} {
	return FromDepends(expr.Id)
}

func (vq *VarsQuery) VisitAssign(expr *exprs.AssignExpr) interface{} {
	// 处理标识符赋值
	if idExpr, ok := expr.Left.(*exprs.IdExpr); ok {
		result := FromAssigns(idExpr.Id)
		if rhs := vq.Execute(expr.Right); rhs != nil {
			result.Combine(rhs)
		}
		return result
	}

	// 处理其他类型的左值
	result := NewVariableSet()
	if leftVars := vq.Execute(expr.Left); leftVars != nil {
		result.SetAssigns(leftVars.GetDepends())
	}

	if rhs := vq.Execute(expr.Right); rhs != nil {
		result.Combine(rhs)
	}
	return result
}

func (vq *VarsQuery) VisitCall(expr *exprs.CallExpr) interface{} {
	result := NewVariableSet()
	for _, arg := range expr.Args {
		if argVars := vq.Execute(arg); argVars != nil {
			result.Combine(argVars)
		}
	}
	return result
}

func (vq *VarsQuery) VisitIf(expr *exprs.IfExpr) interface{} {
	result := NewVariableSet()

	if condVars := vq.Execute(expr.Condition); condVars != nil {
		result.Combine(condVars)
	}

	if thenVars := vq.Execute(expr.ThenBranch); thenVars != nil {
		result.Combine(thenVars)
	}

	if elseVars := vq.Execute(expr.ElseBranch); elseVars != nil {
		result.Combine(elseVars)
	}

	return result
}

func (vq *VarsQuery) VisitGet(expr *exprs.GetExpr) interface{} {
	names := []string{}
	vq.visitGetRecursive(expr, &names)
	id := strings.Join(names, ".")
	return FromDepends(id)
}

func (vq *VarsQuery) visitGetRecursive(e exprs.Expr, names *[]string) {
	switch expr := e.(type) {
	case *exprs.IdExpr:
		*names = append(*names, expr.Id)
	case *exprs.GetExpr:
		vq.visitGetRecursive(expr.Object, names)
		*names = append(*names, expr.Name.Lexeme)
	}
}

func (vq *VarsQuery) VisitSet(expr *exprs.SetExpr) interface{} {
	// 创建虚拟的GetExpr来获取属性路径
	getExpr := exprs.NewGetExpr(expr.Object, expr.Name)
	gets := vq.Execute(getExpr)

	result := NewVariableSet()
	if gets != nil {
		result.SetAssigns(gets.GetDepends())
	}

	if rhs := vq.Execute(expr.Value); rhs != nil {
		result.Combine(rhs)
	}

	return result
}
