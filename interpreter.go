package lox

import "fmt"

type Interpreter struct {
	result any
	err    error
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(e Expr) (any, error) {
	i.result = nil
	i.err = nil
	i.evaluate(e)
	return i.result, i.err
}

func (i *Interpreter) evaluate(e Expr) {
	e.Accept(i)
}

func (i *Interpreter) VisitLiteral(l Literal) {
	i.result = l.literal
}

func (i *Interpreter) VisitGroup(g Group) {
	i.evaluate(g.expr)
}

func (i *Interpreter) VisitBinary(b Binary) {
	current := i.result
	i.evaluate(b.left)
	l := i.result
	i.evaluate(b.right)
	r := i.result
	i.result = current

	switch b.operator.TokenType {
	case BangEqual:
		i.result = l != r
	case EqualEqual:
		i.result = l == r
	case Greater:
		l, r, err := checkNumbers(b.operator, l, r)
		if err != nil {
			i.err = err
			return
		}
		i.result = l > r
	case GreaterEqual:
		l, r, err := checkNumbers(b.operator, l, r)
		if err != nil {
			i.err = err
			return
		}
		i.result = l >= r
	case Less:
		l, r, err := checkNumbers(b.operator, l, r)
		if err != nil {
			i.err = err
			return
		}
		i.result = l < r
	case LessEqual:
		l, r, err := checkNumbers(b.operator, l, r)
		if err != nil {
			i.err = err
			return
		}
		i.result = l <= r
	case Minus:
		l, r, err := checkNumbers(b.operator, l, r)
		if err != nil {
			i.err = err
			return
		}
		i.result = l - r
	case Slash:
		l, r, err := checkNumbers(b.operator, l, r)
		if err != nil {
			i.err = err
			return
		}
		i.result = l / r
	case Star:
		l, r, err := checkNumbers(b.operator, l, r)
		if err != nil {
			i.err = err
			return
		}
		i.result = l * r
	case Plus:
		lNum, rNum, err := checkNumbers(b.operator, l, r)
		if err == nil {
			i.result = lNum + rNum
			return
		}
		lStr, rStr, strErr := checkStrings(b.operator, l, r)
		if strErr == nil {
			i.result = lStr + rStr
			return
		}
		i.err = fmt.Errorf("operands to %s must both be numbers or strings", b.operator.Lexeme)
	}
}

func (i *Interpreter) VisitUnary(u Unary) {
	current := i.result
	i.evaluate(u.right)
	r := i.result
	i.result = current

	switch u.operator.TokenType {
	case Minus:
		n, err := checkNumber(u.operator, r)
		if err != nil {
			i.err = err
			return
		}
		i.result = -n
	case Bang:
		i.result = !isTruthy(r)
	}
}

func isTruthy(object any) bool {
	switch v := object.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}

func checkNumber(operator Token, object any) (float64, error) {
	n, ok := object.(float64)
	if ok {
		return n, nil
	}
	return n, fmt.Errorf("operand to %s must be a number, got %v", operator.Lexeme, object)
}

func checkNumbers(operator Token, left, right any) (float64, float64, error) {
	l, lok := left.(float64)
	r, rok := right.(float64)

	if lok && rok {
		return l, r, nil
	}

	return l, r, fmt.Errorf("operands to %s must both be numbers, got %T, %T",
		operator.Lexeme, left, right)
}

func checkStrings(operator Token, left, right any) (string, string, error) {
	l, lok := left.(string)
	r, rok := right.(string)

	if lok && rok {
		return l, r, nil
	}

	return l, r, fmt.Errorf("operands to %s must both be strings, got %T, %T",
		operator.Lexeme, left, right)
}
