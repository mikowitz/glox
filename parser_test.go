// ABOUTME: Tests for the parser to ensure proper AST construction from tokens
// ABOUTME: Uses visitor pattern to verify parsed expression values
package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// testVisitor is a helper visitor that extracts literal values for testing
type testVisitor struct {
	result any
}

func (tv *testVisitor) VisitBinary(b Binary) {
	tv.result = "binary expression"
}

func (tv *testVisitor) VisitUnary(u Unary) {
	tv.result = "unary expression"
}

func (tv *testVisitor) VisitGroup(g Group) {
	tv.result = "group expression"
}

func (tv *testVisitor) VisitLiteral(l Literal) {
	tv.result = l.literal
}

func TestParser_Expressions(t *testing.T) {
	tests := []struct {
		name         string
		tokens       []Token
		expectedAST  string
		wantErr      bool
		errorMessage string
	}{
		// Literal expressions - testing primary()
		{
			name: "literal: true",
			tokens: []Token{
				NewToken(True, "true", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "true",
		},
		{
			name: "literal: false",
			tokens: []Token{
				NewToken(False, "false", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "false",
		},
		{
			name: "literal: nil",
			tokens: []Token{
				NewToken(Nil, "nil", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "nil",
		},
		{
			name: "literal: integer",
			tokens: []Token{
				NewToken(Number, "123", 123.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "123",
		},
		{
			name: "literal: decimal",
			tokens: []Token{
				NewToken(Number, "45.67", 45.67, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "45.67",
		},
		{
			name: "literal: string",
			tokens: []Token{
				NewToken(String, `"hello"`, "hello", 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "hello",
		},

		// Grouped expressions - testing primary()
		{
			name: "grouped: simple number",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "123", 123.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(group 123)",
		},
		{
			name: "grouped: boolean",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(True, "true", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(group true)",
		},

		// Unary expressions - testing unary()
		{
			name: "unary: negation",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "123", 123.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(- 123)",
		},
		{
			name: "unary: logical not",
			tokens: []Token{
				NewToken(Bang, "!", nil, 1),
				NewToken(True, "true", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(! true)",
		},
		{
			name: "unary: double negation",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(- (- 5))",
		},
		{
			name: "unary: double not",
			tokens: []Token{
				NewToken(Bang, "!", nil, 1),
				NewToken(Bang, "!", nil, 1),
				NewToken(False, "false", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(! (! false))",
		},

		// Factor expressions - testing factor() (multiplication and division)
		{
			name: "factor: multiplication",
			tokens: []Token{
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(* 2 3)",
		},
		{
			name: "factor: division",
			tokens: []Token{
				NewToken(Number, "10", 10.0, 1),
				NewToken(Slash, "/", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(/ 10 2)",
		},
		{
			name: "factor: multiplication with unary",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(* (- 2) 3)",
		},

		// Term expressions - testing term() (addition and subtraction)
		{
			name: "term: addition",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(+ 1 2)",
		},
		{
			name: "term: subtraction",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(- 5 3)",
		},
		{
			name: "term: addition with multiplication (precedence)",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(+ 1 (* 2 3))",
		},
		{
			name: "term: multiplication with addition (precedence)",
			tokens: []Token{
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(+ (* 2 3) 1)",
		},

		// Comparison expressions - testing comparison()
		{
			name: "comparison: greater than",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(> 5 3)",
		},
		{
			name: "comparison: greater than or equal",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(GreaterEqual, ">=", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(>= 5 5)",
		},
		{
			name: "comparison: less than",
			tokens: []Token{
				NewToken(Number, "3", 3.0, 1),
				NewToken(Less, "<", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(< 3 5)",
		},
		{
			name: "comparison: less than or equal",
			tokens: []Token{
				NewToken(Number, "3", 3.0, 1),
				NewToken(LessEqual, "<=", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(<= 3 3)",
		},
		{
			name: "comparison: with arithmetic",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(> (+ 1 2) 3)",
		},

		// Equality expressions - testing equality()
		{
			name: "equality: equal",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(== 1 1)",
		},
		{
			name: "equality: not equal",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(BangEqual, "!=", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(!= 1 2)",
		},
		{
			name: "equality: with comparison",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(True, "true", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(== (> 5 3) true)",
		},

		// Complex nested expressions
		{
			name: "complex: grouped with operator precedence",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(* (group (+ 1 2)) 3)",
		},
		{
			name: "complex: nested groups",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(group (group 5))",
		},
		{
			name: "complex: unary with grouped expression",
			tokens: []Token{
				NewToken(Bang, "!", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(True, "true", nil, 1),
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(False, "false", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(! (group (== true false)))",
		},
		{
			name: "complex: full arithmetic expression",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(+ 1 (* 2 3))",
		},
		{
			name: "complex: comparison with arithmetic on both sides",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(LessEqual, "<=", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(<= (+ 1 2) (- 5 1))",
		},
		{
			name: "complex: multiple unary operators",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(Bang, "!", nil, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(- (! (- 5)))",
		},

		// String expressions
		{
			name: "string: in comparison",
			tokens: []Token{
				NewToken(String, `"hello"`, "hello", 1),
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(String, `"world"`, "world", 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedAST: "(== hello world)",
		},

		// Error cases
		{
			name: "error: missing closing paren",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "123", 123.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			wantErr:      true,
			errorMessage: "expect ')' after expression",
		},
		{
			name: "error: unexpected EOF",
			tokens: []Token{
				NewToken(EOF, "", nil, 1),
			},
			wantErr:      true,
			errorMessage: "expect expression",
		},
		{
			name: "error: unexpected token",
			tokens: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			wantErr:      true,
			errorMessage: "expect expression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			parser := NewParser(tt.tokens)
			expr, err := parser.expression()

			if tt.wantErr {
				asrt.Error(err)
				if tt.errorMessage != "" {
					asrt.Contains(err.Error(), tt.errorMessage)
				}
				return
			}

			asrt.NoError(err)
			asrt.NotNil(expr)

			// Use AST printer to verify structure
			actualAST := printExpr(expr)
			asrt.Equal(tt.expectedAST, actualAST)
		})
	}
}
