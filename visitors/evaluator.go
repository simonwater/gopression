package visitors

import (
	"errors"
	"fmt"

	"github.com/simonwater/gopression/env"
	"github.com/simonwater/gopression/functions/funmgr"
	"github.com/simonwater/gopression/ir"
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/parser"
	"github.com/simonwater/gopression/values"
)

// Evaluator 表达式求值器
type Evaluator struct {
	*ir.BaseVisitor[values.Value]
	env env.Environment
}

func NewEvaluator(env env.Environment) *Evaluator {
	e := &Evaluator{env: env}
	e.BaseVisitor = ir.NewBaseVisitor(e)
	return e
}

func safeExecute(fn func() values.Value) (result values.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			default:
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()

	result = fn()
	return
}

func (e *Evaluator) ExecuteAll(exprs []exprs.Expr) ([]values.Value, error) {
	if len(exprs) == 0 {
		return nil, nil
	}

	results := make([]values.Value, len(exprs))
	for i, expr := range exprs {
		r, err := safeExecute(func() values.Value {
			return e.Execute(expr)
		})
		if err != nil {
			return nil, fmt.Errorf("error evaluating expression at index %d: %w", i, err)
		}
		results[i] = r
	}
	return results, nil
}

func (e *Evaluator) ExecuteSrc(src string) (values.Value, error) {
	if src == "" {
		return values.NewNullValue(), nil
	}
	return safeExecute(func() values.Value {
		p := parser.NewParser(src)
		expr := p.Parse()
		return e.Execute(expr)
	})

}

func (e *Evaluator) Execute(expr exprs.Expr) values.Value {
	if expr == nil {
		return values.NewNullValue()
	}
	return e.Accept(expr)
}

func (e *Evaluator) VisitBinary(expr *exprs.BinaryExpr) values.Value {
	left := e.Execute(expr.Left)
	right := e.Execute(expr.Right)
	r, err := values.BinaryOperate(left, right, expr.Operator.Type)
	if err != nil {
		panic(fmt.Errorf("error evaluating binary expression: %w", err))
	}
	return r
}

func (e *Evaluator) VisitLogic(expr *exprs.LogicExpr) values.Value {
	left := e.Execute(expr.Left)

	if expr.Operator.Type == values.OR {
		if left.IsTruthy() {
			return values.NewBooleanValue(true)
		}
	} else { // AND
		if !left.IsTruthy() {
			return values.NewBooleanValue(false)
		}
	}

	return e.Execute(expr.Right)
}

func (e *Evaluator) VisitLiteral(expr *exprs.LiteralExpr) values.Value {
	return *expr.Value
}

func (e *Evaluator) VisitUnary(expr *exprs.UnaryExpr) values.Value {
	right := e.Execute(expr.Right)
	r, err := values.PreUnaryOperate(right, expr.Operator.Type)
	if err != nil {
		panic(fmt.Errorf("error evaluating unary expression: %w", err))
	}
	return r
}

func (e *Evaluator) VisitId(expr *exprs.IdExpr) values.Value {
	return e.getVariableValue(expr.Id)
}

func (e *Evaluator) VisitAssign(expr *exprs.AssignExpr) values.Value {
	right := e.Execute(expr.Right)

	if idExpr, ok := expr.Left.(*exprs.IdExpr); ok {
		e.env.Put(idExpr.Id, right)
		return right
	}

	panic(errors.New("invalid assignment target"))
}

func (e *Evaluator) VisitCall(expr *exprs.CallExpr) values.Value {
	callee := expr.Callee

	var funcName string
	if idExpr, ok := callee.(*exprs.IdExpr); ok {
		funcName = idExpr.Id
	} else {
		panic(fmt.Errorf("can only call named functions"))
	}

	fn := funmgr.GetFunctionManager().GetFunction(funcName)
	if fn == nil {
		panic(fmt.Errorf("function not found: %s", funcName))
	}

	args := make([]values.Value, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = e.Execute(arg)
	}

	if len(args) != fn.Arity() {
		panic(fmt.Errorf(
			"expected %d arguments but got %d",
			fn.Arity(), len(args),
		))
	}

	r, err := fn.Call(args)
	if err != nil {
		panic(fmt.Errorf("error calling function %s: %w", funcName, err))
	}
	return r
}

func (e *Evaluator) VisitIf(expr *exprs.IfExpr) values.Value {
	cond := e.Execute(expr.Condition)
	if cond.IsTruthy() {
		return e.Execute(expr.ThenBranch)
	} else if expr.ElseBranch != nil {
		return e.Execute(expr.ElseBranch)
	}
	return values.NewNullValue()
}

func (e *Evaluator) VisitGet(expr *exprs.GetExpr) values.Value {
	object := e.Execute(expr.Object)
	if object.IsInstance() {
		obj := object.AsInstance()
		r, ok := obj.Get(expr.Name.Lexeme)
		if ok {
			return r
		} else {
			return values.NewNullValue()
		}
	}
	panic(errors.New("only instances have properties"))
}

func (e *Evaluator) VisitSet(expr *exprs.SetExpr) values.Value {
	object := e.Execute(expr.Object)
	if !object.IsInstance() {
		panic(errors.New("only instances have fields"))
	}

	value := e.Execute(expr.Value)
	obj := object.AsInstance()
	obj.Set(expr.Name.Lexeme, value)
	return value
}

func (e *Evaluator) getVariableValue(id string) values.Value {
	return e.env.GetOrDefault(id, values.NewNullValue())
}
