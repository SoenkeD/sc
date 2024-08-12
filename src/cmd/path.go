package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/SoenkeD/sc/src/generator"
	"github.com/SoenkeD/sc/src/generator/parseuml"
	"github.com/SoenkeD/sc/src/generator/stage2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var routeFile string

func init() {
	pathCmd.Flags().StringVarP(&config.RepoRoot, "root", "r", "", "base path to generate ctl in")
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))

	pathCmd.Flags().StringVarP(&config.Module, "module", "m", "", "name of the module e.g. github.com/SoenkeD/sc")
	viper.BindPFlag("module", rootCmd.PersistentFlags().Lookup("module"))

	pathCmd.Flags().StringVarP(&ctlName, "name", "n", "", "name of the ctl")

	pathCmd.Flags().StringVarP(&routeFile, "route", "i", "", "path to the route input file (contains JSON with key \"route\")")

	rootCmd.AddCommand(pathCmd)
}

var tab = strings.Repeat(" ", 4)

func renderUml(name string, lines []string) string {
	umlLines := []string{
		"@startuml " + name,
		"<style>\n    .visitedState {\n        FontColor Green\n    }\n</style>\n",
	}

	umlLines = append(umlLines, lines...)

	return strings.Join(umlLines, "\n")
}

func getOpenCompoundStateLines(steps, route []string, stateIdx, usedCurrentStateName string) (lines []string, parentCalled bool) {

	for stepIdx, step := range steps {
		if stepIdx == len(steps)-1 {

			if strings.HasSuffix(stateIdx, "End") || strings.HasSuffix(stateIdx, "Start") {
				continue
			}

			var visitColor string
			if slices.Contains(route, strings.ReplaceAll(stateIdx, "/", "")+"State") {
				visitColor = "<<visitedState>>"
				parentCalled = true
			}

			lines = append(lines, fmt.Sprintf("%sstate %s %s", strings.Repeat(tab, stepIdx-1), usedCurrentStateName, visitColor))

			break
		}

		if step == "" {
			continue
		}

		var visitColor string
		if slices.Contains(route, strings.Join(steps[0:stepIdx+1], "")+"StartState") {
			visitColor = "#Green"
		}

		lines = append(lines, fmt.Sprintf("%sstate %s %s{", strings.Repeat(tab, stepIdx-1), step, visitColor))
	}

	return
}

func getCloseCompoundStateLines(steps, writeAfterLines []string) (lines []string) {

	for stepIdx, step := range steps {
		if stepIdx == len(steps)-1 {
			break
		}
		if step == "" { // on root level
			continue
		}
		lines = append(lines, fmt.Sprintf("%s}", strings.Repeat(tab, len(steps)-stepIdx-2)))

		if stepIdx == 1 {
			lines = append(lines, writeAfterLines...)
		}
	}

	return
}

func getStateActionUmlLines(state stage2.ParsedState, steps, route []string, parentCalled bool, usedCurrentStateName string) (lines []string) {

	var prevActionFailed bool
	for _, action := range state.Actions {

		var actionColor string
		if prevActionFailed {
			actionColor = "<color:Black>"
		} else if parentCalled && route[len(route)-1] == action.Action+"Action" {
			actionColor = "<color:Red>"
			prevActionFailed = true
		}

		stateActionStr := fmt.Sprintf("%s%s: %s%s%s", strings.Repeat(tab, len(steps)-2), usedCurrentStateName, actionColor, parseuml.DO_SLASH_LINE, action.Action)

		if len(action.ActionParams) > 0 {
			stateActionStr += fmt.Sprintf("(%s)", strings.Join(action.ActionParams, ", "))
		}

		lines = append(lines, stateActionStr)
	}

	return
}

func getTransitionUmlLine(transitions []string, transition stage2.ParseTransitionStage2, stateIdx, usedCurrentStateName string) string {
	var wasVisited bool

	var arrowColor string

	if slices.Contains(transitions, "/"+strings.ReplaceAll(stateIdx, "/", "")+"/"+transition.GetId()) {
		arrowColor = "#Green"
		wasVisited = true
	}
	arrowStr := stage2.PrintTransitionType(transition.Type, arrowColor)

	transitionParts := strings.Split(transition.Target, "/")

	targetName := transitionParts[len(transitionParts)-1]
	if targetName == "End" {
		targetName = "[*]"
	} else if targetName == "Start" {
		targetName = transitionParts[len(transitionParts)-2]
	}

	var guardStr string
	if transition.Guard != "" {

		var negationStr string
		if transition.Negation {
			negationStr = "! "
		}

		var guardParamStr string
		if len(transition.GuardParams) > 0 {
			guardParamStr = fmt.Sprintf("(%s)", strings.Join(transition.GuardParams, ", "))
		}

		guardStr = fmt.Sprintf("[ %s%s%s ]", negationStr, transition.Guard, guardParamStr)
	}

	var actionStr string
	if transition.Action != "" {
		var actionParamStr string
		if len(transition.ActionParams) > 0 {
			actionParamStr = fmt.Sprintf("(%s)", strings.Join(transition.ActionParams, ", "))
		}
		actionStr = fmt.Sprintf(" / %s%s", transition.Action, actionParamStr)
	}

	delimiter := ""
	if transition.Action != "" || transition.Guard != "" {
		delimiter = ": "
		if wasVisited {
			delimiter += "<color:Green> "
		}
	}

	if wasVisited && usedCurrentStateName == "[*]" {
		usedCurrentStateName += "#Green"
	}

	return fmt.Sprintf("%s %s %s%s%s %s", usedCurrentStateName, arrowStr, targetName, delimiter, guardStr, actionStr)
}

func getStateUmlLines(stateIdx string, state stage2.ParsedState, route, transitions []string) []string {

	var lines []string

	var writeAfterLines []string

	steps := strings.Split(stateIdx, "/")
	currentStateName := steps[len(steps)-1]
	usedCurrentStateName := currentStateName
	if currentStateName == "Start" {
		usedCurrentStateName = "[*]"
	}

	if currentStateName == "End" {
		usedCurrentStateName = steps[len(steps)-2]
	}

	// open compound states
	openCompoundLines, parentCalled := getOpenCompoundStateLines(steps, route, stateIdx, usedCurrentStateName)
	lines = append(lines, openCompoundLines...)

	// state actions
	lines = append(lines, getStateActionUmlLines(state, steps, route, parentCalled, usedCurrentStateName)...)

	// transitions
	for _, transition := range state.Transitions {

		add2UmlStr := getTransitionUmlLine(transitions, transition, stateIdx, usedCurrentStateName)
		if currentStateName == "End" {
			writeAfterLines = append(writeAfterLines, strings.Repeat(tab, len(steps)-3)+add2UmlStr)
		} else {
			lines = append(lines, strings.Repeat(tab, len(steps)-2)+add2UmlStr)
		}

	}

	// close compound states
	lines = append(lines, getCloseCompoundStateLines(steps, writeAfterLines)...)

	return lines
}

type RouteInputJson struct {
	Route []string `json:"route"`
}

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Colorize the path of a controller run",
	RunE: func(cmd *cobra.Command, args []string) error {

		// read route infos
		routeFileBytes, err := os.ReadFile(routeFile)
		if err != nil {
			return err
		}

		var input RouteInputJson
		err = json.Unmarshal(routeFileBytes, &input)
		if err != nil {
			return err
		}

		// prepare generation input
		umlFile := filepath.Join(config.RepoRoot, config.CtlDir, ctlName, ctlName+".plantuml")
		code, err := generator.ParseUmlFile(umlFile)
		if err != nil {
			return err
		}

		_, transitions := stage2.ExtractVisitedTransactions(input.Route)

		// generate uml code
		var lines []string
		for stateIdx, state := range code.States {

			stateLines := getStateUmlLines(stateIdx, state, input.Route, transitions)

			lines = append(lines, stateLines...)
			lines = append(lines, "")
		}

		umlFileContent := renderUml(code.Name, lines)

		// write output
		err = os.WriteFile(ctlName+".route.plantuml", []byte(umlFileContent), 0777)
		if err != nil {
			return err
		}

		return nil
	},
}
