package aftercompilation

import (
	"log"
	"os"
	"path/filepath"

	"github.com/SoenkeD/sc/src/generator/diskformat"
)

func WriteToDisk(gen diskformat.Generation, fileExtension string, enableGeneratedFileExtension, enableFileCapitalization, forceWriteGenerated bool) error {

	for _, dir := range gen.Dirs {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	for _, file := range gen.Files {

		if forceWriteGenerated && file.MarkAsGenerated {
			file.ForceWrite = true
		}

		genFilePath := file.GetFilePath(fileExtension, enableGeneratedFileExtension, enableFileCapitalization)
		filePath := filepath.Join(gen.BasePath, genFilePath)
		if file.PathStartsAtRepoRoot {
			filePath = filepath.Join(gen.RepoRoot, genFilePath)
		}

		if !file.ForceWrite {
			_, err := os.Stat(filePath)
			if !os.IsNotExist(err) {
				continue
			}
		}

		if len(file.Content) == 0 {
			log.Printf("INFO: skipping empty file %s\n", filePath)
			continue
		}

		log.Printf("INFO: writing file %s\n", filePath)

		err := os.WriteFile(filePath, file.Content, os.ModePerm)
		if err != nil {
			return err
		}

	}

	return nil
}
