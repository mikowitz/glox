package main

import (
	"bufio"
	"fmt"
	"os"

	lox "github.com/mikowitz/glox"
)

func main() {
	runtime := &lox.Runtime{}
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: glox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		// if err := runFile(os.Args[1]); err != nil {
		runFile(runtime, os.Args[1])
		if runtime.HadSyntaxError {
			os.Exit(65)
		}
	} else {
		runPrompt(runtime)
	}
}

func runFile(runtime *lox.Runtime, filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		os.Exit(66)
	}
	run(runtime, string(bytes))
}

func runPrompt(runtime *lox.Runtime) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "> ")
		if scanner.Scan() {
			run(runtime, scanner.Text())
			runtime.HadSyntaxError = false
			runtime.HadRuntimeError = false
		} else {
			os.Exit(74)
		}
	}
}

func run(runtime *lox.Runtime, source string) {
	scanner := lox.NewScanner(runtime, source)
	scanner.ScanTokens()
	if runtime.HadSyntaxError {
		return
	}

	parser := lox.NewParser(runtime, scanner.Tokens)
	expr, _ := parser.Parse()

	if runtime.HadSyntaxError {
		return
	}

	interpreter := lox.NewInterpreter(runtime)
	result, _ := interpreter.Interpret(expr)
	if runtime.HadRuntimeError {
		return
	}

	fmt.Println(result)
}
