package aftercompilation

import (
	"log"
	"path/filepath"
	"slices"
	"strings"

	"github.com/SoenkeD/sc/src/utils"
)

func CollectUnnecessary(dir string, codeIDs []string, fileExtension string, excluded []string) ([]string, error) {

	// to lower
	for idx, codeID := range codeIDs {
		codeIDs[idx] = strings.ToLower(codeID)
	}

	for idx, exclude := range excluded {
		excluded[idx] = strings.ToLower(exclude)
	}

	exists, err := utils.FileOrDirExists(dir)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	files, err := utils.ListFilesInDir(dir)
	if err != nil {
		return nil, err
	}
	var unusedFiles []string
	for _, file := range files {
		fileCodeID := strings.TrimSuffix(strings.TrimSuffix(strings.TrimPrefix(file, "zz_gen_"), "."+fileExtension), "_test")
		fileCodeID = strings.ToLower(fileCodeID)

		if !slices.Contains(codeIDs, fileCodeID) && !slices.Contains(excluded, fileCodeID) {
			unusedFiles = append(unusedFiles, filepath.Join(dir, file))
		}
	}

	return unusedFiles, nil
}

func WarnUnnecessaryFiles(actionsDir string, actions []string, guardsDir string, guards []string, fileExtension string) error {

	notNeededActionFiles, err := CollectUnnecessary(actionsDir, actions, fileExtension, []string{"actions"})
	if err != nil {
		return err
	}

	notNeededGuardFiles, err := CollectUnnecessary(guardsDir, guards, fileExtension, []string{"guards"})
	if err != nil {
		return err
	}

	if len(notNeededActionFiles) > 0 {
		log.Printf("INFO: Can remove these action files: \n%s\n", strings.Join(notNeededActionFiles, "\n"))
	}

	if len(notNeededGuardFiles) > 0 {
		log.Printf("INFO: Can remove these guard files: \n%s\n", strings.Join(notNeededGuardFiles, "\n"))
	}

	return nil
}
