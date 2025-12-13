package lox

import (
	"fmt"
	"os"
	"strconv"
)

type Scanner struct {
	source         string
	Tokens         []Token
	start, current int
	line           int
	hadError       bool
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		line:   1,
	}
}

func (s *Scanner) ScanTokens() error {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return err
		}
	}

	s.Tokens = append(s.Tokens, NewToken(EOF, "", nil, s.line))
	return nil
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case ',':
		s.addToken(Comma)
	case '.':
		s.addToken(Dot)
	case '-':
		s.addToken(Minus)
	case '+':
		s.addToken(Plus)
	case ';':
		s.addToken(Semicolon)
	case '*':
		s.addToken(Star)
	case '!':
		if s.match('=') {
			s.addToken(BangEqual)
		} else {
			s.addToken(Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual)
		} else {
			s.addToken(Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual)
		} else {
			s.addToken(Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual)
		} else {
			s.addToken(Greater)
		}
	case '/':
		if s.match('/') {
			s.handleComment()
		} else {
			s.addToken(Slash)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	case '"':
		return s.handleString()
	default:
		if isDigit(c) {
			return s.handleNumber()
		} else if isAlpha(c) {
			return s.handleIdentifier()
		}
		s.hadError = true
		return s.reportError("unexpected character")
	}
	return nil
}

func (s *Scanner) addToken(tokenType TokenType) {
	lexeme := s.source[s.start:s.current]
	s.Tokens = append(s.Tokens, NewToken(tokenType, lexeme, nil, s.line))
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	s.Tokens = append(s.Tokens, NewToken(tokenType, lexeme, literal, s.line))
}

func (s *Scanner) advance() rune {
	char := s.source[s.current]
	s.current++
	return rune(char)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return rune(0)
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) match(r rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != byte(r) {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) handleComment() {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
}

func (s *Scanner) handleString() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.hadError = true
		return s.reportError("unterminated string")
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(String, value)
	return nil
}

func (s *Scanner) handleNumber() error {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
	}

	for isDigit(s.peek()) {
		s.advance()
	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		return s.reportError("unparseable number")
	}
	s.addTokenWithLiteral(Number, value)
	return nil
}

func (s *Scanner) handleIdentifier() error {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	lexeme := s.source[s.start:s.current]
	tokenType := Identifier
	if tt, ok := keywords[lexeme]; ok {
		tokenType = tt
	}
	s.addToken(tokenType)
	return nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) reportError(msg string) error {
	err := fmt.Errorf("[line %d] %w: %s", s.line, ErrLoxSyntax, msg)
	fmt.Fprintf(os.Stderr, "%v\n", err)
	return err
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}
