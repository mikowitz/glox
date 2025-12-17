package lox

import "errors"

var (
	ErrLoxSyntax  = errors.New("syntax error")
	ErrLoxRuntime = errors.New("runtime error")
)
