package cmd

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/SoenkeD/sc/src/generator/diskformat"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	exportCmd.Flags().StringVarP(&config.RepoRoot, "root", "r", "", "base path to generate ctl in")
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))

	exportCmd.Flags().StringVarP(&config.Module, "module", "m", "", "name of the module e.g. github.com/SoenkeD/sc")
	viper.BindPFlag("module", rootCmd.PersistentFlags().Lookup("module"))

	exportCmd.Flags().StringVarP(&ctlName, "name", "n", "", "name of the ctl")

	rootCmd.AddCommand(exportCmd)
}

type exportFile struct {
	Controller    string
	SrcFile       string
	TargetFile    string
	TargetContent []byte
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export to templates",
	RunE: func(cmd *cobra.Command, args []string) error {

		var exportFiles []*exportFile

		// collect export files
		for _, export := range config.Exports {
			ctlPath := filepath.Join(config.CtlDir, export.Controller)

			for _, item := range export.Items {

				suffix := "/*"
				if strings.HasSuffix(item.Src, suffix) {

					srcDir := filepath.Join(ctlPath, item.Src[:len(item.Src)-len(suffix)])
					log.Println(item.Src, srcDir)

					files, err := os.ReadDir(srcDir)
					if err != nil {
						return err
					}

					for _, entry := range files {
						if !entry.IsDir() {
							fileName := entry.Name()
							fileExt := filepath.Ext(fileName)

							if strings.HasSuffix(fileName, "_test"+fileExt) {
								continue
							}

							if fileExt == "."+config.Language {

								if slices.Contains(item.Excluded, fileName) {
									continue
								}

								exportFiles = append(exportFiles, &exportFile{
									Controller: export.Controller,
									SrcFile:    filepath.Join(srcDir, fileName),
									TargetFile: filepath.Join(
										item.To,
										strings.ReplaceAll(
											fileName,
											fileExt,
											".tpl",
										),
									),
								})
							}
						}
					}
				} else {
					srcPath := filepath.Join(ctlPath, item.Src)
					extractedFilenameParts := strings.Split(srcPath, "/")
					extractedFilename := extractedFilenameParts[len(extractedFilenameParts)-1]

					exportFiles = append(exportFiles, &exportFile{
						Controller: export.Controller,
						SrcFile:    srcPath,
						TargetFile: filepath.Join(
							item.To,
							strings.ReplaceAll(
								extractedFilename,
								"."+config.Language,
								".tpl",
							),
						),
					})
				}

			}
		}

		// prepare files
		for _, exp := range exportFiles {
			srcContentBytes, err := os.ReadFile(exp.SrcFile)
			if err != nil {
				return err
			}

			srcContent := string(srcContentBytes)

			srcImport := filepath.Join(config.Module, config.CtlDir, exp.Controller)
			srcImport = diskformat.TransformImport(srcImport, config.ImportPathSeparator)

			exp.TargetContent = []byte(strings.ReplaceAll(srcContent, srcImport, "{{ .ImportRoot }}"))
		}

		// write to disk
		for _, exp := range exportFiles {
			err := os.WriteFile(
				exp.TargetFile,
				exp.TargetContent,
				0777,
			)
			if err != nil {
				return err
			}
		}

		return nil
	},
}
