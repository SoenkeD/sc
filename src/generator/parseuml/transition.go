package parseuml

import (
	"fmt"
	"strings"

	"github.com/SoenkeD/sc/src/types"
)

const TransitionNormal = "-->"
const TransitionHappy = "-[bold]->"
const TransitionError = "-[dotted]->"

func ParseTransitionType(arrow string) (taType types.TransitionType, err error) {

	switch arrow {
	case TransitionNormal:
		taType = types.TransitionTypeNormal
	case "-[bold]->":
		taType = types.TransitionTypeHappy // TODO add happy handling
	case "-[dotted]->":
		taType = types.TransitionTypeError
	default:
		err = fmt.Errorf("unknown transition type: %s", arrow)
		return
	}

	return
}

func ParseTransition(tokens []string, linePart2 string) (ta ParsedTransition, err error) {

	if len(tokens) != 3 {
		err = fmt.Errorf("expected %d tokens got %d", 3, len(tokens))
		return
	}

	var negation bool
	var guard string
	var guardParams []string
	var action string
	var actionParams []string

	if len(linePart2) > 0 {

		transactionConditionParts := strings.Split(linePart2, "/")

		guardPart := strings.TrimSpace(transactionConditionParts[0])
		if len(guardPart) > 0 {

			guard, guardParams, negation, err = ParseGuard(guardPart)
			if err != nil {
				return
			}
		}

		if len(transactionConditionParts) == 2 {
			actionPart := strings.TrimSpace(transactionConditionParts[1])
			action, actionParams, err = ParseParams(actionPart)
			if err != nil {
				return
			}
		}
	}

	taType, err := ParseTransitionType(tokens[1])
	if err != nil {
		return
	}

	ta = ParsedTransition{
		Type:         taType,
		Start:        tokens[0],
		Target:       tokens[2],
		Guard:        guard,
		GuardParams:  guardParams,
		Action:       action,
		ActionParams: actionParams,
		Negation:     negation,
	}

	return
}
