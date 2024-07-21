package diskformat

import (
	"github.com/SoenkeD/sc/src/generator/templates"
)

type GenerateStateMachineTplInput struct {
	ImportRoot string
	Name       string
	States     []string
}

func generateStateMachine(ctlName, importRoot, stateMachineTpl string, states []string) (GeneratedFile, error) {

	code, err := templates.ExecTemplate(ctlName, stateMachineTpl, GenerateStateMachineTplInput{
		Name:       ctlName,
		ImportRoot: importRoot,
		States:     states,
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

type GenerateBaseReconcilerTplInput struct {
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

	reconcilerFile, err := templates.ExecTemplate("reconciler", tplIn.TemplatedBaseFiles["reconciler"], GenerateBaseReconcilerTplInput{
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
		GeneratedFile{
			Path:            ctlDirName,
			Type:            "ctl",
			MarkAsGenerated: true,
			Name:            "reconciler",
			Content:         []byte(reconcilerFile),
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
