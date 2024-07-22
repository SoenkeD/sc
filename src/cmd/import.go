package cmd

import (
	"github.com/SoenkeD/sc/src/types"
	"github.com/SoenkeD/sc/src/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	importCmd.Flags().StringVarP(&config.RepoRoot, "root", "r", "", "base path to generate ctl in")
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))

	importCmd.Flags().StringVarP(&config.Module, "module", "m", "", "name of the module e.g. github.com/SoenkeD/sc")
	viper.BindPFlag("module", rootCmd.PersistentFlags().Lookup("module"))

	rootCmd.AddCommand(importCmd)
}

func importAll(imports []types.Import, localPathPrefix string) error {

	for _, importItem := range imports {

		err := utils.DownloadFolder(
			importItem.RepoOwner,
			importItem.RepoName,
			importItem.RepoPath,
			importItem.LocalPath,
			importItem.Token,
			importItem.RepoPath,
			localPathPrefix,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import templates from sources",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := importAll(config.Imports, "")
		if err != nil {
			return err
		}

		return nil
	},
}
