package lib

import (
	"fmt"
	"os"
)

func Fatal(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	_, _ = os.Stderr.WriteString("❗ " + msg + "\n")
	os.Exit(1)
}

func Info(message string, args ...interface{}) {
	line := fmt.Sprintf(message, args...)
	fmt.Printf(line + "\n")
}
