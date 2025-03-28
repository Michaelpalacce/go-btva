package prompt

import (
	"syscall"

	"golang.org/x/term"
)

// AskPass will ask the user for a password. Typed text will be hidden
func AskPass(prompts ...string) (string, error) {
	printPrompts(prompts...)

	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	return string(bytepw), err
}
