package templates

import (
	"bytes"
	"fmt"
	"slices"
	"text/template"

	"github.com/SoenkeD/sc/src/generator/stage2"
	"github.com/SoenkeD/sc/src/utils"
)

func LoadTemplate(targetCodes, templatedCodes map[string]string, module, ctlDir, ctlName string) map[string]string {
	templates := map[string]string{}

	for actionID := range targetCodes {

		if actionTemplate, ok := templatedCodes[actionID]; ok {
			templates[actionID] = actionTemplate
		}
	}

	return templates
}

func CollectExisting(codes map[string]string, dir, filePostFix string) []string {
	var existing []string

	for ID := range codes {

		exists := utils.CheckStateMachineFileExists(dir, fmt.Sprintf("%s%s", ID, filePostFix))
		if exists {
			existing = append(existing, ID)
		}
	}
	return existing
}

type Codes struct {
	Actions         map[string]string
	ActionTests     map[string]string
	ActionTemplates map[string]string

	Guards         map[string]string
	GuardTests     map[string]string
	GuardTemplates map[string]string

	States []string
}

type filledTemplates struct {
	Actions     map[string]string
	ActionTests map[string]string

	Guards     map[string]string
	GuardTests map[string]string

	States []string
}

func executeDefaultTemplates(st2 stage2.ParseUmlStage2, tpl GenerateTemplatesInput) (filledTemplates, error) {

	actionCodes, err := GenerateActions(st2, tpl.TemplatedBaseFiles["oneaction"])
	if err != nil {
		return filledTemplates{}, err
	}
	actionTestCodes, err := GenerateActionTests(st2, tpl.TemplatedBaseFiles["oneaction_test"])
	if err != nil {
		return filledTemplates{}, err
	}

	guardCodes, err := GenerateGuards(st2, tpl.TemplatedBaseFiles["oneguard"])
	if err != nil {
		return filledTemplates{}, err
	}

	guardTestCodes, err := GenerateGuardTests(st2, tpl.TemplatedBaseFiles["oneguard_test"])
	if err != nil {
		return filledTemplates{}, err
	}

	stateCodes, err := GenerateStates(st2, tpl)
	if err != nil {
		return filledTemplates{}, err
	}

	return filledTemplates{
		Actions:     actionCodes,
		ActionTests: actionTestCodes,
		Guards:      guardCodes,
		GuardTests:  guardTestCodes,
		States:      stateCodes,
	}, nil
}

func removeExisting(defaults filledTemplates, actionsDir, guardsDir, language string) (filledTemplates, error) {

	languageFilePostfix := "." + language

	tpls := filledTemplates{
		Actions:     map[string]string{},
		ActionTests: map[string]string{},
		Guards:      map[string]string{},
		GuardTests:  map[string]string{},
		States:      defaults.States,
	}

	existingActions := CollectExisting(defaults.Actions, actionsDir, languageFilePostfix)
	for actionID, action := range defaults.Actions {

		if !slices.Contains(existingActions, actionID) {
			tpls.Actions[actionID] = action
		}

	}

	existingActionTests := CollectExisting(defaults.Actions, actionsDir, "_test"+languageFilePostfix)
	for actionID, action := range defaults.ActionTests {

		if !slices.Contains(existingActionTests, actionID) {
			tpls.ActionTests[actionID] = action
		}

	}

	existingGuards := CollectExisting(defaults.Guards, guardsDir, languageFilePostfix)
	for guardID, guard := range defaults.Guards {

		if !slices.Contains(existingGuards, guardID) {
			tpls.Guards[guardID] = guard
		}

	}

	existingGuardTests := CollectExisting(defaults.GuardTests, guardsDir, "_test"+languageFilePostfix)
	for guardID, guard := range defaults.GuardTests {

		if !slices.Contains(existingGuardTests, guardID) {
			tpls.GuardTests[guardID] = guard
		}

	}

	return tpls, nil
}

type GenerateTemplatesInput struct {
	TemplatedActions              map[string]string
	TemplatedGuards               map[string]string
	TemplatedStateExtensions      map[string]string
	TemplatedBaseFiles            map[string]string
	TemplatedControllerExtensions map[string]string
	TemplatedPerController        map[string]string
}

func ExecuteTemplates(st2 stage2.ParseUmlStage2, tpl GenerateTemplatesInput, repoRoot, module, ctlDir, ctlName, actionsDir, guardsDir, language string) (Codes, error) {

	languageFilePostfix := "." + language

	defaults, err := executeDefaultTemplates(st2, tpl)
	if err != nil {
		return Codes{}, err
	}

	// load custom templates
	templatedActionCodes := LoadTemplate(defaults.Actions, tpl.TemplatedActions, module, ctlDir, ctlName)
	templatedGuardCodes := LoadTemplate(defaults.Guards, tpl.TemplatedGuards, module, ctlDir, ctlName)

	// overwrite with existing
	clearedTpls, err := removeExisting(defaults, actionsDir, guardsDir, languageFilePostfix)
	if err != nil {
		return Codes{}, nil
	}

	existingActions := CollectExisting(defaults.Actions, actionsDir, languageFilePostfix)
	for _, actionID := range existingActions {
		defaults.Actions[actionID] = ""
	}
	existingActionTests := CollectExisting(defaults.Actions, actionsDir, "_test"+languageFilePostfix)
	for _, actionID := range existingActionTests {
		defaults.ActionTests[actionID] = ""
	}

	existingGuards := CollectExisting(defaults.Guards, guardsDir, languageFilePostfix)
	for _, guardID := range existingGuards {
		defaults.Guards[guardID] = ""
	}
	existingGuardTests := CollectExisting(defaults.GuardTests, guardsDir, "_test"+languageFilePostfix)
	for _, guardID := range existingGuardTests {
		defaults.GuardTests[guardID] = ""
	}

	return Codes{
		Actions:         clearedTpls.Actions,
		ActionTests:     clearedTpls.ActionTests,
		ActionTemplates: templatedActionCodes,
		Guards:          clearedTpls.Guards,
		GuardTests:      clearedTpls.GuardTests,
		GuardTemplates:  templatedGuardCodes,
		States:          clearedTpls.States,
	}, nil
}

func ExecTemplate(name, tpl string, input any, funcMap *template.FuncMap) (string, error) {

	if funcMap == nil {
		funcMap = &template.FuncMap{}
	}

	tmpl, err := template.New(name).Funcs(*funcMap).Parse(tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, input)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
