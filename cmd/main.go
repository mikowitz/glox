package main

import (
	"bufio"
	"fmt"
	"os"

	lox "github.com/mikowitz/glox"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: glox [script]")
		os.Exit(64)
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
		return 66
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
			os.Exit(74)
		}
	}
}

func run(source string) int {
	scanner := lox.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 65 // Syntax error
	}

	parser := lox.NewParser(tokens)
	expr, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 65 // Syntax error
	}

	interpreter := lox.NewInterpreter()
	result, err := interpreter.Interpret(expr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 70 // Runtime error
	}

	fmt.Println(result)
	return 0
}
