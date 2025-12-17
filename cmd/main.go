package main

import (
	"bufio"
	"fmt"
	"os"

	lox "github.com/mikowitz/glox"
)

const (
	ExitSuccess      = 0
	ExitUsageError   = 64
	ExitSyntaxError  = 65
	ExitInputError   = 66
	ExitRuntimeError = 70
	ExitIOError      = 74
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: glox [script]")
		os.Exit(ExitUsageError)
	} else if len(os.Args) == 2 {
		exitCode := runFile(os.Args[1])
		os.Exit(exitCode)
	} else {
		runPrompt()
	}
}

func runFile(filename string) int {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return ExitInputError
	}
	return run(string(bytes))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "> ")
		if scanner.Scan() {
			run(scanner.Text())
		} else {
			os.Exit(ExitIOError)
		}
	}
}

func run(source string) int {
	scanner := lox.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitSyntaxError
	}

	fmt.Printf("%+#v\n", tokens)

	parser := lox.NewParser(tokens)
	expr, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitSyntaxError
	}

	fmt.Printf("%+#v\n", expr)

	interpreter := lox.NewInterpreter()
	result, err := interpreter.Interpret(expr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitRuntimeError
	}

	fmt.Println(result)
	return ExitSuccess
}
