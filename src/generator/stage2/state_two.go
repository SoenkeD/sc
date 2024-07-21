package stage2

import (
	"fmt"
	"slices"
	"strings"

	"github.com/SoenkeD/sc/src/generator/parseuml"
	"github.com/SoenkeD/sc/src/types"
)

func AddUnique(sl []string, unique string) []string {

	if slices.Contains(sl, unique) {
		return sl
	}

	return append(sl, unique)
}

func PathJoin(basePath, addPath string) string {
	if basePath == "/" {
		return basePath + addPath
	}

	return strings.Join([]string{basePath, addPath}, "/")
}

func Stage2(stage1 parseuml.ParseUmlStage1) (uml ParseUmlStage2, err error) {
	uml.Name = stage1.Name
	uml.States = map[string]ParsedState{}
	uml.CompoundStates = map[string]string{}

	var stateStack1 []string
	for _, line := range stage1.LinesInOrder {

		statePath := "/" + strings.Join(stateStack1, "/")
		parts := strings.Split(line, ":")

		if len(parts) != 2 {
			err = fmt.Errorf("failed to parse stage 1 line\n%s", line)
			return
		}

		if parts[0] == "state_group" {
			stateStack1 = append(stateStack1, parts[1])
			uml.CompoundStates[PathJoin(statePath, parts[1])] = parts[1]
			continue
		}

		if parts[0] == "state_group_closing" {
			stateStack1 = stateStack1[0 : len(stateStack1)-1]
			continue
		}
	}

	var stateStack []string
	for _, line := range stage1.LinesInOrder {

		statePath := "/" + strings.Join(stateStack, "/")

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			err = fmt.Errorf("failed to parse stage 1 line\n%s", line)
			return
		}

		if parts[0] == "state_group" {
			stateStack = append(stateStack, parts[1])
			continue
		}

		if parts[0] == "state_group_closing" {
			uml.AddEndState(statePath)
			stateStack = stateStack[0 : len(stateStack)-1]
			continue
		}

		if parts[0] == "state_action" {

			stateAction := stage1.StateActions[parts[1]]

			addStatePath := PathJoin(statePath, stateAction.State)
			uml.AddAction2State(addStatePath, stateAction.Action, stateAction.ActionParams)
			uml.Actions = AddUnique(uml.Actions, stateAction.Action)
			continue
		}

		if parts[0] == "transition" {

			transition := stage1.Transitions[parts[1]]
			targetName := strings.ReplaceAll(transition.Target, "[*]", "End")

			targetPath := PathJoin(statePath, targetName)

			transitionStart := strings.ReplaceAll(transition.Start, "[*]", "Start")
			transitionStatePath := PathJoin(statePath, transitionStart)

			uml.AddTransition2State(
				transitionStatePath,
				ParseTransitionStage2{
					Action:       transition.Action,
					ActionParams: transition.ActionParams,
					Guard:        transition.Guard,
					GuardParams:  transition.GuardParams,
					Type:         transition.Type,
					Target:       targetPath,
					Negation:     transition.Negation,
				},
			)
			continue
		}
	}

	uml.AddEndState("/") // add end state

	return
}

func CheckForHappyPath(states map[string]ParsedState) error {
	for stateID, state := range states {

		if stateID == "/End" {
			continue
		}

		if len(state.Transitions) == 0 {
			return fmt.Errorf("state %s: every state must have any outgoing transitions", stateID)
		}

		if state.Transitions[len(state.Transitions)-1].Type != types.TransitionTypeHappy {
			return fmt.Errorf("state %s: the last transition of each state must be a happy path transition", stateID)
		}
	}

	return nil
}
