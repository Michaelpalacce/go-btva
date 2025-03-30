package prompt

import (
	"fmt"
	"strings"
)

// AskYesNo will ask the user for a yes or no answer and return true of false
func AskYesNo(prompts ...string) (bool, error) {
	printPrompts(prompts...)
	fmt.Print(" [Y/n]: ")

	var text string
	_, err := fmt.Scanln(&text)

	return IsYesAnswer(text), ignoreUnexpectedNewline(err)
}

// ignoreUnexpectedNewline will return nil if the err is unexpected newline, as that is ... actually expected
func ignoreUnexpectedNewline(err error) error {
	if err == nil || strings.Contains(err.Error(), "unexpected newline") {
		return nil
	}
	return err
}
