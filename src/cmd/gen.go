package cmd

import (
	"path/filepath"

	"github.com/SoenkeD/sc/src/generator"
	"github.com/SoenkeD/sc/src/generator/templates"
	"github.com/SoenkeD/sc/src/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ctlName string
var clearUnnecessary bool

func init() {
	genCmd.Flags().StringVarP(&config.RepoRoot, "root", "r", "", "base path to generate ctl in")
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))

	genCmd.Flags().StringVarP(&config.Module, "module", "m", "", "name of the module e.g. github.com/SoenkeD/sc")
	viper.BindPFlag("module", rootCmd.PersistentFlags().Lookup("module"))

	genCmd.Flags().StringVarP(&ctlName, "name", "n", "", "name of the ctl")

	genCmd.Flags().BoolVarP(&clearUnnecessary, "clear", "u", false, "if true it deletes unnecessary actions & guards")

	rootCmd.AddCommand(genCmd)
}

func ReadTemplates() (templates.GenerateTemplatesInput, error) {

	out := templates.GenerateTemplatesInput{
		TemplatedActions:              map[string]string{},
		TemplatedGuards:               map[string]string{},
		TemplatedStateExtensions:      map[string]string{},
		TemplatedBaseFiles:            map[string]string{},
		TemplatedControllerExtensions: map[string]string{},
	}

	for _, tplPackage := range config.Templates {
		actionsPath := filepath.Join(tplPackage.Dir, "actions")
		err := utils.ReadTplFilesInDir(actionsPath, out.TemplatedActions)
		if err != nil {
			return templates.GenerateTemplatesInput{}, err
		}

		guardsPath := filepath.Join(tplPackage.Dir, "guards")
		err = utils.ReadTplFilesInDir(guardsPath, out.TemplatedGuards)
		if err != nil {
			return templates.GenerateTemplatesInput{}, err
		}

		stateExtensionsPath := filepath.Join(tplPackage.Dir, "state")
		err = utils.ReadTplFilesInDir(stateExtensionsPath, out.TemplatedStateExtensions)
		if err != nil {
			return templates.GenerateTemplatesInput{}, err
		}

		baseFilesPath := filepath.Join(tplPackage.Dir, "base")
		err = utils.ReadTplFilesInDir(baseFilesPath, out.TemplatedBaseFiles)
		if err != nil {
			return templates.GenerateTemplatesInput{}, err
		}

		controllerExtensionsFilesPath := filepath.Join(tplPackage.Dir, "controller")
		err = utils.ReadTplFilesInDir(controllerExtensionsFilesPath, out.TemplatedControllerExtensions)
		if err != nil {
			return templates.GenerateTemplatesInput{}, err
		}
	}

	return out, nil
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate the state machine code",
	RunE: func(cmd *cobra.Command, args []string) error {

		// read from template packages
		templates, err := ReadTemplates()
		if err != nil {
			return err
		}

		err = generator.Generate(
			config,
			ctlName,
			templates,
			clearUnnecessary,
		)
		if err != nil {
			return err
		}

		return nil
	},
}
