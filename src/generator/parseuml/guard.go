package parseuml

import (
	"fmt"
	"strings"
)

func ParseGuard(guardPart string) (guard string, guardParams []string, negation bool, err error) {

	guardPart = strings.TrimSpace(guardPart)

	if len(guardPart) == 0 {
		// no guard found, returning defaults
		return
	}

	inputLen := len(guardPart)

	guardPart = strings.ReplaceAll(guardPart, "[", "")
	guardPart = strings.ReplaceAll(guardPart, "]", "")

	if inputLen-2 != len(guardPart) {
		err = fmt.Errorf("missing [] around the guard")
		return
	}

	guardPart = strings.TrimSpace(guardPart)

	if strings.HasPrefix(guardPart, "!") {
		negation = true
		guardPart = strings.ReplaceAll(guardPart, "!", "")
	}

	guard, guardParams, err = ParseParams(guardPart)
	if err != nil {
		return
	}

	return
}
