package util

import (
	"fmt"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// ReadPasswordFromConsole ...
func ReadPasswordFromConsole() (string, error) {
	passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()
	password := string(passwordBytes)
	return password, nil
}
