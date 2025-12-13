package lox

import "errors"

var (
	ErrLoxSyntax  = errors.New("syntax error")
	ErrLoxRuntime = errors.New("runtime error")
)

type Runtime struct {
	hadError        bool
	hadRuntimeError bool
}
