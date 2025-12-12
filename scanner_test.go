// ABOUTME: Tests for the scanner to ensure proper tokenization of source code
// ABOUTME: Uses table-driven tests to verify various token combinations
package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanTokens(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []Token
	}{
		{
			name:   "empty string",
			source: "",
			expected: []Token{
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "single left paren",
			source: "(",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "single right paren",
			source: ")",
			expected: []Token{
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "matching parens",
			source: "()",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "left brace",
			source: "{",
			expected: []Token{
				NewToken(LeftBrace, "{", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "right brace",
			source: "}",
			expected: []Token{
				NewToken(RightBrace, "}", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "matching braces",
			source: "{}",
			expected: []Token{
				NewToken(LeftBrace, "{", nil, 1),
				NewToken(RightBrace, "}", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "comma",
			source: ",",
			expected: []Token{
				NewToken(Comma, ",", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "dot",
			source: ".",
			expected: []Token{
				NewToken(Dot, ".", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "minus",
			source: "-",
			expected: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "plus",
			source: "+",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "semicolon",
			source: ";",
			expected: []Token{
				NewToken(Semicolon, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "star",
			source: "*",
			expected: []Token{
				NewToken(Star, "*", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "multiple operators",
			source: "+-*",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "all delimiters",
			source: "(){}",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(LeftBrace, "{", nil, 1),
				NewToken(RightBrace, "}", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "punctuation mix",
			source: ".,;",
			expected: []Token{
				NewToken(Dot, ".", nil, 1),
				NewToken(Comma, ",", nil, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "nested parens and braces",
			source: "({()})",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(LeftBrace, "{", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(RightBrace, "}", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "expression with delimiters",
			source: "{(+)*(-);};",
			expected: []Token{
				NewToken(LeftBrace, "{", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(RightBrace, "}", nil, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "comma separated values",
			source: "(,,,)",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Comma, ",", nil, 1),
				NewToken(Comma, ",", nil, 1),
				NewToken(Comma, ",", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "multiple statements",
			source: ";;;;;",
			expected: []Token{
				NewToken(Semicolon, ";", nil, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "bang",
			source: "!",
			expected: []Token{
				NewToken(Bang, "!", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "bang equal",
			source: "!=",
			expected: []Token{
				NewToken(BangEqual, "!=", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "equal",
			source: "=",
			expected: []Token{
				NewToken(Equal, "=", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "equal equal",
			source: "==",
			expected: []Token{
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "less",
			source: "<",
			expected: []Token{
				NewToken(Less, "<", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "less equal",
			source: "<=",
			expected: []Token{
				NewToken(LessEqual, "<=", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "greater",
			source: ">",
			expected: []Token{
				NewToken(Greater, ">", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "greater equal",
			source: ">=",
			expected: []Token{
				NewToken(GreaterEqual, ">=", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "comparison expression",
			source: "<=>",
			expected: []Token{
				NewToken(LessEqual, "<=", nil, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "equality check",
			source: "==!=",
			expected: []Token{
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(BangEqual, "!=", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "mixed single and double operators",
			source: "!===",
			expected: []Token{
				NewToken(BangEqual, "!=", nil, 1),
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "all comparison operators",
			source: "<><=>=",
			expected: []Token{
				NewToken(Less, "<", nil, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(LessEqual, "<=", nil, 1),
				NewToken(GreaterEqual, ">=", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "bang with parentheses",
			source: "!()",
			expected: []Token{
				NewToken(Bang, "!", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "assignment expression",
			source: "=+",
			expected: []Token{
				NewToken(Equal, "=", nil, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "comparison with arithmetic",
			source: "(+)>(-)",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "not equal comparison",
			source: "!=!",
			expected: []Token{
				NewToken(BangEqual, "!=", nil, 1),
				NewToken(Bang, "!", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "tokens with spaces",
			source: "+ - *",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "tokens with tabs",
			source: "+\t-\t*",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "multiple spaces between tokens",
			source: "+   -",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "leading and trailing spaces",
			source: "  +  ",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "slash",
			source: "/",
			expected: []Token{
				NewToken(Slash, "/", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "simple comment",
			source: "// this is a comment",
			expected: []Token{
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "line ending with comment",
			source: "+- // comment at end",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "tokens before comment",
			source: "(+) // this is a comment",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "comment with special chars",
			source: "// comment with @#$ special chars!",
			expected: []Token{
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "division vs comment",
			source: "/+",
			expected: []Token{
				NewToken(Slash, "/", nil, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "slash followed by other operator",
			source: "/*",
			expected: []Token{
				NewToken(Slash, "/", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "newline between tokens",
			source: "+\n-",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Minus, "-", nil, 2),
				NewToken(EOF, "", nil, 2),
			},
		},
		{
			name:   "multiple newlines",
			source: "+\n\n-",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Minus, "-", nil, 3),
				NewToken(EOF, "", nil, 3),
			},
		},
		{
			name:   "tokens on different lines",
			source: "(\n+\n)\n;",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Plus, "+", nil, 2),
				NewToken(RightParen, ")", nil, 3),
				NewToken(Semicolon, ";", nil, 4),
				NewToken(EOF, "", nil, 4),
			},
		},
		{
			name:   "mixed whitespace and newlines",
			source: "+  \n  -",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Minus, "-", nil, 2),
				NewToken(EOF, "", nil, 2),
			},
		},
		{
			name:   "comment after newline",
			source: "+\n//comment",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(EOF, "", nil, 2),
			},
		},
		{
			name:   "multiple lines with comments",
			source: "+ // first line\n- // second line",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Minus, "-", nil, 2),
				NewToken(EOF, "", nil, 2),
			},
		},
		{
			name:   "newline after comment",
			source: "//comment\n+",
			expected: []Token{
				NewToken(Plus, "+", nil, 2),
				NewToken(EOF, "", nil, 2),
			},
		},
		{
			name:   "empty lines between tokens",
			source: "+\n\n\n-",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Minus, "-", nil, 4),
				NewToken(EOF, "", nil, 4),
			},
		},
		{
			name:   "empty string",
			source: `""`,
			expected: []Token{
				NewToken(String, `""`, "", 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "simple string",
			source: `"hello"`,
			expected: []Token{
				NewToken(String, `"hello"`, "hello", 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "string with spaces",
			source: `"hello world"`,
			expected: []Token{
				NewToken(String, `"hello world"`, "hello world", 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "string with special characters",
			source: `"hello @#$%!"`,
			expected: []Token{
				NewToken(String, `"hello @#$%!"`, "hello @#$%!", 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "multiline string",
			source: "\"line1\nline2\"",
			expected: []Token{
				NewToken(String, "\"line1\nline2\"", "line1\nline2", 2),
				NewToken(EOF, "", nil, 2),
			},
		},
		{
			name:   "multiline string with multiple lines",
			source: "\"line1\nline2\nline3\"",
			expected: []Token{
				NewToken(String, "\"line1\nline2\nline3\"", "line1\nline2\nline3", 3),
				NewToken(EOF, "", nil, 3),
			},
		},
		{
			name:   "string between tokens",
			source: `+ "test" -`,
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(String, `"test"`, "test", 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "multiple strings",
			source: `"first" "second"`,
			expected: []Token{
				NewToken(String, `"first"`, "first", 1),
				NewToken(String, `"second"`, "second", 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "string on different line",
			source: "+\n\"test\"",
			expected: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(String, `"test"`, "test", 2),
				NewToken(EOF, "", nil, 2),
			},
		},
		{
			name:   "empty multiline string",
			source: "\"\n\"",
			expected: []Token{
				NewToken(String, "\"\n\"", "\n", 2),
				NewToken(EOF, "", nil, 2),
			},
		},
		{
			name:   "simple integer",
			source: "1234",
			expected: []Token{
				NewToken(Number, "1234", 1234.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "simple decimal",
			source: "12.34",
			expected: []Token{
				NewToken(Number, "12.34", 12.34, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "single digit",
			source: "0",
			expected: []Token{
				NewToken(Number, "0", 0.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "zero with decimal",
			source: "0.5",
			expected: []Token{
				NewToken(Number, "0.5", 0.5, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "integer with operators",
			source: "123+456",
			expected: []Token{
				NewToken(Number, "123", 123.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "456", 456.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "decimal with operators",
			source: "12.5+34.7",
			expected: []Token{
				NewToken(Number, "12.5", 12.5, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "34.7", 34.7, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "number with spaces",
			source: "123 + 456",
			expected: []Token{
				NewToken(Number, "123", 123.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "456", 456.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "number in parentheses",
			source: "(123)",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "123", 123.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "arithmetic expression with numbers",
			source: "12 + 34 * 56",
			expected: []Token{
				NewToken(Number, "12", 12.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "34", 34.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "56", 56.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "decimal arithmetic expression",
			source: "1.5 * 2.0",
			expected: []Token{
				NewToken(Number, "1.5", 1.5, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "2.0", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "numbers on multiple lines",
			source: "123\n456",
			expected: []Token{
				NewToken(Number, "123", 123.0, 1),
				NewToken(Number, "456", 456.0, 2),
				NewToken(EOF, "", nil, 2),
			},
		},
		{
			name:   "mixed integers and decimals",
			source: "100 + 3.14",
			expected: []Token{
				NewToken(Number, "100", 100.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "3.14", 3.14, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "number with semicolon",
			source: "42;",
			expected: []Token{
				NewToken(Number, "42", 42.0, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "complex expression with numbers",
			source: "(12.5 + 34) * 56.7",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "12.5", 12.5, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "34", 34.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "56.7", 56.7, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "numbers with comparison operators",
			source: "100 > 50",
			expected: []Token{
				NewToken(Number, "100", 100.0, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(Number, "50", 50.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "decimal comparison",
			source: "3.14 <= 3.15",
			expected: []Token{
				NewToken(Number, "3.14", 3.14, 1),
				NewToken(LessEqual, "<=", nil, 1),
				NewToken(Number, "3.15", 3.15, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "number subtraction",
			source: "100-50",
			expected: []Token{
				NewToken(Number, "100", 100.0, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "50", 50.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "number division",
			source: "100/25",
			expected: []Token{
				NewToken(Number, "100", 100.0, 1),
				NewToken(Slash, "/", nil, 1),
				NewToken(Number, "25", 25.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "large integer",
			source: "999999",
			expected: []Token{
				NewToken(Number, "999999", 999999.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "many decimal places",
			source: "3.14159265",
			expected: []Token{
				NewToken(Number, "3.14159265", 3.14159265, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "simple identifier",
			source: "foo",
			expected: []Token{
				NewToken(Identifier, "foo", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier with underscore prefix",
			source: "_test",
			expected: []Token{
				NewToken(Identifier, "_test", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier with numbers",
			source: "test123",
			expected: []Token{
				NewToken(Identifier, "test123", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier with underscores",
			source: "foo_bar",
			expected: []Token{
				NewToken(Identifier, "foo_bar", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier starting with uppercase",
			source: "Foo",
			expected: []Token{
				NewToken(Identifier, "Foo", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier all uppercase",
			source: "FOO",
			expected: []Token{
				NewToken(Identifier, "FOO", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier with multiple underscores",
			source: "__test__",
			expected: []Token{
				NewToken(Identifier, "__test__", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier starting with underscore and number",
			source: "_123",
			expected: []Token{
				NewToken(Identifier, "_123", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: and",
			source: "and",
			expected: []Token{
				NewToken(And, "and", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: class",
			source: "class",
			expected: []Token{
				NewToken(Class, "class", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: else",
			source: "else",
			expected: []Token{
				NewToken(Else, "else", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: false",
			source: "false",
			expected: []Token{
				NewToken(False, "false", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: fun",
			source: "fun",
			expected: []Token{
				NewToken(Fun, "fun", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: for",
			source: "for",
			expected: []Token{
				NewToken(For, "for", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: if",
			source: "if",
			expected: []Token{
				NewToken(If, "if", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: nil",
			source: "nil",
			expected: []Token{
				NewToken(Nil, "nil", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: or",
			source: "or",
			expected: []Token{
				NewToken(Or, "or", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: print",
			source: "print",
			expected: []Token{
				NewToken(Print, "print", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: return",
			source: "return",
			expected: []Token{
				NewToken(Return, "return", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: super",
			source: "super",
			expected: []Token{
				NewToken(Super, "super", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: this",
			source: "this",
			expected: []Token{
				NewToken(This, "this", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: true",
			source: "true",
			expected: []Token{
				NewToken(True, "true", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: var",
			source: "var",
			expected: []Token{
				NewToken(Var, "var", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word: while",
			source: "while",
			expected: []Token{
				NewToken(While, "while", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word with wrong case is identifier",
			source: "And",
			expected: []Token{
				NewToken(Identifier, "And", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word uppercase is identifier",
			source: "TRUE",
			expected: []Token{
				NewToken(Identifier, "TRUE", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "reserved word mixed case is identifier",
			source: "False",
			expected: []Token{
				NewToken(Identifier, "False", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifiers with operators",
			source: "foo + bar",
			expected: []Token{
				NewToken(Identifier, "foo", nil, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Identifier, "bar", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier with number",
			source: "foo + 123",
			expected: []Token{
				NewToken(Identifier, "foo", nil, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "123", 123.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier in parentheses",
			source: "(foo)",
			expected: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Identifier, "foo", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "complex expression with identifiers",
			source: "foo * bar + 123",
			expected: []Token{
				NewToken(Identifier, "foo", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Identifier, "bar", nil, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "123", 123.0, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifiers on multiple lines",
			source: "foo\nbar",
			expected: []Token{
				NewToken(Identifier, "foo", nil, 1),
				NewToken(Identifier, "bar", nil, 2),
				NewToken(EOF, "", nil, 2),
			},
		},
		{
			name:   "reserved words in expression",
			source: "if true or false",
			expected: []Token{
				NewToken(If, "if", nil, 1),
				NewToken(True, "true", nil, 1),
				NewToken(Or, "or", nil, 1),
				NewToken(False, "false", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "var declaration",
			source: "var foo = 123;",
			expected: []Token{
				NewToken(Var, "var", nil, 1),
				NewToken(Identifier, "foo", nil, 1),
				NewToken(Equal, "=", nil, 1),
				NewToken(Number, "123", 123.0, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "function declaration",
			source: "fun foo() {}",
			expected: []Token{
				NewToken(Fun, "fun", nil, 1),
				NewToken(Identifier, "foo", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(LeftBrace, "{", nil, 1),
				NewToken(RightBrace, "}", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "class declaration",
			source: "class Foo {}",
			expected: []Token{
				NewToken(Class, "class", nil, 1),
				NewToken(Identifier, "Foo", nil, 1),
				NewToken(LeftBrace, "{", nil, 1),
				NewToken(RightBrace, "}", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "if statement",
			source: "if (x > 10) print x;",
			expected: []Token{
				NewToken(If, "if", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Identifier, "x", nil, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(Number, "10", 10.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(Print, "print", nil, 1),
				NewToken(Identifier, "x", nil, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "while loop",
			source: "while (x < 10) x = x + 1;",
			expected: []Token{
				NewToken(While, "while", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Identifier, "x", nil, 1),
				NewToken(Less, "<", nil, 1),
				NewToken(Number, "10", 10.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(Identifier, "x", nil, 1),
				NewToken(Equal, "=", nil, 1),
				NewToken(Identifier, "x", nil, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "for loop",
			source: "for (var i = 0; i < 10; i = i + 1) {}",
			expected: []Token{
				NewToken(For, "for", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Var, "var", nil, 1),
				NewToken(Identifier, "i", nil, 1),
				NewToken(Equal, "=", nil, 1),
				NewToken(Number, "0", 0.0, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(Identifier, "i", nil, 1),
				NewToken(Less, "<", nil, 1),
				NewToken(Number, "10", 10.0, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(Identifier, "i", nil, 1),
				NewToken(Equal, "=", nil, 1),
				NewToken(Identifier, "i", nil, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(LeftBrace, "{", nil, 1),
				NewToken(RightBrace, "}", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "return statement",
			source: "return foo;",
			expected: []Token{
				NewToken(Return, "return", nil, 1),
				NewToken(Identifier, "foo", nil, 1),
				NewToken(Semicolon, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "this and super",
			source: "this.foo super.bar",
			expected: []Token{
				NewToken(This, "this", nil, 1),
				NewToken(Dot, ".", nil, 1),
				NewToken(Identifier, "foo", nil, 1),
				NewToken(Super, "super", nil, 1),
				NewToken(Dot, ".", nil, 1),
				NewToken(Identifier, "bar", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier that starts with reserved word",
			source: "ifTrue",
			expected: []Token{
				NewToken(Identifier, "ifTrue", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier that contains reserved word",
			source: "myClassVariable",
			expected: []Token{
				NewToken(Identifier, "myClassVariable", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
		{
			name:   "identifier that ends with reserved word",
			source: "isTrue",
			expected: []Token{
				NewToken(Identifier, "isTrue", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)
			scanner := NewScanner(tt.source)
			err := scanner.ScanTokens()

			asrt.False(scanner.hadError)
			asrt.NoError(err)
			asrt.Equal(tt.expected, scanner.Tokens)
		})
	}
}

func TestScanTokens_InvalidCharacters(t *testing.T) {
	tests := []struct {
		name             string
		source           string
		expectedErrorMsg string
	}{
		{
			name:             "at symbol",
			source:           "@",
			expectedErrorMsg: "[line 1] Error: unexpected character",
		},
		{
			name:             "hash symbol",
			source:           "#",
			expectedErrorMsg: "[line 1] Error: unexpected character",
		},
		{
			name:             "invalid char in expression",
			source:           "(+@)",
			expectedErrorMsg: "[line 1] Error: unexpected character",
		},
		{
			name:             "multiple invalid chars",
			source:           "@#",
			expectedErrorMsg: "[line 1] Error: unexpected character",
		},
		{
			name:             "error on line 2",
			source:           "+\n@",
			expectedErrorMsg: "[line 2] Error: unexpected character",
		},
		{
			name:             "error on line 3",
			source:           "+\n-\n@",
			expectedErrorMsg: "[line 3] Error: unexpected character",
		},
		{
			name:             "error after empty lines",
			source:           "+\n\n\n@",
			expectedErrorMsg: "[line 4] Error: unexpected character",
		},
		{
			name:             "error after comment",
			source:           "// comment\n@",
			expectedErrorMsg: "[line 2] Error: unexpected character",
		},
		{
			name:             "error on same line as valid tokens",
			source:           "+ - @",
			expectedErrorMsg: "[line 1] Error: unexpected character",
		},
		{
			name:             "error after multiple lines of valid code",
			source:           "+\n-\n*\n/\n@",
			expectedErrorMsg: "[line 5] Error: unexpected character",
		},
		{
			name:             "unterminated string on single line",
			source:           `"hello`,
			expectedErrorMsg: "[line 1] Error: unterminated string",
		},
		{
			name:             "unterminated string with one newline",
			source:           "\"hello\nworld",
			expectedErrorMsg: "[line 2] Error: unterminated string",
		},
		{
			name:             "unterminated string with two newlines",
			source:           "\"line1\nline2\nline3",
			expectedErrorMsg: "[line 3] Error: unterminated string",
		},
		{
			name:             "unterminated string with three newlines",
			source:           "\"a\nb\nc\nd",
			expectedErrorMsg: "[line 4] Error: unterminated string",
		},
		{
			name:             "unterminated string after valid tokens",
			source:           "+\n\"unterminated",
			expectedErrorMsg: "[line 2] Error: unterminated string",
		},
		{
			name:             "unterminated multiline string after valid code",
			source:           "+\n-\n\"start\nmore",
			expectedErrorMsg: "[line 4] Error: unterminated string",
		},
		{
			name:             "unterminated empty string",
			source:           `"`,
			expectedErrorMsg: "[line 1] Error: unterminated string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asrt := assert.New(t)

			scanner := NewScanner(tt.source)
			err := scanner.ScanTokens()

			asrt.True(scanner.hadError)

			asrt.Error(err)
			asrt.Equal(tt.expectedErrorMsg, err.Error())
		})
	}
}
