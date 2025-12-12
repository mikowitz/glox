package lox

type TokenType int

const (
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	Identifier
	String
	Number

	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

type Token struct {
	TokenType
	Lexeme string
	Object any
	Line   int
}

func NewToken(tokenType TokenType, lexeme string, object any, line int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Object:    object,
		Line:      line,
	}
}
