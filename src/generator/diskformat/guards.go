package diskformat

import (
	"fmt"

	"github.com/SoenkeD/sc/src/generator/templates"
)

func generateGuards(importRoot, guardTpl string, guards, templatedGuards map[string]string) ([]GeneratedFile, error) {
	var files []GeneratedFile

	for guardID, guard := range guards {

		if guard == "" {
			continue
		}

		guardCode, err := generateGuard(guardID, guard, importRoot, guardTpl, templatedGuards)
		if err != nil {
			return nil, err
		}
		files = append(files, guardCode)
	}

	return files, nil
}

type GenerateGuardTplInput struct {
	ImportRoot string
	Code       string
}

func generateGuard(guardID, guardCode, importRoot, guardTpl string, templatedGuards map[string]string) (GeneratedFile, error) {

	code, err := templates.ExecTemplate(guardID, guardTpl, GenerateActionTestTplInput{
		ImportRoot: importRoot,
		Code:       guardCode,
	}, getFuncMap())
	if err != nil {
		return GeneratedFile{}, err
	}

	var isTemplated bool
	if _, ok := templatedGuards[guardID]; ok {

		code, err = templates.ExecTemplate(guardID, templatedGuards[guardID], GenerateActionTestTplInput{
			ImportRoot: importRoot,
			Code:       guardCode,
		}, getFuncMap())
		if err != nil {
			return GeneratedFile{}, err
		}

		isTemplated = true
	}
	return GeneratedFile{
		Type:            "guard",
		MarkAsGenerated: isTemplated,
		Path:            "guards/",
		Name:            guardID,
		Content:         []byte(code),
	}, nil
}

func generateGuardTests(importRoot, guardTestTpl string, guardTests, templatedGuards map[string]string) ([]GeneratedFile, error) {
	var files []GeneratedFile

	for guardID, guard := range guardTests {

		if guard == "" {
			continue
		}

		// skip templated guards
		if _, ok := templatedGuards[guardID]; ok {
			continue
		}

		guardCodes, err := generateGuardTest(guardID, guard, importRoot, guardTestTpl)
		if err != nil {
			return nil, err
		}
		files = append(files, guardCodes)
	}

	return files, nil
}

func generateGuardTest(guardID, guardTest, importRoot, guardTestTpl string) (GeneratedFile, error) {

	code, err := templates.ExecTemplate(guardID, guardTestTpl, GenerateActionTestTplInput{
		ImportRoot: importRoot,
		Code:       guardTest,
	}, getFuncMap())
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Type:            "guard_test",
		MarkAsGenerated: false,
		Path:            "guards/",
		Name:            fmt.Sprintf("%s_test", guardID),
		Content:         []byte(code),
	}, nil
}
