package prompt

import "fmt"

// AskText will ask the user for text.
func AskText() (string, error) {
	var text string
	_, err := fmt.Scanln(&text)
	return text, err
}
