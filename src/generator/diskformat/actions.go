package diskformat

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/SoenkeD/sc/src/generator/templates"
)

func generateActions(importRoot, actionTpl string, actions, templatedActions map[string]string) ([]GeneratedFile, error) {
	var files []GeneratedFile

	for actionID, action := range actions {

		if action == "" {
			continue
		}

		actionCode, err := generateAction(actionID, action, importRoot, actionTpl, templatedActions)
		if err != nil {
			return nil, err
		}
		files = append(files, actionCode)
	}

	return files, nil
}

type GenerateActionTplInput struct {
	ImportRoot string
	Code       string
}

func generateAction(actionID, actionCode, importRoot, actionTpl string, templatedActions map[string]string) (GeneratedFile, error) {

	code, err := templates.ExecTemplate(actionID, actionTpl, GenerateActionTestTplInput{
		ImportRoot: importRoot,
		Code:       actionCode,
	}, &template.FuncMap{
		"replaceAll": strings.ReplaceAll,
	})
	if err != nil {
		return GeneratedFile{}, err
	}

	var isTemplated bool
	if _, ok := templatedActions[actionID]; ok {

		tplCode, err := templates.ExecTemplate(actionID, templatedActions[actionID], GenerateActionTestTplInput{
			ImportRoot: importRoot,
		}, &template.FuncMap{})
		if err != nil {
			return GeneratedFile{}, err
		}

		code = tplCode
		isTemplated = true
	}

	return GeneratedFile{
		Type:            "action",
		MarkAsGenerated: isTemplated,
		Path:            "actions/",
		Name:            actionID,
		Content:         []byte(code),
	}, nil
}

func generateActionTests(importRoot, actionTestTpl string, actionTests, templatedActions map[string]string) ([]GeneratedFile, error) {
	var files []GeneratedFile

	for actionID, action := range actionTests {

		if action == "" {
			continue
		}

		// skip templated actions
		if _, ok := templatedActions[actionID]; ok {
			continue
		}

		actionTest, err := generateActionTest(actionID, action, importRoot, actionTestTpl)
		if err != nil {
			return nil, err
		}

		files = append(files, actionTest)
	}

	return files, nil
}

type GenerateActionTestTplInput struct {
	ImportRoot string
	Code       string
}

func generateActionTest(actionID, actionTest, importRoot, actionTestTpl string) (GeneratedFile, error) {

	code, err := templates.ExecTemplate(actionID, actionTestTpl, GenerateActionTestTplInput{
		ImportRoot: importRoot,
		Code:       actionTest,
	}, getFuncMap())
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Type:            "action_test",
		MarkAsGenerated: false,
		Path:            "actions/",
		Name:            fmt.Sprintf("%s_test", actionID),
		Content:         []byte(code),
	}, nil
}
