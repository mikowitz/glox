package lox

type Visitor interface {
	VisitBinary(b Binary)
	VisitUnary(u Unary)
	VisitGroup(g Group)
	VisitLiteral(l Literal)
	VisitExprStmt(e ExprStmt)
	VisitPrintStmt(p PrintStmt)
}

type Expr interface {
	Accept(v Visitor)
}

type Binary struct {
	left, right Expr
	operator    Token
}

func (b Binary) Accept(v Visitor) {
	v.VisitBinary(b)
}

type Unary struct {
	right    Expr
	operator Token
}

func (u Unary) Accept(v Visitor) {
	v.VisitUnary(u)
}

type Group struct {
	expr Expr
}

func (g Group) Accept(v Visitor) {
	v.VisitGroup(g)
}

type Literal struct {
	literal any
}

func (l Literal) Accept(v Visitor) {
	v.VisitLiteral(l)
}
