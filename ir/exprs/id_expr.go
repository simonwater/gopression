package exprs

// IdExpr 标识符表达式
type IdExpr struct {
	Id string
}

func NewIdExpr(id string) *IdExpr {
	return &IdExpr{Id: id}
}

func (e *IdExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitId(e)
}

func (e *IdExpr) String() string {
	return e.Id
}
