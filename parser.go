package lox

import (
	"errors"
	"fmt"
	"slices"
)

type Parser struct {
	tokens  []Token
	current int
	errors  []error
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() ([]Stmt, error) {
	statements := []Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			p.errors = append(p.errors, err)
		}
		statements = append(statements, stmt)
	}
	return statements, errors.Join(p.errors...)
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(Print) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Semicolon, "expect ';' after value")
	if err != nil {
		return nil, err
	}

	return PrintStmt{expression: expr}, err
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Semicolon, "expect ';' after expression")
	if err != nil {
		return nil, err
	}
	return ExprStmt{expression: expr}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqual, EqualEqual) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = Binary{left: expr, right: right, operator: op}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = Binary{left: expr, right: right, operator: op}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(Minus, Plus) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = Binary{left: expr, right: right, operator: op}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(Slash, Star) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = Binary{left: expr, right: right, operator: op}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(Bang, Minus) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Unary{operator: op, right: right}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(False) {
		return Literal{literal: false}, nil
	}
	if p.match(True) {
		return Literal{literal: true}, nil
	}
	if p.match(Nil) {
		return Literal{literal: nil}, nil
	}

	if p.match(Number, String) {
		return Literal{literal: p.previous().Object}, nil
	}

	if p.match(LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RightParen, "expect ')' after expression")
		if err != nil {
			return nil, err
		}

		return Group{expr: expr}, nil
	}

	err := p.reportError("expect expression")
	return nil, err
}

func (p *Parser) consume(tokenType TokenType, msg string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	err := p.reportError(msg)
	return p.peek(), err
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	if slices.ContainsFunc(tokenTypes, func(tokenType TokenType) bool {
		return p.check(tokenType)
	}) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) reportError(msg string) error {
	token := p.peek()
	location := fmt.Sprintf("at '%s'", token.Lexeme)
	if p.isAtEnd() {
		location = "at end"
	}

	err := fmt.Errorf("[line %d] %w %s: %s", token.Line, ErrLoxSyntax, location, msg)
	return err
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == Semicolon {
			return
		}

		switch p.peek().TokenType {
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}

		p.advance()
	}
}
