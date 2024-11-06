package stage2

import (
	"fmt"
	"strings"

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
	Options      []string
}

func (tr *ParseTransitionStage2) GetId() string {

	negationStr := "false"
	if tr.Negation {
		negationStr = "true"
	}

	return strings.Join(
		[]string{
			tr.Action,
			strings.Join(tr.ActionParams, ","),
			tr.Guard,
			strings.Join(tr.GuardParams, ","),
			string(tr.Type),
			strings.ReplaceAll(tr.Target, "/", ""),
			negationStr,
			strings.Join(tr.Options, ","),
		},
		"/",
	)
}

type ParsedStateActionStage2 struct {
	Action       string
	ActionParams []string
}

type ParsedState struct {
	Actions     []ParsedStateActionStage2
	Transitions []ParseTransitionStage2
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

func ExtractVisitedTransactions(route []string) (states, transitions []string) {

	for _, station := range route {

		if strings.HasSuffix(station, "State") {

			nextState := strings.TrimSuffix(station, "State")
			states = append(states, nextState)
		}

		// is a transition
		if strings.HasSuffix(station, "/true") || strings.HasSuffix(station, "/false") {
			parts := strings.Split(station, "/")

			parts[0] = strings.TrimSuffix(parts[0], "State")
			parts[1] = strings.TrimSuffix(parts[1], "Action")
			parts[3] = strings.TrimSuffix(parts[3], "Guard")
			parts[6] = strings.TrimSuffix(parts[6], "State")

			transitions = append(transitions, "/"+strings.Join(parts, "/"))
		}
	}

	return
}

func PrintTransitionType(taType types.TransitionType, color string) string {

	options := []string{}

	if taType == types.TransitionTypeHappy {
		options = append(options, "bold")
	}

	if taType == types.TransitionTypeError {
		options = append(options, "dotted")
	}

	if color != "" {
		options = append(options, color)
	}

	optionsStr := ""
	if len(options) > 0 {
		optionsStr = "[" + strings.Join(options, ",") + "]"
	}

	return fmt.Sprintf("-%s->", optionsStr)
}
