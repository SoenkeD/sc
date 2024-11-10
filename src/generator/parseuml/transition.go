package parseuml

import (
	"fmt"
	"slices"
	"strings"

	"github.com/SoenkeD/sc/src/types"
)

const TransitionNormal = "-->"
const TransitionHappy = "-[bold]->"
const TransitionError = "-[dotted]->"

func GetTransitionArgs(arrow string) ([]string, error) {

	if !strings.HasPrefix(arrow, "-") || strings.HasPrefix(strings.TrimSuffix(strings.TrimPrefix(arrow, "-"), "->"), "-") {
		return nil, fmt.Errorf("no valid arrow: not starting with -")
	}
	arrow = strings.TrimPrefix(arrow, "-")

	if !strings.HasSuffix(arrow, "->") {
		return nil, fmt.Errorf("no valid arrow: not ending with ->")
	}
	arrow = strings.TrimSuffix(arrow, "->")

	if len(arrow) == 0 {
		// has no arguments
		return nil, nil
	}

	// check if it has options in []
	if strings.HasSuffix(arrow, "]") {

		var args []string
		var optStr string
		if strings.HasPrefix(arrow, "[") {
			// has no directions but options
			opts := strings.TrimPrefix(arrow, "[")
			opts = strings.TrimSuffix(opts, "]")

			optStr = opts
		} else {
			// has direction and options
			parts := strings.Split(arrow, "[")

			if len(parts) != 2 {
				return nil, fmt.Errorf("")
			}

			args = append(args, parts[0])
			optStr = strings.TrimSuffix(parts[1], "]")
		}

		strArgs := strings.Split(optStr, ",")
		for _, strArg := range strArgs {
			args = append(args, strings.TrimSpace(strArg))
		}

		return args, nil

	}

	// only has direction
	return []string{arrow}, nil
}

func ParseTransitionType(arrow string) (taType types.TransitionType, options []string, err error) {

	args, err := GetTransitionArgs(arrow)
	if err != nil {
		return
	}

	if slices.Contains(args, "bold") {
		return types.TransitionTypeHappy, args, nil
	}

	if slices.Contains(args, "dotted") {
		return types.TransitionTypeError, args, nil
	}

	return types.TransitionTypeNormal, args, nil
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

	taType, options, err := ParseTransitionType(tokens[1])
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
		Options:      options,
	}

	return
}
