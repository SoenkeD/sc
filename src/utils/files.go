package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func FileOrDirExists(path string) (exists bool, err error) {

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		exists = false
		err = nil
		return
	}

	if err != nil {
		return
	}

	exists = true

	return
}

func ListFilesInDir(dir string) (fileNames []string, err error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	fileNames = make([]string, 0, len(files))
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return
}

func ListFilesInDirRecursive(dir string) (fileNames []string, err error) {

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range dirEntries {
		entryPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			subDirFiles, err := ListFilesInDirRecursive(entryPath)
			if err != nil {
				return nil, err
			}
			fileNames = append(fileNames, subDirFiles...)
		} else {
			fileNames = append(fileNames, entryPath)
		}
	}

	return
}

func ReadTplFilesInDir(dir string, writeMap map[string]string) error {

	exists, err := FileOrDirExists(dir)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	tplFiles, err := ListFilesInDir(dir)
	if err != nil {
		return err
	}

	for _, tplFile := range tplFiles {

		if !strings.HasSuffix(tplFile, ".tpl") {
			continue
		}

		fullTplPath := filepath.Join(dir, tplFile)
		tplBytes, readErr := os.ReadFile(fullTplPath)
		if readErr != nil {
			return readErr
		}

		actionName := strings.TrimSuffix(tplFile, ".tpl")
		writeMap[actionName] = string(tplBytes)
	}

	return nil
}

func CheckStateMachineFileExists(filePath, filename string) bool {

	fileName := filepath.Join(filePath, filename)
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		return true
	}

	fileName = filepath.Join(filePath, "zz_gen_"+filename)
	_, err = os.Stat(fileName)
	return !os.IsNotExist(err)
}

func CreateDirs(dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func WriteFile(filePath, fileStr string) error {
	err := os.WriteFile(filePath, []byte(fileStr), 0755)
	if err != nil {
		return err
	}

	return nil
}
