package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var template string
var extensions *[]string

func init() {
	extendCmd.Flags().StringVarP(&config.RepoRoot, "root", "r", "", "base path to generate ctl in")
	extendCmd.Flags().StringVarP(&config.Module, "module", "m", "", "name of the module e.g. github.com/SoenkeD/sc")
	extendCmd.Flags().StringVar(&template, "template", "", "relative path from repo root to the template file")
	extensions = extendCmd.Flags().StringSlice("extension", []string{}, "relative path from repo root to the template file")
	extendCmd.Flags().StringVarP(&ctlName, "name", "n", "", "name of the ctl")

	rootCmd.AddCommand(extendCmd)
}

func removeFirstLine(input string) string {
	lines := strings.Split(input, "\n")
	return strings.Join(lines[1:], "\n")
}

var extendCmd = &cobra.Command{
	Use:   "extend",
	Short: "Extend the .plantuml file for a controller",
	RunE: func(cmd *cobra.Command, args []string) error {

		// read template file
		templateBytes, err := os.ReadFile(template)
		if err != nil {
			return err
		}
		templateFile := removeFirstLine(string(templateBytes))

		targetFileLines := []string{
			"@startuml " + ctlName,
			templateFile,
		}

		// read extension files
		ctlDir := filepath.Join(config.RepoRoot, config.CtlDir, ctlName)
		for _, extension := range *extensions {
			extensionFilePath := filepath.Join(ctlDir, extension)
			extensionFile, err := os.ReadFile(extensionFilePath)
			if err != nil {
				return err
			}

			targetFileLines = append(targetFileLines, removeFirstLine(string(extensionFile)))
		}

		// write to target file
		targetFilePath := filepath.Join(ctlDir, ctlName+".plantuml")
		err = os.WriteFile(targetFilePath, []byte(strings.Join(targetFileLines, "\n")), 0777)
		if err != nil {
			return err
		}

		return nil
	},
}
