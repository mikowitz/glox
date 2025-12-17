package lox

type Stmt interface {
	Accept(v Visitor)
}

type ExprStmt struct {
	expression Expr
}

func (e ExprStmt) Accept(v Visitor) {
	v.VisitExprStmt(e)
}

type PrintStmt struct {
	expression Expr
}

func (p PrintStmt) Accept(v Visitor) {
	v.VisitPrintStmt(p)
}
