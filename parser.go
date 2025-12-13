package lox

import (
	"fmt"
	"os"
	"slices"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
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

		return Binary{left: expr, right: right, operator: op}, nil
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

		return Binary{left: expr, right: right, operator: op}, nil
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

		return Binary{left: expr, right: right, operator: op}, nil
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

		return Binary{left: expr, right: right, operator: op}, nil
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
		_, err = p.consume(RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return Group{expr: expr}, nil
	}

	token := p.peek()
	var err error
	if token.TokenType == EOF {
		err = fmt.Errorf("[line %d] Error at end: expect expression", token.Line)
	} else {
		err = fmt.Errorf("[line %d] Error at '%s': expect expression", token.Line, token.Lexeme)
	}

	fmt.Fprintf(os.Stderr, "%v\n", err)
	return nil, err
}

func (p *Parser) consume(tokenType TokenType, msg string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	token := p.peek()
	var err error
	if token.TokenType == EOF {
		err = fmt.Errorf("[line %d] Error at end: %s", token.Line, msg)
	} else {
		err = fmt.Errorf("[line %d] Error at '%s': %s", token.Line, token.Lexeme, msg)
	}

	fmt.Fprintf(os.Stderr, "%v\n", err)
	return token, err
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	hasToken := slices.ContainsFunc(tokenTypes, func(tt TokenType) bool {
		return p.check(tt)
	})
	if hasToken {
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
