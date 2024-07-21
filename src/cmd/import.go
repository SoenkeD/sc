package cmd

import (
	"github.com/SoenkeD/sc/src/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	importCmd.Flags().StringVarP(&config.RepoRoot, "root", "r", "", "base path to generate ctl in")
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))

	importCmd.Flags().StringVarP(&config.Module, "module", "m", "", "name of the module e.g. github.com/SoenkeD/sc")
	viper.BindPFlag("module", rootCmd.PersistentFlags().Lookup("module"))

	importCmd.Flags().StringVarP(&ctlName, "name", "n", "", "name of the ctl")

	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import templates from sources",
	RunE: func(cmd *cobra.Command, args []string) error {

		for _, importItem := range config.Imports {

			err := utils.DownloadFolder(
				importItem.RepoOwner,
				importItem.RepoName,
				importItem.RepoPath,
				importItem.LocalPath,
				importItem.Token,
				importItem.RepoPath,
			)

			if err != nil {
				return err
			}
		}

		return nil
	},
}
