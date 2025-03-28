package prompt

import "strings"

// IsYesAnswer will return true if it is essentially a yes
// Yes,Y,y,yes,yes    ,    yes are some options
func IsYesAnswer(answer string) bool {
	answer = strings.TrimSpace(answer)
	answer = strings.ToLower(answer)
	return answer == "y" ||
		answer == "yes"
}
