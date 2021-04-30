package main

import (
	"os"

	"github.com/sh-miyoshi/hekate/cmd/debug-api/cmd"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
)

func main() {
	if err := os.MkdirAll("./tmp", 0775); err != nil {
		print.Error("Failed to create tmp directory: %v", err)
		os.Exit(1)
	}

	cmd.Execute()
}
