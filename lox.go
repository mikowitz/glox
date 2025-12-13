package lox

import (
	"fmt"
	"os"
)

func reportError(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s", line, where, message)
}
