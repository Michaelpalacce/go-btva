package prompt

import "fmt"

func printPrompts(prompts ...string) {
	promptLength := len(prompts)

	for i := 0; i < promptLength; i++ {
		if i >= promptLength-1 {
			fmt.Print(prompts[i])
		} else {
			fmt.Println(prompts[i])
		}
	}
}
