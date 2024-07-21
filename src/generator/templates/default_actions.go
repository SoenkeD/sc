package templates

import (
	"github.com/SoenkeD/sc/src/generator/stage2"
)

type GenerateActionTplInput struct {
	Name string
}

func GenerateActions(uml stage2.ParseUmlStage2, oneActionTpl string) (map[string]string, error) {
	actionCodes := make(map[string]string, len(uml.Actions))
	for _, action := range uml.Actions {

		code, err := ExecTemplate(action, oneActionTpl, GenerateActionTplInput{
			Name: action,
		}, nil)
		if err != nil {
			return nil, err
		}

		actionCodes[action] = code
	}

	return actionCodes, nil
}

type GenerateActionTestTplInput struct {
	Name string
}

func GenerateActionTests(uml stage2.ParseUmlStage2, oneActionTestTpl string) (map[string]string, error) {
	actionTestCodes := make(map[string]string, len(uml.Actions))
	for _, action := range uml.Actions {

		code, err := ExecTemplate(action, oneActionTestTpl, GenerateActionTestTplInput{
			Name: action,
		}, nil)
		if err != nil {
			return nil, err
		}

		actionTestCodes[action] = code
	}

	return actionTestCodes, nil
}
