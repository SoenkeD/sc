package parseuml

import (
	"fmt"
	"strings"
)

func ParseParams(inputStr string) (name string, params []string, err error) {

	inputStr = strings.TrimSpace(inputStr)

	callParts := strings.Split(inputStr, "(")
	if len(callParts) == 1 {
		name = inputStr
		return
	}

	if len(callParts) == 2 {
		name = callParts[0]
		paramsStr := strings.ReplaceAll(callParts[1], ")", "")
		params = strings.Split(paramsStr, ",")

		for paramIdx, param := range params {
			params[paramIdx] = strings.TrimSpace(param)
		}
		return
	}

	err = fmt.Errorf("invalid input string given: expected 1 | 2 got %d, %s", len(callParts), inputStr)

	return
}
