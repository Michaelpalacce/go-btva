package prompt

import "fmt"

// AskText will ask the user for text.
func AskText(prompts ...string) (string, error) {
	printPrompts(prompts...)
	fmt.Print(" ")

	var text string
	_, err := fmt.Scanln(&text)
	return text, err
}
