package cmd

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/SoenkeD/sc/src/types"
)

var cfgFile string
var config types.Config

func readCfg() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigType("yaml")
		viper.SetConfigName("sc")
		viper.AddConfigPath("./sc")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s \n", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Error unmarshalling config: %s \n", err)
		os.Exit(1)
	}
}

func validateCfg() {
	validate := validator.New()
	valErr := validate.Struct(config)
	if valErr != nil {
		log.Println("failed to validate config", valErr)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to the config file")
}

var rootCmd = &cobra.Command{
	Use:   "sc",
	Short: "SC is a state chart code generator",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "init" && cmd.Name() != "version" {
			readCfg()
			validateCfg()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
