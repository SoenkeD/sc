package parseuml

import (
	"fmt"
	"strings"

	"github.com/SoenkeD/sc/src/types"
)

type ParsedTransition struct {
	Type         types.TransitionType
	Start        string
	Target       string
	Guard        string
	GuardParams  []string
	Action       string
	ActionParams []string
	Negation     bool
	Options      []string
}

type ParseStateAction struct {
	State        string
	Action       string
	ActionParams []string
}

type ParseUmlStage1 struct {
	Name         string
	LinesInOrder []string
	StateGroups  map[string]string
	Transitions  map[string]ParsedTransition
	StateActions map[string]ParseStateAction
}

const LINE_STATE_GROUP_CLOSING = "state_group_closing:"
const LINE_STATE_GROUP_OPENING = "state_group:"
const LINE_ACTION = "state_action:"
const LINE_TRANSITION = "transition:"

func (uml *ParseUmlStage1) Init() {
	uml.StateGroups = map[string]string{}
	uml.Transitions = map[string]ParsedTransition{}
	uml.StateActions = map[string]ParseStateAction{}
}

func (uml *ParseUmlStage1) AddStateGroupClosing() {
	uml.LinesInOrder = append(uml.LinesInOrder, LINE_STATE_GROUP_CLOSING)
}

func (uml *ParseUmlStage1) AddStateGroup(stateGroup string) {
	uml.LinesInOrder = append(uml.LinesInOrder, fmt.Sprintf("%s%s", LINE_STATE_GROUP_OPENING, stateGroup))
	uml.StateGroups[stateGroup] = stateGroup
}

func (uml *ParseUmlStage1) AddTransition(parsedTa ParsedTransition) {
	merged := string(parsedTa.Type) + parsedTa.Start + parsedTa.Target +
		parsedTa.Guard + strings.Join(parsedTa.GuardParams, ",") +
		parsedTa.Action + strings.Join(parsedTa.ActionParams, ",")
	uml.LinesInOrder = append(uml.LinesInOrder, fmt.Sprintf("%s%s", LINE_TRANSITION, merged))
	uml.Transitions[merged] = parsedTa
}

func (uml *ParseUmlStage1) AddStateAction(state, action string, actionParams []string) {
	merged := state + action + strings.Join(actionParams, "")
	uml.LinesInOrder = append(uml.LinesInOrder, fmt.Sprintf("%s%s", LINE_ACTION, merged))
	uml.StateActions[merged] = ParseStateAction{
		State:        state,
		Action:       action,
		ActionParams: actionParams,
	}
}
