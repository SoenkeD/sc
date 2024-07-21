package diskformat

import (
	"html/template"
	"path/filepath"
	"strings"
	"unicode"
)

type GenerationInput struct {
	CtlName                  string
	RepoRoot                 string
	RelativeCtlRoot          string
	Module                   string
	Actions                  map[string]string
	TemplatedActions         map[string]string
	ActionTests              map[string]string
	Guards                   map[string]string
	TemplatedGuards          map[string]string
	GuardTests               map[string]string
	TemplatedStateExtensions map[string]string
	States                   []string
}

type GeneratedFile struct {
	Type            string
	MarkAsGenerated bool
	Name            string
	Path            string
	Content         []byte
	ForceWrite      bool
}

func (file GeneratedFile) GetFilePath(fileExtension string, enableGeneratedFileExtension, enableFileCapitalization bool) string {
	fileName := file.Name

	// capitalize first letter
	if enableFileCapitalization && len(file.Name) > 0 {
		fileName = string(unicode.ToUpper([]rune(file.Name)[0])) + file.Name[1:]
	}

	// add file extension
	fileName += "." + fileExtension

	// add generation prefix
	if file.MarkAsGenerated && enableGeneratedFileExtension {
		fileName = "zz_gen_" + fileName
	}
	return filepath.Join(file.Path, fileName)
}

type Generation struct {
	BasePath string
	Dirs     []string
	Files    []GeneratedFile
}

func getFuncMap() *template.FuncMap {
	return &template.FuncMap{
		"replaceAll": strings.ReplaceAll,
		"toUpper":    strings.ToUpper,
	}
}

func TransformImport(importPath, separator string) string {
	return strings.ReplaceAll(importPath, "/", separator)
}
