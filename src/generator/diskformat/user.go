package diskformat

import (
	"github.com/SoenkeD/sc/src/generator/templates"
)

type GenerateCtlTplInput struct {
	Name       string
	ImportRoot string
}

func generateCtl(ctlName, importRoot, initCtlTpl string) (GeneratedFile, error) {

	code, err := templates.ExecTemplate(ctlName, initCtlTpl, GenerateCtlTplInput{
		Name:       ctlName,
		ImportRoot: importRoot,
	}, getFuncMap())
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Type:            "ctl",
		MarkAsGenerated: true,
		Name:            "initctl",
		Content:         []byte(code),
	}, nil
}

type GenerateStateTplInput struct {
	Name       string
	ImportRoot string
}

func generateState(stateTpl, importRoot string) (GeneratedFile, error) {

	code, err := templates.ExecTemplate("state", stateTpl, GenerateStateTplInput{
		ImportRoot: importRoot,
	}, getFuncMap())
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Type:            "state",
		MarkAsGenerated: false,
		Path:            "state/",
		Name:            "ExtendedState",
		Content:         []byte(code),
	}, nil
}

type GenerateCtxTplInput struct {
	ImportRoot string
}

func generateCtx(ctxTpl, importRoot string) (GeneratedFile, error) {

	code, err := templates.ExecTemplate("context", ctxTpl, GenerateCtxTplInput{
		ImportRoot: importRoot,
	}, getFuncMap())
	if err != nil {
		return GeneratedFile{}, err
	}

	return GeneratedFile{
		Type:            "context",
		MarkAsGenerated: false,
		Path:            "state/",
		Name:            "Ctx",
		Content:         []byte(code),
	}, nil
}
