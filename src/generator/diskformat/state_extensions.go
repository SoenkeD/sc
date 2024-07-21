package diskformat

func generateStateExtensions(templatedStateExtensions map[string]string) []GeneratedFile {
	var files []GeneratedFile

	for extensionId, extension := range templatedStateExtensions {

		files = append(files, generateStateExtension(extensionId, extension))
	}

	return files
}

func generateStateExtension(extensionId, extension string) GeneratedFile {

	return GeneratedFile{
		Type:            "state",
		MarkAsGenerated: true,
		Path:            "state/",
		Name:            extensionId + ".go",
		Content:         []byte(extension),
	}
}
