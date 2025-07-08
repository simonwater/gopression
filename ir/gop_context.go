package ir

import "github.com/simonwater/gopression/util"

type GopContext struct {
	tracer      *util.Tracer
	execContext *ExecuteContext
}

func NewGopContext() *GopContext {
	ctx := GopContext{
		tracer: util.NewTracer(),
	}
	execCtx := NewExecuteContext(&ctx)
	ctx.execContext = execCtx
	return &ctx
}

func (ctx *GopContext) GetTracer() *util.Tracer {
	return ctx.tracer
}

func (ctx *GopContext) GetExecContext() *ExecuteContext {
	return ctx.execContext
}

func (ctx *GopContext) SetExecContext(execCtx *ExecuteContext) {
	ctx.execContext = execCtx
}

func (ctx *GopContext) PrepareExecute(exprInfos []*ExprInfo) {
	ctx.execContext.PreExecute(exprInfos)
}
