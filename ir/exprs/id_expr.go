package exprs

// IdExpr 标识符表达式
type IdExpr struct {
	ID string
}

func NewIdExpr(id string) *IdExpr {
	return &IdExpr{ID: id}
}

func (e *IdExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitId(e)
}

func (e *IdExpr) String() string {
	return e.ID
}
