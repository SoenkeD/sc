package diskformat

func generateStateExtensions(templatedStateExtensions map[string]string, language string) []GeneratedFile {
	var files []GeneratedFile

	for extensionId, extension := range templatedStateExtensions {

		files = append(files, generateStateExtension(extensionId, extension, language))
	}

	return files
}

func generateStateExtension(extensionId, extension, language string) GeneratedFile {

	return GeneratedFile{
		Type:            "state",
		MarkAsGenerated: true,
		Path:            "state/",
		Name:            extensionId + "." + language,
		Content:         []byte(extension),
	}
}
