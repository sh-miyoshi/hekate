package print

import (
	"fmt"
	"os"

	"github.com/sh-miyoshi/hekate/pkg/jwtctl/config"
)

// Debug method output debug message is run as debug mode
func Debug(format string, a ...interface{}) {
	if config.Get().EnableDebug {
		fmt.Printf(format, a...)
	}
}

// Print ...
func Print(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

// Error ...
func Error(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "[ERROR] %s", msg)
}
