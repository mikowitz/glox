package lox

import (
	"fmt"
	"os"
)

type Interpreter struct {
	runtime   *Runtime
	result    any
	lastError error
}

func NewInterpreter(lox *Runtime) *Interpreter {
	return &Interpreter{
		runtime: lox,
	}
}

func (i *Interpreter) Interpret(e Expr) (any, error) {
	i.result = nil
	i.lastError = nil
	i.runtime.HadRuntimeError = false
	i.evaluate(e)
	return i.result, i.lastError
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
		l, r, ok := i.checkNumbers(b.operator, l, r)
		if !ok {
			return
		}
		i.result = l > r
	case GreaterEqual:
		l, r, ok := i.checkNumbers(b.operator, l, r)
		if !ok {
			return
		}
		i.result = l >= r
	case Less:
		l, r, ok := i.checkNumbers(b.operator, l, r)
		if !ok {
			return
		}
		i.result = l < r
	case LessEqual:
		l, r, ok := i.checkNumbers(b.operator, l, r)
		if !ok {
			return
		}
		i.result = l <= r
	case Minus:
		l, r, ok := i.checkNumbers(b.operator, l, r)
		if !ok {
			return
		}
		i.result = l - r
	case Slash:
		l, r, ok := i.checkNumbers(b.operator, l, r)
		if !ok {
			return
		}
		i.result = l / r
	case Star:
		l, r, ok := i.checkNumbers(b.operator, l, r)
		if !ok {
			return
		}
		i.result = l * r
	case Plus:
		// Try numbers first
		if lNum, lOk := l.(float64); lOk {
			if rNum, rOk := r.(float64); rOk {
				i.result = lNum + rNum
				return
			}
		}
		// Try strings
		if lStr, lOk := l.(string); lOk {
			if rStr, rOk := r.(string); rOk {
				i.result = lStr + rStr
				return
			}
		}
		// Both failed, report error
		err := fmt.Errorf("operands to %s must both be numbers or strings", b.operator.Lexeme)
		i.reportError(err, b.operator)
	}
}

func (i *Interpreter) VisitUnary(u Unary) {
	current := i.result
	i.evaluate(u.right)
	r := i.result
	i.result = current

	switch u.operator.TokenType {
	case Minus:
		n, ok := i.checkNumber(u.operator, r)
		if !ok {
			return
		}
		i.result = -n
	case Bang:
		i.result = !isTruthy(r)
	}
}

func (i *Interpreter) reportError(err error, token Token) {
	i.runtime.HadRuntimeError = true
	err = fmt.Errorf("[line %d] %w: %s", token.Line, ErrLoxRuntime, err.Error())
	i.lastError = err
	fmt.Fprintf(os.Stderr, "%v\n", err)
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

func (i *Interpreter) checkNumber(operator Token, object any) (float64, bool) {
	n, ok := object.(float64)
	if ok {
		return n, true
	}
	err := fmt.Errorf("operand to %s must be a number", operator.Lexeme)
	i.reportError(err, operator)
	return 0, false
}

func (i *Interpreter) checkNumbers(operator Token, left, right any) (float64, float64, bool) {
	l, lok := left.(float64)
	r, rok := right.(float64)

	if lok && rok {
		return l, r, true
	}

	err := fmt.Errorf("operands to %s must both be numbers", operator.Lexeme)
	i.reportError(err, operator)
	return 0, 0, false
}

func (i *Interpreter) checkStrings(operator Token, left, right any) (string, string, bool) {
	l, lok := left.(string)
	r, rok := right.(string)

	if lok && rok {
		return l, r, true
	}

	err := fmt.Errorf("operands to %s must both be strings", operator.Lexeme)
	i.reportError(err, operator)
	return "", "", false
}
