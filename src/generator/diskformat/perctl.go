package diskformat

import (
	"log"

	"github.com/SoenkeD/sc/src/generator/templates"
)

func generatePerControllerAll(templatedPerController, perCtlDir map[string]string, ctlName, importRoot, module string) ([]GeneratedFile, error) {
	var files []GeneratedFile

	for perControllerId, perController := range templatedPerController {

		targetDir, ok := perCtlDir[perControllerId]
		if !ok {
			log.Printf("Info: skipping: failed to get a target dir for perController ID=%s\n", perControllerId)
			continue
		}

		file, err := generatePerController(perControllerId, perController, targetDir, ctlName, importRoot, module)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

type GeneratePerCtlTplInput struct {
	Name       string
	ImportRoot string
	Module     string
}

func generatePerController(perCtlId, perCtl, targetDir, ctlName, importRoot, module string) (GeneratedFile, error) {

	code, err := templates.ExecTemplate(perCtlId, perCtl, GeneratePerCtlTplInput{
		Name:       ctlName,
		ImportRoot: importRoot,
		Module:     module,
	}, getFuncMap())
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Type:                 "per_ctl",
		MarkAsGenerated:      true,
		Path:                 targetDir,
		PathStartsAtRepoRoot: true,
		Name:                 ctlName + "_" + perCtlId,
		Content:              []byte(code),
	}, nil
}
