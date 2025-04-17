package prompt

import "strings"

// IsYesAnswer will return true if it is essentially a yes
// Yes,Y,y,yes,yes    ,    yes are some options
// Defaults to true
func IsYesAnswer(answer string) bool {
	if answer == "" {
		return true
	}

	answer = strings.TrimSpace(answer)
	answer = strings.ToLower(answer)
	return answer == "y" ||
		answer == "yes"
}

// isAbortAnswer will return true if it is essentially an abort
// Abort,A,a,abort,abort    ,    abort are some options
// Defaults to false
func isAbortAnswer(answer string) bool {
	if answer == "" {
		return false
	}

	answer = strings.TrimSpace(answer)
	answer = strings.ToLower(answer)
	return answer == "a" ||
		answer == "abort"
}
