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

// AskYesNo will ask the user for a yes or no answer and return true of false
func AskYesNoAbort(prompts ...string) (bool, error) {
	printPrompts(prompts...)
	fmt.Print(" [Y/n/a]: ")

	var text string
	_, err := fmt.Scanln(&text)

	return IsYesAnswer(text), ignoreAbortAnswer(text, err)
}

// ignoreAbortAnswer will return an error if the user typed "a/A/Abort/abort"
func ignoreAbortAnswer(text string, err error) error {
	if ignoreUnexpectedNewline(err) != nil {
		return err
	}

	if isAbort := isAbortAnswer(text); isAbort {
		return fmt.Errorf("user aborted")
	}

	return nil
}

// ignoreUnexpectedNewline will return nil if the err is unexpected newline, as that is ... actually expected
func ignoreUnexpectedNewline(err error) error {
	if err == nil || strings.Contains(err.Error(), "unexpected newline") {
		return nil
	}
	return err
}
