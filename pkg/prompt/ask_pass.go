package prompt

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

// AskPass will ask the user for a password. Typed text will be hidden
// Print an empty line cause otherwise next log will be on same line
func AskPass(prompts ...string) (string, error) {
	printPrompts(prompts...)
	fmt.Print(" ")

	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	return string(bytepw), err
}
