package main

import (
	"io"
	"os"
	"testing"
)

// captureOutput captures stdout and stderr during test execution
func captureOutput(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()

	// Save original stdout/stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// Create pipes for stdout and stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	// Run the function
	fn()

	// Close writers and read output
	wOut.Close()
	wErr.Close()

	outBytes, _ := io.ReadAll(rOut)
	errBytes, _ := io.ReadAll(rErr)

	return string(outBytes), string(errBytes)
}

// writeTempFile creates a temporary .lox file with the given content
func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test-*.lox")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	return tmpFile.Name()
}

func TestRunFile_SimpleExpression(t *testing.T) {
	// Create a file with a simple expression
	filename := writeTempFile(t, "1 + 2;\n")

	// Run the file and capture output
	var exitCode int
	stdout, stderr := captureOutput(t, func() {
		exitCode = runFile(filename)
	})

	// Verify exit code is success
	if exitCode != ExitSuccess {
		t.Errorf("expected exit code %d, got %d", ExitSuccess, exitCode)
	}

	// Verify output
	if stdout != "3\n" {
		t.Errorf("expected stdout '3\\n', got %q", stdout)
	}

	// Verify no errors
	if stderr != "" {
		t.Errorf("expected no stderr, got %q", stderr)
	}
}
