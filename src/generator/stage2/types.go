package stage2

import (
	"github.com/SoenkeD/sc/src/types"
)

type ParseTransitionStage2 struct {
	Action       string
	ActionParams []string
	Guard        string
	GuardParams  []string
	Type         types.TransitionType
	Target       string
	Negation     bool
}

type ParsedStateActionStage2 struct {
	Action       string
	ActionParams []string
}

type ParsedState struct {
	Actions     []ParsedStateActionStage2
	Transitions []ParseTransitionStage2
}

type ParsedTransition struct {
}

type ParseUmlStage2 struct {
	Name           string
	Actions        []string
	Guards         []string
	States         map[string]ParsedState
	CompoundStates map[string]string
}

func (uml *ParseUmlStage2) AddAction2State(state, action string, actionParams []string) {
	appState := uml.States[state]
	appState.Actions = append(appState.Actions, ParsedStateActionStage2{
		Action:       action,
		ActionParams: actionParams,
	})
	uml.States[state] = appState
}

func (uml *ParseUmlStage2) AddTransition2State(state string, parsedTa ParseTransitionStage2) {

	_, ok := uml.CompoundStates[state]
	if ok {
		state += "/End"
	}

	_, ok = uml.CompoundStates[parsedTa.Target]
	if ok {
		parsedTa.Target += "/Start"
	}

	appState := uml.States[state]
	appState.Transitions = append(appState.Transitions, parsedTa)
	uml.States[state] = appState

	targetState := uml.States[parsedTa.Target]
	uml.States[parsedTa.Target] = targetState

	if parsedTa.Guard != "" {
		uml.Guards = AddUnique(uml.Guards, parsedTa.Guard)
	}

	if parsedTa.Action != "" {
		uml.Actions = AddUnique(uml.Actions, parsedTa.Action)
	}
}

func (uml *ParseUmlStage2) AddStartingState(state string) {
	state = PathJoin(state, "Start")
	appState := uml.States[state]
	uml.States[state] = appState
}

func (uml *ParseUmlStage2) AddEndState(state string) {
	state = PathJoin(state, "End")
	appState := uml.States[state]
	uml.States[state] = appState
}
