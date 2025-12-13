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
		if err := runFile(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(65)
		}
	} else {
		if err := runPrompt(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
}

func runFile(filename string) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return run(string(bytes))
}

func runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "> ")
		if scanner.Scan() {
			_ = run(scanner.Text())
		} else {
			break
		}
	}
	return scanner.Err()
}

func run(source string) error {
	scanner := lox.NewScanner(source)
	err := scanner.ScanTokens()
	if err != nil {
		return err
	}

	parser := lox.NewParser(scanner.Tokens)
	expr, err := parser.Parse()
	if err != nil {
		return err
	}
	// fmt.Println(printExpr(expr))
	interpreter := lox.NewInterpreter()

	result, err := interpreter.Interpret(expr)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	fmt.Println(result)

	return nil
}
