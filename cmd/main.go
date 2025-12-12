package main

import (
	"bufio"
	"fmt"
	"os"

	lox "github.com/mikowitz/glox"
)

type Lox struct {
	hadError bool
}

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
	fmt.Fprintln(os.Stdout, "received:", source)
	scanner := lox.NewScanner(source)
	err := scanner.ScanTokens()
	if err != nil {
		return err
	}
	fmt.Println(scanner.Tokens)
	return nil
}
