package main

import (
	"fmt"
	"os"
)

func fatal(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	_, _ = os.Stderr.WriteString("❗ " + msg + "\n")
	os.Exit(1)
}

func info(message string, args ...interface{}) {
	line := fmt.Sprintf(message, args...)
	fmt.Printf(line + "\n")
}
