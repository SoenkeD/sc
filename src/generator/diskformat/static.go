package diskformat

import (
	"github.com/SoenkeD/sc/src/generator/templates"
)

type GenerateStateMachineTplInput struct {
	ImportRoot string
	Name       string
	States     []string
	HasActions bool
	HasGuards  bool
}

func generateStateMachine(ctlName, importRoot, stateMachineTpl string, states []string, hasActions, hasGuards bool) (GeneratedFile, error) {

	code, err := templates.ExecTemplate(ctlName, stateMachineTpl, GenerateStateMachineTplInput{
		Name:       ctlName,
		ImportRoot: importRoot,
		States:     states,
		HasActions: hasActions,
		HasGuards:  hasGuards,
	}, getFuncMap())
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Type:            "ctl",
		MarkAsGenerated: true,
		Name:            "sm",
		Content:         []byte(code),
		ForceWrite:      true,
	}, nil
}

type GenerateBaseCtlTplInput struct {
	ImportRoot string
}

type GenerateBaseTypesTplInput struct {
	ImportRoot string
}

type GenerateBaseUtilsTplInput struct {
	ImportRoot string
}

type GenerateStateExtensionsTplInput struct {
	ImportRoot string
}

func generateCtlFiles(ctlDirName, importRoot string, tplIn templates.GenerateTemplatesInput, separator string) ([]GeneratedFile, error) {
	var files []GeneratedFile

	ctlFile, err := templates.ExecTemplate("ctl", tplIn.TemplatedBaseFiles["ctl"], GenerateBaseCtlTplInput{
		ImportRoot: TransformImport(importRoot, separator),
	}, getFuncMap())
	if err != nil {
		return nil, err
	}

	files = append(
		files,
		GeneratedFile{
			Path:            ctlDirName,
			Type:            "ctl",
			MarkAsGenerated: true,
			Name:            "ctl",
			Content:         []byte(ctlFile),
		},
	)

	for file, fileContent := range tplIn.TemplatedControllerExtensions {

		ctlFile, err := templates.ExecTemplate(file, fileContent, GenerateStateExtensionsTplInput{
			ImportRoot: TransformImport(importRoot, separator),
		}, getFuncMap())
		if err != nil {
			return nil, err
		}

		files = append(files, GeneratedFile{
			Path:            ctlDirName,
			Type:            "ctl",
			MarkAsGenerated: true,
			Name:            file,
			Content:         []byte(ctlFile),
		})

	}

	return files, nil
}
