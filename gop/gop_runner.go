package gop

import (
	"github.com/simonwater/gopression/chk"
	"github.com/simonwater/gopression/env"
	"github.com/simonwater/gopression/exec"
	"github.com/simonwater/gopression/ir"
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/parser"
	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/visitors"
)

type ExecuteMode int

const (
	SyntaxTree ExecuteMode = iota
	ChunkVM
)

type GopRunner struct {
	needSort    bool
	executeMode ExecuteMode
	context     *ir.GopContext
}

func NewGopRunner() *GopRunner {
	return &GopRunner{
		needSort:    true,
		executeMode: SyntaxTree,
		context:     ir.NewGopContext(),
	}
}

func (r *GopRunner) IsNeedSort() bool {
	return r.needSort
}

func (r *GopRunner) SetNeedSort(needSort bool) {
	r.needSort = needSort
}

func (r *GopRunner) IsTrace() bool {
	return r.context.GetTracer().IsEnable()
}

func (r *GopRunner) SetTrace(isTrace bool) {
	r.context.GetTracer().SetEnable(isTrace)
}

func (r *GopRunner) GetExecuteMode() ExecuteMode {
	return r.executeMode
}

func (r *GopRunner) SetExecuteMode(executeMode ExecuteMode) {
	r.executeMode = executeMode
}

func (r *GopRunner) Execute(expression string, ev ...env.Environment) interface{} {
	var e env.Environment
	if len(ev) == 0 {
		e = env.NewDefaultEnvironment()
	} else {
		e = ev[0]
	}
	result := r.executeBatch([]string{expression}, e)
	if len(result) == 0 {
		return nil
	}
	return result[0]
}

func (r *GopRunner) ExecuteBatch(expressions []string, ev ...env.Environment) []interface{} {
	var e env.Environment
	if len(ev) == 0 {
		e = env.NewDefaultEnvironment()
	} else {
		e = ev[0]
	}
	return r.executeBatch(expressions, e)
}

func (r *GopRunner) executeBatch(expressions []string, env env.Environment) []interface{} {
	tracer := r.context.GetTracer()
	tracer.StartTimerWithMsg("开始。公式总数：%d", len(expressions))

	exprs := r.Parse(expressions)
	exprInfos := r.Analyze(exprs)

	var result []interface{}
	if r.executeMode == ChunkVM {
		chunk := r.CompileIR(exprInfos)
		result = r.RunChunk(chunk, env)
	} else {
		result = r.RunIR(exprInfos, env)
	}

	tracer.EndTimer("结束。")
	return result
}

func (r *GopRunner) RunIR(exprInfos []*ir.ExprInfo, ev env.Environment) []interface{} {
	tracer := r.context.GetTracer()
	tracer.StartTimer()

	variables := make(map[string]bool)
	for _, info := range exprInfos {
		for name := range info.GetPrecursors() {
			variables[name] = true
		}
		for name := range info.GetSuccessors() {
			variables[name] = true
		}
	}

	fields := make([]*util.Field, 0, len(variables))
	for v := range variables {
		fields = append(fields, util.NewField(v))
	}

	flag := ev.BeforeExecute(fields)
	tracer.EndTimer("完成执行环境初始化。")
	if !flag {
		return nil
	}

	tracer.StartTimerWithMsg("执行")
	n := len(exprInfos)
	result := make([]interface{}, n)

	for _, info := range exprInfos {
		expr := info.GetExpr()
		evaluator := visitors.NewEvaluator(ev)
		v := evaluator.Execute(expr)
		result[info.GetIndex()] = v.GetValue()
	}

	tracer.EndTimer("执行完成。")
	return result
}

func (r *GopRunner) RunChunk(chunk *chk.Chunk, ev env.Environment) []interface{} {
	tracer := r.context.GetTracer()
	tracer.StartTimer()

	chunkReader := chk.NewChunkReader(chunk, tracer)
	variables := chunkReader.GetVariables()

	fields := make([]*util.Field, 0, len(variables))
	for _, v := range variables {
		fields = append(fields, util.NewField(v))
	}

	flag := ev.BeforeExecute(fields)
	tracer.EndTimer("完成执行环境初始化。")
	if !flag {
		return nil
	}

	tracer.StartTimerWithMsg("执行")
	vm := exec.NewVM(tracer)
	exResults, _ := vm.ExecuteWithReader(chunkReader, ev)

	result := make([]interface{}, len(exResults))
	for _, res := range exResults {
		result[res.GetIndex()] = res.GetResult().GetValue()
	}

	tracer.EndTimer("执行完成。")
	return result
}

func (r *GopRunner) Parse(expressions []string) []exprs.Expr {
	tracer := r.context.GetTracer()
	tracer.StartTimerWithMsg("解析")

	result := make([]exprs.Expr, 0, len(expressions))
	for _, src := range expressions {
		parser := parser.NewParser(src)
		result = append(result, parser.Parse())
	}

	tracer.EndTimer("完成表达式解析。")
	return result
}

func (r *GopRunner) Analyze(exprs []exprs.Expr) []*ir.ExprInfo {
	tracer := r.context.GetTracer()
	tracer.StartTimerWithMsg("分析")

	exprInfos := make([]*ir.ExprInfo, len(exprs))
	for i, expr := range exprs {
		exprInfos[i] = ir.NewExprInfo(expr, i)
	}

	r.context.PrepareExecute(exprInfos)
	sortedInfos := r.sortExprs(exprInfos)

	tracer.EndTimer("完成表达式分析。")
	return sortedInfos
}

func (r *GopRunner) CompileSource(expressions []string) *chk.Chunk {
	tracer := r.context.GetTracer()
	tracer.StartTimerWithMsg("编译源码")

	exprs := r.Parse(expressions)
	exprInfos := r.Analyze(exprs)
	chunk := r.CompileIR(exprInfos)

	tracer.EndTimer("完成表达式编译。")
	return chunk
}

func (r *GopRunner) CompileIR(exprInfos []*ir.ExprInfo) *chk.Chunk {
	tracer := r.context.GetTracer()
	tracer.StartTimerWithMsg("编译中间表示")

	compiler := visitors.NewOpCodeCompiler(tracer, len(exprInfos))
	compiler.BeginCompile()

	for _, info := range exprInfos {
		compiler.Compile(info)
	}

	result := compiler.EndCompile()
	tracer.EndTimer("完成表达式编译。")
	return result
}

func (r *GopRunner) sortExprs(exprInfos []*ir.ExprInfo) []*ir.ExprInfo {
	if r.needSort && len(exprInfos) >= 1 && r.context.GetExecContext().HasAssign() {
		sorter := ir.NewExprSorter(r.context)
		sortedExprInfos, _ := sorter.Sort()
		return sortedExprInfos
	}
	return exprInfos
}
