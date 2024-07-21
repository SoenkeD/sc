package templates

import (
	"log"
	"strings"
	"text/template"

	"github.com/SoenkeD/sc/src/generator/stage2"
	"github.com/SoenkeD/sc/src/types"
)

type GenerateStateTplInput struct {
	Name  string
	State stage2.ParsedState
}

func GenerateState(stateIdx, oneStateTpl string, state stage2.ParsedState) (string, error) {

	funcMap := template.FuncMap{
		"join":       strings.Join,
		"replaceAll": strings.ReplaceAll,
		"trans": func(inType types.TransitionType) string {
			var taType string
			switch inType {
			case types.TransitionTypeNormal:
				taType = "TransitionTypeNormal"
			case types.TransitionTypeHappy:
				taType = "TransitionTypeHappy"
			case types.TransitionTypeError:
				taType = "TransitionTypeError"
			default:
				taType = "TransitionTypeNormal"
				log.Println("WARNING: cannot understand transition type", inType)
			}

			return taType
		},
		"toUpper": strings.ToUpper,
		"typesToUpper": func(inType types.TransitionType) string {
			return strings.ToUpper(string(inType))
		},
	}

	code, err := ExecTemplate(stateIdx, oneStateTpl, GenerateStateTplInput{
		Name:  strings.ReplaceAll(stateIdx, "/", ""),
		State: state,
	}, &funcMap)
	if err != nil {
		return "", err
	}

	return code, nil
}

func GenerateStates(uml stage2.ParseUmlStage2, tpl GenerateTemplatesInput) ([]string, error) {

	stateCodes := []string{}

	for stateIdx, state := range uml.States {

		stateCode, err := GenerateState(stateIdx, tpl.TemplatedBaseFiles["onestate"], state)
		if err != nil {
			return nil, err
		}

		stateCodes = append(stateCodes, stateCode)
	}

	return stateCodes, nil
}
