package templates

import (
	"github.com/SoenkeD/sc/src/generator/stage2"
)

type GenerateGuardTplInput struct {
	Name string
}

func GenerateGuards(uml stage2.ParseUmlStage2, oneGuardTpl string) (map[string]string, error) {
	guardCodes := make(map[string]string, len(uml.Guards))
	for _, guard := range uml.Guards {

		code, err := ExecTemplate(guard, oneGuardTpl, GenerateGuardTplInput{
			Name: guard,
		}, nil)
		if err != nil {
			return nil, err
		}

		guardCodes[guard] = code
	}

	return guardCodes, nil
}

type GenerateGuardTestTplInput struct {
	Name string
}

func GenerateGuardTests(uml stage2.ParseUmlStage2, oneGuardTestTpl string) (map[string]string, error) {
	guardTestCodes := make(map[string]string, len(uml.Guards))
	for _, guard := range uml.Guards {

		code, err := ExecTemplate(guard, oneGuardTestTpl, GenerateGuardTestTplInput{
			Name: guard,
		}, nil)
		if err != nil {
			return nil, err
		}

		guardTestCodes[guard] = code
	}

	return guardTestCodes, nil
}
