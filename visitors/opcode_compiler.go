package visitors

import (
	"github.com/simonwater/gopression/chk"
	"github.com/simonwater/gopression/functions/funmgr"
	"github.com/simonwater/gopression/ir"
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/values"
)

const ADDRESS_SIZE = 4 // 地址大小（4字节）

// OpCodeCompiler 字节码编译器
type OpCodeCompiler struct {
	chunkWriter *chk.ChunkWriter
	varSet      map[string]bool
	tracer      *util.Tracer
}

// NewOpCodeCompiler 创建新的编译器
func NewOpCodeCompiler(tracer *util.Tracer, chunkCapacity ...int) *OpCodeCompiler {
	compiler := &OpCodeCompiler{
		varSet: make(map[string]bool),
		tracer: tracer,
	}

	if len(chunkCapacity) > 0 && chunkCapacity[0] > 0 {
		compiler.chunkWriter = chk.NewChunkWriter(chunkCapacity[0], tracer)
	} else {
		compiler.chunkWriter = chk.NewChunkWriter(0, tracer)
	}

	return compiler
}

// BeginCompile 开始编译
func (c *OpCodeCompiler) BeginCompile() {
	c.chunkWriter.Clear()
	c.varSet = make(map[string]bool)
}

// Compile 编译表达式信息
func (c *OpCodeCompiler) Compile(exprInfo *ir.ExprInfo) {
	expr := exprInfo.GetExpr()
	order := exprInfo.GetIndex()
	c.CompileExpr(expr, order)

	// 添加前置和后继变量到集合
	for name := range exprInfo.GetPrecursors() {
		c.varSet[name] = true
	}
	for name := range exprInfo.GetSuccessors() {
		c.varSet[name] = true
	}
}

// CompileExpr 编译单个表达式
func (c *OpCodeCompiler) CompileExpr(expr exprs.Expr, order int) {
	c.emitOp(chk.OP_BEGIN, order)
	c.execute(expr)
	c.emitOp(chk.OP_END)
}

// EndCompile 结束编译并返回块
func (c *OpCodeCompiler) EndCompile() *chk.Chunk {
	c.emitOp(chk.OP_EXIT)

	// 转换变量集合为切片
	vars := make([]string, 0, len(c.varSet))
	for name := range c.varSet {
		vars = append(vars, name)
	}

	c.chunkWriter.SetVariables(vars)
	return c.chunkWriter.Flush()
}

func (c *OpCodeCompiler) execute(expr exprs.Expr) any {
	if expr == nil {
		return nil
	}
	exprs.VisitExpr(expr, c)
	return nil
}

// 实现表达式访问者接口
func (c *OpCodeCompiler) VisitBinary(expr *exprs.BinaryExpr) any {
	c.execute(expr.Left)
	c.execute(expr.Right)

	switch expr.Operator.Type {
	case values.PLUS:
		c.emitOp(chk.OP_ADD)
	case values.MINUS:
		c.emitOp(chk.OP_SUBTRACT)
	case values.STAR:
		c.emitOp(chk.OP_MULTIPLY)
	case values.SLASH:
		c.emitOp(chk.OP_DIVIDE)
	case values.PERCENT:
		c.emitOp(chk.OP_MODE)
	case values.STARSTAR:
		c.emitOp(chk.OP_POWER)
	case values.GREATER:
		c.emitOp(chk.OP_GREATER)
	case values.GREATER_EQUAL:
		c.emitOp(chk.OP_GREATER_EQUAL)
	case values.LESS:
		c.emitOp(chk.OP_LESS)
	case values.LESS_EQUAL:
		c.emitOp(chk.OP_LESS_EQUAL)
	case values.BANG_EQUAL:
		c.emitOp(chk.OP_BANG_EQUAL)
	case values.EQUAL_EQUAL:
		c.emitOp(chk.OP_EQUAL_EQUAL)
	}
	return nil
}

func (c *OpCodeCompiler) VisitLogic(expr *exprs.LogicExpr) any {
	c.execute(expr.Left)

	if expr.Operator.Type == values.AND { // AND
		jumper := c.emitJump(chk.OP_JUMP_IF_FALSE)
		c.emitOp(chk.OP_POP)
		c.execute(expr.Right)
		c.patchJump(jumper)
	} else { // OR
		jumper1 := c.emitJump(chk.OP_JUMP_IF_FALSE)
		jumper2 := c.emitJump(chk.OP_JUMP)
		c.patchJump(jumper1)
		c.emitOp(chk.OP_POP)
		c.execute(expr.Right)
		c.patchJump(jumper2)
	}
	return nil
}

func (c *OpCodeCompiler) VisitLiteral(expr *exprs.LiteralExpr) any {
	c.emitConstant(expr.Value)
	return nil
}

func (c *OpCodeCompiler) VisitUnary(expr *exprs.UnaryExpr) any {
	c.execute(expr.Right)

	switch expr.Operator.Type {
	case values.BANG:
		c.emitOp(chk.OP_NOT)
	case values.MINUS:
		c.emitOp(chk.OP_NEGATE)
	}
	return nil
}

func (c *OpCodeCompiler) VisitId(expr *exprs.IdExpr) any {
	value := values.NewStringValue(expr.Id)
	constIndex := c.makeConstant(&value)
	c.emitOp(chk.OP_GET_GLOBAL, constIndex)
	return nil
}

func (c *OpCodeCompiler) VisitAssign(expr *exprs.AssignExpr) any {
	c.execute(expr.Right)

	// 假设左侧是IdExpr
	if idExpr, ok := expr.Left.(*exprs.IdExpr); ok {
		value := values.NewStringValue(idExpr.Id)
		constIndex := c.makeConstant(&value)
		c.emitOp(chk.OP_SET_GLOBAL, constIndex)
	} else {
		// 处理其他类型的左值
		c.execute(expr.Left)
	}
	return nil
}

func (c *OpCodeCompiler) VisitCall(expr *exprs.CallExpr) any {
	// 假设被调用者是IdExpr
	if idExpr, ok := expr.Callee.(*exprs.IdExpr); ok {
		name := idExpr.Id
		fn := funmgr.GetFunctionManager().GetFunction(name)
		if fn == nil {
			panic("未定义的函数: " + name)
		}

		if len(expr.Args) != fn.Arity() {
			panic("参数数量不匹配")
		}

		// 编译所有参数
		for _, arg := range expr.Args {
			c.execute(arg)
		}

		value := values.NewStringValue(name)
		constIndex := c.makeConstant(&value)
		c.emitOp(chk.OP_CALL, constIndex)
	} else {
		panic("不支持的调用表达式")
	}
	return nil
}

func (c *OpCodeCompiler) VisitIf(expr *exprs.IfExpr) any {
	c.execute(expr.Condition)
	elseJumper := c.emitJump(chk.OP_JUMP_IF_FALSE)
	c.emitOp(chk.OP_POP)
	c.execute(expr.ThenBranch)
	endJumper := c.emitJump(chk.OP_JUMP)
	c.patchJump(elseJumper)
	c.emitOp(chk.OP_POP)

	if expr.ElseBranch != nil {
		c.execute(expr.ElseBranch)
	} else {
		c.emitOp(chk.OP_NULL)
	}

	c.patchJump(endJumper)
	return nil
}

func (c *OpCodeCompiler) VisitGet(expr *exprs.GetExpr) any {
	c.execute(expr.Object)
	value := values.NewStringValue(expr.Name.Lexeme)
	constIndex := c.makeConstant(&value)
	c.emitOp(chk.OP_GET_PROPERTY, constIndex)
	return nil
}

func (c *OpCodeCompiler) VisitSet(expr *exprs.SetExpr) any {
	c.execute(expr.Value)
	c.execute(expr.Object)
	value := values.NewStringValue(expr.Name.Lexeme)
	constIndex := c.makeConstant(&value)
	c.emitOp(chk.OP_SET_PROPERTY, constIndex)
	return nil
}

// emitJump 发出跳转指令并返回跳转地址位置
func (c *OpCodeCompiler) emitJump(jumpCode chk.OpCode) int {
	c.emitOp(jumpCode)
	c.emitInt(0xffffffff) // 占位符
	return c.chunkWriter.Position() - ADDRESS_SIZE
}

// patchJump 修补跳转指令
func (c *OpCodeCompiler) patchJump(index int) {
	offset := c.chunkWriter.Position() - index - ADDRESS_SIZE
	c.chunkWriter.UpdateInt(index, int32(offset))
}

// emitOp 发出操作码
func (c *OpCodeCompiler) emitOp(opCode chk.OpCode, args ...int) {
	c.chunkWriter.WriteCode(opCode)
	if len(args) > 0 {
		c.emitInt(args[0])
	}
}

// emitConstant 发出常量指令
func (c *OpCodeCompiler) emitConstant(value *values.Value) {
	index := c.makeConstant(value)
	c.emitOp(chk.OP_CONSTANT, index)
}

// makeConstant 创建常量并返回索引
func (c *OpCodeCompiler) makeConstant(value *values.Value) int {
	index, err := c.chunkWriter.AddConstant(value)
	if err != nil {
		panic("创建常量失败: " + err.Error())
	}
	return index
}

// emitInt 发出整数
func (c *OpCodeCompiler) emitInt(value int) {
	c.chunkWriter.WriteInt(int32(value))
}
