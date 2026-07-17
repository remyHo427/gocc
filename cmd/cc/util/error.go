package util

import (
	"fmt"
	"os"
	"runtime/debug"
)

func Exit_with_error(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	fmt.Print(string(debug.Stack()))
	os.Exit(1)
}

func Exit_with_printf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Print(string(debug.Stack()))
	os.Exit(1)
}

func Exit_with_println(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Print(string(debug.Stack()))
	os.Exit(1)
}
