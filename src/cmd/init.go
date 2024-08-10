package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/SoenkeD/sc/src/generator/templates"
	"github.com/SoenkeD/sc/src/types"
	"github.com/SoenkeD/sc/src/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func init() {
	initCmd.Flags().StringVarP(&config.RepoRoot, "root", "r", "", "base path to generate ctl in")
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))

	initCmd.Flags().StringVarP(&config.Module, "module", "m", "", "name of the module e.g. github.com/SoenkeD/sc")
	viper.BindPFlag("module", rootCmd.PersistentFlags().Lookup("module"))

	initCmd.Flags().StringVarP(&initCtl, "name", "n", "", "name of the controller to create e.g. demo")
	initCmd.Flags().StringVarP(&setupSourceRepo, "setup", "s", "", "url of the setup repository")
	initCmd.Flags().StringVarP(&containerDriver, "container", "d", "docker", "container provider docker | podman. Defaults to docker")

	rootCmd.AddCommand(initCmd)
}

var (
	initCtl         string
	setupSourceRepo string
	containerDriver string
)

type projectSetupInput struct {
	ScYaml                 string       `yaml:"scYaml"`
	DefaultPlantUml        string       `yaml:"defaultPlantUml"`
	AfterStructureCreation []string     `yaml:"afterStructureCreation"`
	AfterAll               []string     `yaml:"afterAll"`
	Files                  []fileConfig `yaml:"files"`
}

type fileConfig struct {
	Src    string `yaml:"src"`
	Target string `yaml:"target"`
}

type templateInputInitFiles struct {
	Cfg       types.Config
	InitCtl   string
	Container string
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init an empty project",
	RunE: func(cmd *cobra.Command, args []string) error {

		tempDir, err := os.MkdirTemp("", "sc-installs-")
		if err != nil {
			log.Fatalf("Failed to create temp directory: %v", err)
		}

		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				log.Printf("Failed to remove install script temp directory: %v", err)
			}
		}()

		owner, repo, fPath, err := utils.ExtractGitHubInfo(setupSourceRepo)
		if err != nil {
			return err
		}

		err = utils.DownloadFolder(owner, repo, fPath, tempDir, "", fPath, "")
		if err != nil {
			return err
		}

		setupFBytes, err := os.ReadFile(filepath.Join(tempDir, "setup.yaml"))
		if err != nil {
			return err
		}

		var input projectSetupInput
		err = yaml.Unmarshal(setupFBytes, &input)
		if err != nil {
			return err
		}

		err = utils.CreateDirs(config.RepoRoot)
		if err != nil {
			return err
		}

		scPath := filepath.Join(config.RepoRoot, "sc")
		err = utils.CreateDirs(scPath)
		if err != nil {
			return err
		}

		scYamlBytes, err := os.ReadFile(filepath.Join(tempDir, input.ScYaml))
		if err != nil {
			return err
		}

		var initCfg types.Config
		err = yaml.Unmarshal(scYamlBytes, &initCfg)
		if err != nil {
			return err
		}

		initCfg.RepoRoot = config.RepoRoot
		initCfg.Module = config.Module
		config = initCfg

		outYml, err := yaml.Marshal(initCfg)
		if err != nil {
			return err
		}

		scFile := filepath.Join(scPath, "sc.yaml")
		err = utils.WriteFile(scFile, string(outYml))
		if err != nil {
			return err
		}

		ctlPath := filepath.Join(config.RepoRoot, "src/controller")
		initCtlPath := filepath.Join(ctlPath, initCtl)
		err = utils.CreateDirs(initCtlPath)
		if err != nil {
			return err
		}

		umlBytes, err := os.ReadFile(filepath.Join(tempDir, input.DefaultPlantUml))
		if err != nil {
			return err
		}

		umlFile := filepath.Join(initCtlPath, initCtl+".plantuml")
		err = utils.WriteFile(umlFile, string(umlBytes))
		if err != nil {
			return err
		}

		for _, cmd := range input.AfterStructureCreation {
			err = execInitCmd(cmd)
			if err != nil {
				return err
			}
		}

		writeFiles := map[string]string{}
		for _, cFile := range input.Files {

			content, err := os.ReadFile(filepath.Join(tempDir, cFile.Src))
			if err != nil {
				return err
			}

			writeContent, err := templates.ExecTemplate(filepath.Join(config.RepoRoot, cFile.Src), string(content), templateInputInitFiles{
				Cfg:       config,
				InitCtl:   initCtl,
				Container: containerDriver,
			}, nil)
			if err != nil {
				return err
			}

			targetPath := filepath.Join(config.RepoRoot, cFile.Target)
			writeFiles[targetPath] = writeContent
		}

		for cFilePath, cFile := range writeFiles {
			err = utils.CreateDirs(filepath.Dir(cFilePath))
			if err != nil {
				return err
			}

			err = utils.WriteFile(cFilePath, cFile)
			if err != nil {
				return err
			}
		}

		err = importAll(config.Imports, config.RepoRoot)
		if err != nil {
			return err
		}

		for _, cmd := range input.AfterAll {
			err = execInitCmd(cmd)
			if err != nil {
				return err
			}
		}

		log.Println("Successfully initialized the project")
		log.Println("To get started edit at least the Print action to print the first input argument and run 'make run'")

		return nil
	},
}

func execInitCmd(cmd string) error {
	execCmd, err := templates.ExecTemplate(cmd, cmd, templateInputInitFiles{
		Cfg:       config,
		InitCtl:   initCtl,
		Container: containerDriver,
	}, nil)
	if err != nil {
		return err
	}

	if !utils.UserConfirm("Execute '" + execCmd + "'? (y/n) ") {
		return fmt.Errorf("user stopped the init script")
	}

	_, err = utils.ExecuteCommand(execCmd, config.RepoRoot)
	if err != nil {
		return err
	}

	return nil
}