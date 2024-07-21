package parseuml

import (
	"fmt"
	"strings"
)

const DO_SLASH_LINE = "do / "

func ParseStateActionFromCode(contentLinePart string) (action string, actionParams []string, err error) {
	firstPart := strings.TrimSpace(contentLinePart)

	if !strings.HasPrefix(firstPart, DO_SLASH_LINE) {
		err = fmt.Errorf("cannot parse: prefix not found")
		return
	}

	actionLineParts := strings.Split(firstPart, DO_SLASH_LINE)
	if len(actionLineParts) != 2 {
		err = fmt.Errorf("state action line")
		return
	}

	action, actionParams, err = ParseParams(actionLineParts[1])

	return
}
