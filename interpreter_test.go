// ABOUTME: Tests for the interpreter to ensure correct evaluation of expressions
// ABOUTME: Covers all operations, type checking, and error handling
package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpreter_Literals(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
	}{
		{
			name:     "literal: true",
			expr:     Literal{literal: true},
			expected: true,
		},
		{
			name:     "literal: false",
			expr:     Literal{literal: false},
			expected: false,
		},
		{
			name:     "literal: nil",
			expr:     Literal{literal: nil},
			expected: nil,
		},
		{
			name:     "literal: integer",
			expr:     Literal{literal: 123.0},
			expected: 123.0,
		},
		{
			name:     "literal: decimal",
			expr:     Literal{literal: 45.67},
			expected: 45.67,
		},
		{
			name:     "literal: string",
			expr:     Literal{literal: "hello"},
			expected: "hello",
		},
		{
			name:     "literal: zero",
			expr:     Literal{literal: 0.0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			interp := NewInterpreter()
			result, err := interp.Interpret(tt.expr)

			if tt.wantErr {
				asrt.Error(err)
				return
			}

			asrt.NoError(err)
			asrt.Equal(tt.expected, result)
		})
	}
}

func TestInterpreter_UnaryOperations(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
		errMsg   string
	}{
		{
			name: "unary: negate positive number",
			expr: Unary{
				operator: NewToken(Minus, "-", nil, 1),
				right:    Literal{literal: 123.0},
			},
			expected: -123.0,
		},
		{
			name: "unary: negate negative number",
			expr: Unary{
				operator: NewToken(Minus, "-", nil, 1),
				right:    Literal{literal: -45.0},
			},
			expected: 45.0,
		},
		{
			name: "unary: negate zero",
			expr: Unary{
				operator: NewToken(Minus, "-", nil, 1),
				right:    Literal{literal: 0.0},
			},
			expected: -0.0,
		},
		{
			name: "unary: logical not true",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: true},
			},
			expected: false,
		},
		{
			name: "unary: logical not false",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: false},
			},
			expected: true,
		},
		{
			name: "unary: logical not nil",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: nil},
			},
			expected: true,
		},
		{
			name: "unary: logical not number (truthy)",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: 123.0},
			},
			expected: false,
		},
		{
			name: "unary: logical not zero (truthy)",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: 0.0},
			},
			expected: false,
		},
		{
			name: "unary: double negation",
			expr: Unary{
				operator: NewToken(Minus, "-", nil, 1),
				right: Unary{
					operator: NewToken(Minus, "-", nil, 1),
					right:    Literal{literal: 5.0},
				},
			},
			expected: 5.0,
		},
		{
			name: "unary: double not",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right: Unary{
					operator: NewToken(Bang, "!", nil, 1),
					right:    Literal{literal: true},
				},
			},
			expected: true,
		},
		{
			name: "unary: negate non-number (error)",
			expr: Unary{
				operator: NewToken(Minus, "-", nil, 1),
				right:    Literal{literal: "hello"},
			},
			wantErr: true,
			errMsg:  "operand to - must be a number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			interp := NewInterpreter()
			result, err := interp.Interpret(tt.expr)

			if tt.wantErr {
				asrt.Error(err)
				if tt.errMsg != "" {
					asrt.Contains(err.Error(), tt.errMsg)
				}
				return
			}

			asrt.NoError(err)
			asrt.Equal(tt.expected, result)
		})
	}
}

func TestInterpreter_ArithmeticOperations(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
		errMsg   string
	}{
		{
			name: "arithmetic: addition",
			expr: Binary{
				left:     Literal{literal: 1.0},
				operator: NewToken(Plus, "+", nil, 1),
				right:    Literal{literal: 2.0},
			},
			expected: 3.0,
		},
		{
			name: "arithmetic: subtraction",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(Minus, "-", nil, 1),
				right:    Literal{literal: 3.0},
			},
			expected: 2.0,
		},
		{
			name: "arithmetic: multiplication",
			expr: Binary{
				left:     Literal{literal: 2.0},
				operator: NewToken(Star, "*", nil, 1),
				right:    Literal{literal: 3.0},
			},
			expected: 6.0,
		},
		{
			name: "arithmetic: division",
			expr: Binary{
				left:     Literal{literal: 10.0},
				operator: NewToken(Slash, "/", nil, 1),
				right:    Literal{literal: 2.0},
			},
			expected: 5.0,
		},
		{
			name: "arithmetic: negative result",
			expr: Binary{
				left:     Literal{literal: 3.0},
				operator: NewToken(Minus, "-", nil, 1),
				right:    Literal{literal: 5.0},
			},
			expected: -2.0,
		},
		{
			name: "arithmetic: division with decimal result",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(Slash, "/", nil, 1),
				right:    Literal{literal: 2.0},
			},
			expected: 2.5,
		},
		{
			name: "arithmetic: multiplication with zero",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(Star, "*", nil, 1),
				right:    Literal{literal: 0.0},
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			interp := NewInterpreter()
			result, err := interp.Interpret(tt.expr)

			if tt.wantErr {
				asrt.Error(err)
				if tt.errMsg != "" {
					asrt.Contains(err.Error(), tt.errMsg)
				}
				return
			}

			asrt.NoError(err)
			asrt.Equal(tt.expected, result)
		})
	}
}

func TestInterpreter_StringOperations(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
		errMsg   string
	}{
		{
			name: "string: concatenation",
			expr: Binary{
				left:     Literal{literal: "hello"},
				operator: NewToken(Plus, "+", nil, 1),
				right:    Literal{literal: "world"},
			},
			expected: "helloworld",
		},
		{
			name: "string: concatenation with space",
			expr: Binary{
				left:     Literal{literal: "hello "},
				operator: NewToken(Plus, "+", nil, 1),
				right:    Literal{literal: "world"},
			},
			expected: "hello world",
		},
		{
			name: "string: concatenation empty strings",
			expr: Binary{
				left:     Literal{literal: ""},
				operator: NewToken(Plus, "+", nil, 1),
				right:    Literal{literal: ""},
			},
			expected: "",
		},
		{
			name: "string: mixed types with plus (error)",
			expr: Binary{
				left:     Literal{literal: "hello"},
				operator: NewToken(Plus, "+", nil, 1),
				right:    Literal{literal: 123.0},
			},
			wantErr: true,
			errMsg:  "operands to + must both be numbers or strings",
		},
		{
			name: "string: number plus string (error)",
			expr: Binary{
				left:     Literal{literal: 123.0},
				operator: NewToken(Plus, "+", nil, 1),
				right:    Literal{literal: "hello"},
			},
			wantErr: true,
			errMsg:  "operands to + must both be numbers or strings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			interp := NewInterpreter()
			result, err := interp.Interpret(tt.expr)

			if tt.wantErr {
				asrt.Error(err)
				if tt.errMsg != "" {
					asrt.Contains(err.Error(), tt.errMsg)
				}
				return
			}

			asrt.NoError(err)
			asrt.Equal(tt.expected, result)
		})
	}
}

func TestInterpreter_ComparisonOperations(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
		errMsg   string
	}{
		{
			name: "comparison: greater than (true)",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(Greater, ">", nil, 1),
				right:    Literal{literal: 3.0},
			},
			expected: true,
		},
		{
			name: "comparison: greater than (false)",
			expr: Binary{
				left:     Literal{literal: 3.0},
				operator: NewToken(Greater, ">", nil, 1),
				right:    Literal{literal: 5.0},
			},
			expected: false,
		},
		{
			name: "comparison: greater than equal (greater)",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(GreaterEqual, ">=", nil, 1),
				right:    Literal{literal: 3.0},
			},
			expected: true,
		},
		{
			name: "comparison: greater than equal (equal)",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(GreaterEqual, ">=", nil, 1),
				right:    Literal{literal: 5.0},
			},
			expected: true,
		},
		{
			name: "comparison: greater than equal (less)",
			expr: Binary{
				left:     Literal{literal: 3.0},
				operator: NewToken(GreaterEqual, ">=", nil, 1),
				right:    Literal{literal: 5.0},
			},
			expected: false,
		},
		{
			name: "comparison: less than (true)",
			expr: Binary{
				left:     Literal{literal: 3.0},
				operator: NewToken(Less, "<", nil, 1),
				right:    Literal{literal: 5.0},
			},
			expected: true,
		},
		{
			name: "comparison: less than (false)",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(Less, "<", nil, 1),
				right:    Literal{literal: 3.0},
			},
			expected: false,
		},
		{
			name: "comparison: less than equal (less)",
			expr: Binary{
				left:     Literal{literal: 3.0},
				operator: NewToken(LessEqual, "<=", nil, 1),
				right:    Literal{literal: 5.0},
			},
			expected: true,
		},
		{
			name: "comparison: less than equal (equal)",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(LessEqual, "<=", nil, 1),
				right:    Literal{literal: 5.0},
			},
			expected: true,
		},
		{
			name: "comparison: less than equal (greater)",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(LessEqual, "<=", nil, 1),
				right:    Literal{literal: 3.0},
			},
			expected: false,
		},
		{
			name: "comparison: with negative numbers",
			expr: Binary{
				left:     Literal{literal: -5.0},
				operator: NewToken(Less, "<", nil, 1),
				right:    Literal{literal: -3.0},
			},
			expected: true,
		},
		{
			name: "comparison: non-number operand (error)",
			expr: Binary{
				left:     Literal{literal: "hello"},
				operator: NewToken(Greater, ">", nil, 1),
				right:    Literal{literal: 5.0},
			},
			wantErr: true,
			errMsg:  "operands to > must both be numbers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			interp := NewInterpreter()
			result, err := interp.Interpret(tt.expr)

			if tt.wantErr {
				asrt.Error(err)
				if tt.errMsg != "" {
					asrt.Contains(err.Error(), tt.errMsg)
				}
				return
			}

			asrt.NoError(err)
			asrt.Equal(tt.expected, result)
		})
	}
}

func TestInterpreter_EqualityOperations(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
	}{
		{
			name: "equality: numbers equal",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(EqualEqual, "==", nil, 1),
				right:    Literal{literal: 5.0},
			},
			expected: true,
		},
		{
			name: "equality: numbers not equal",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(EqualEqual, "==", nil, 1),
				right:    Literal{literal: 3.0},
			},
			expected: false,
		},
		{
			name: "equality: strings equal",
			expr: Binary{
				left:     Literal{literal: "hello"},
				operator: NewToken(EqualEqual, "==", nil, 1),
				right:    Literal{literal: "hello"},
			},
			expected: true,
		},
		{
			name: "equality: strings not equal",
			expr: Binary{
				left:     Literal{literal: "hello"},
				operator: NewToken(EqualEqual, "==", nil, 1),
				right:    Literal{literal: "world"},
			},
			expected: false,
		},
		{
			name: "equality: booleans equal",
			expr: Binary{
				left:     Literal{literal: true},
				operator: NewToken(EqualEqual, "==", nil, 1),
				right:    Literal{literal: true},
			},
			expected: true,
		},
		{
			name: "equality: booleans not equal",
			expr: Binary{
				left:     Literal{literal: true},
				operator: NewToken(EqualEqual, "==", nil, 1),
				right:    Literal{literal: false},
			},
			expected: false,
		},
		{
			name: "equality: nil equals nil",
			expr: Binary{
				left:     Literal{literal: nil},
				operator: NewToken(EqualEqual, "==", nil, 1),
				right:    Literal{literal: nil},
			},
			expected: true,
		},
		{
			name: "equality: different types",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(EqualEqual, "==", nil, 1),
				right:    Literal{literal: "5"},
			},
			expected: false,
		},
		{
			name: "inequality: numbers not equal",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(BangEqual, "!=", nil, 1),
				right:    Literal{literal: 3.0},
			},
			expected: true,
		},
		{
			name: "inequality: numbers equal",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(BangEqual, "!=", nil, 1),
				right:    Literal{literal: 5.0},
			},
			expected: false,
		},
		{
			name: "inequality: different types",
			expr: Binary{
				left:     Literal{literal: 5.0},
				operator: NewToken(BangEqual, "!=", nil, 1),
				right:    Literal{literal: "5"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			interp := NewInterpreter()
			result, err := interp.Interpret(tt.expr)

			asrt.NoError(err)
			asrt.Equal(tt.expected, result)
		})
	}
}

func TestInterpreter_GroupedExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
	}{
		{
			name: "grouped: simple number",
			expr: Group{
				expr: Literal{literal: 123.0},
			},
			expected: 123.0,
		},
		{
			name: "grouped: boolean",
			expr: Group{
				expr: Literal{literal: true},
			},
			expected: true,
		},
		{
			name: "grouped: arithmetic expression",
			expr: Group{
				expr: Binary{
					left:     Literal{literal: 1.0},
					operator: NewToken(Plus, "+", nil, 1),
					right:    Literal{literal: 2.0},
				},
			},
			expected: 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			interp := NewInterpreter()
			result, err := interp.Interpret(tt.expr)

			asrt.NoError(err)
			asrt.Equal(tt.expected, result)
		})
	}
}

func TestInterpreter_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
	}{
		{
			name: "complex: (1 + 2) * 3",
			expr: Binary{
				left: Group{
					expr: Binary{
						left:     Literal{literal: 1.0},
						operator: NewToken(Plus, "+", nil, 1),
						right:    Literal{literal: 2.0},
					},
				},
				operator: NewToken(Star, "*", nil, 1),
				right:    Literal{literal: 3.0},
			},
			expected: 9.0,
		},
		{
			name: "complex: !(true == false)",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right: Group{
					expr: Binary{
						left:     Literal{literal: true},
						operator: NewToken(EqualEqual, "==", nil, 1),
						right:    Literal{literal: false},
					},
				},
			},
			expected: true,
		},
		{
			name: "complex: -5 * 3",
			expr: Binary{
				left: Unary{
					operator: NewToken(Minus, "-", nil, 1),
					right:    Literal{literal: 5.0},
				},
				operator: NewToken(Star, "*", nil, 1),
				right:    Literal{literal: 3.0},
			},
			expected: -15.0,
		},
		{
			name: "complex: 1 + 2 > 2",
			expr: Binary{
				left: Binary{
					left:     Literal{literal: 1.0},
					operator: NewToken(Plus, "+", nil, 1),
					right:    Literal{literal: 2.0},
				},
				operator: NewToken(Greater, ">", nil, 1),
				right:    Literal{literal: 2.0},
			},
			expected: true,
		},
		{
			name: "complex: nested arithmetic (10 / 2) + (3 * 4)",
			expr: Binary{
				left: Group{
					expr: Binary{
						left:     Literal{literal: 10.0},
						operator: NewToken(Slash, "/", nil, 1),
						right:    Literal{literal: 2.0},
					},
				},
				operator: NewToken(Plus, "+", nil, 1),
				right: Group{
					expr: Binary{
						left:     Literal{literal: 3.0},
						operator: NewToken(Star, "*", nil, 1),
						right:    Literal{literal: 4.0},
					},
				},
			},
			expected: 17.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			interp := NewInterpreter()
			result, err := interp.Interpret(tt.expr)

			asrt.NoError(err)
			asrt.Equal(tt.expected, result)
		})
	}
}

func TestInterpreter_TypeErrors(t *testing.T) {
	tests := []struct {
		name   string
		expr   Expr
		errMsg string
	}{
		{
			name: "error: subtract non-numbers",
			expr: Binary{
				left:     Literal{literal: "hello"},
				operator: NewToken(Minus, "-", nil, 1),
				right:    Literal{literal: "world"},
			},
			errMsg: "operands to - must both be numbers",
		},
		{
			name: "error: multiply non-numbers",
			expr: Binary{
				left:     Literal{literal: true},
				operator: NewToken(Star, "*", nil, 1),
				right:    Literal{literal: false},
			},
			errMsg: "operands to * must both be numbers",
		},
		{
			name: "error: divide non-numbers",
			expr: Binary{
				left:     Literal{literal: "hello"},
				operator: NewToken(Slash, "/", nil, 1),
				right:    Literal{literal: 5.0},
			},
			errMsg: "operands to / must both be numbers",
		},
		{
			name: "error: compare string and number",
			expr: Binary{
				left:     Literal{literal: "hello"},
				operator: NewToken(Less, "<", nil, 1),
				right:    Literal{literal: 5.0},
			},
			errMsg: "operands to < must both be numbers",
		},
		{
			name: "error: negate boolean",
			expr: Unary{
				operator: NewToken(Minus, "-", nil, 1),
				right:    Literal{literal: true},
			},
			errMsg: "operand to - must be a number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			interp := NewInterpreter()
			_, err := interp.Interpret(tt.expr)

			asrt.Error(err)
			asrt.Contains(err.Error(), tt.errMsg)
		})
	}
}

func TestInterpreter_Truthiness(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected bool
	}{
		{
			name: "truthiness: nil is falsy",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: nil},
			},
			expected: true,
		},
		{
			name: "truthiness: false is falsy",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: false},
			},
			expected: true,
		},
		{
			name: "truthiness: true is truthy",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: true},
			},
			expected: false,
		},
		{
			name: "truthiness: zero is truthy",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: 0.0},
			},
			expected: false,
		},
		{
			name: "truthiness: non-zero number is truthy",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: 123.0},
			},
			expected: false,
		},
		{
			name: "truthiness: empty string is truthy",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: ""},
			},
			expected: false,
		},
		{
			name: "truthiness: non-empty string is truthy",
			expr: Unary{
				operator: NewToken(Bang, "!", nil, 1),
				right:    Literal{literal: "hello"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			interp := NewInterpreter()
			result, err := interp.Interpret(tt.expr)

			asrt.NoError(err)
			asrt.Equal(tt.expected, result)
		})
	}
}
