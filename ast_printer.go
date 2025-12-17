package lox

import "fmt"

// astPrinter is a visitor that converts an expression tree to a string representation
// This makes it easy to verify the structure of parsed expressions
type AstPrinter struct {
	result string
}

func (ap *AstPrinter) Result() string {
	return ap.result
}

func (ap *AstPrinter) VisitBinary(b Binary) {
	left := printExpr(b.left)
	right := printExpr(b.right)
	ap.result = fmt.Sprintf("(%s %s %s)", b.operator.Lexeme, left, right)
}

func (ap *AstPrinter) VisitUnary(u Unary) {
	right := printExpr(u.right)
	ap.result = fmt.Sprintf("(%s %s)", u.operator.Lexeme, right)
}

func (ap *AstPrinter) VisitGroup(g Group) {
	inner := printExpr(g.expr)
	ap.result = fmt.Sprintf("(group %s)", inner)
}

func (ap *AstPrinter) VisitLiteral(l Literal) {
	if l.literal == nil {
		ap.result = "nil"
	} else {
		ap.result = fmt.Sprintf("%v", l.literal)
	}
}

func (ap *AstPrinter) VisitExprStmt(e ExprStmt) {}

func (ap *AstPrinter) VisitPrintStmt(p PrintStmt) {}

// printExpr is a helper function to convert any expression to its string representation
func printExpr(expr Expr) string {
	printer := &AstPrinter{}
	expr.Accept(printer)
	return printer.result
}
