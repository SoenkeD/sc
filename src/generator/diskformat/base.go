package diskformat

import (
	"github.com/SoenkeD/sc/src/generator/templates"
)

type GenerateActionsTplInput struct {
	ImportRoot string
	Actions    map[string]string
}

func generateBaseActions(importRoot, actionsTpl string, actions map[string]string, forceWrite bool) (GeneratedFile, error) {

	code, err := templates.ExecTemplate("actions", actionsTpl, GenerateActionsTplInput{
		ImportRoot: importRoot,
		Actions:    actions,
	}, getFuncMap())
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Type:            "actions",
		MarkAsGenerated: true,
		Path:            "actions/",
		Name:            "actions",
		ForceWrite:      forceWrite,
		Content:         []byte(code),
	}, nil
}

func generateBaseActionsTest(actionTestsTpl string) GeneratedFile {
	return GeneratedFile{
		Type:            "actions_test",
		MarkAsGenerated: true,
		Path:            "actions/",
		Name:            "actions_test",
		Content:         []byte(actionTestsTpl),
	}
}

type GenerateGuardsTplInput struct {
	ImportRoot string
	Guards     map[string]string
}

func generateBaseGuards(importRoot, guardsTpl string, guards map[string]string, forceWrite bool) (GeneratedFile, error) {

	code, err := templates.ExecTemplate("guards", guardsTpl, GenerateGuardsTplInput{
		ImportRoot: importRoot,
		Guards:     guards,
	}, getFuncMap())
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Type:            "guards",
		MarkAsGenerated: true,
		Path:            "guards/",
		Name:            "guards",
		ForceWrite:      forceWrite,
		Content:         []byte(code),
	}, nil
}

func generateBaseGuardsTest(guardsTestTpl string) GeneratedFile {
	return GeneratedFile{
		Type:            "guards_test",
		MarkAsGenerated: true,
		Path:            "guards/",
		Name:            "guards_test",
		Content:         []byte(guardsTestTpl),
	}
}
