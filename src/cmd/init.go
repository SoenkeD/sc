package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	initCmd.Flags().StringVarP(&config.CtlDir, "ctl", "", "", "relative path of the controller directory")
	viper.BindPFlag("ctl", rootCmd.PersistentFlags().Lookup("ctl"))

	initCmd.Flags().StringVarP(&initCtl, "name", "n", "", "name of the controller to create e.g. demo")
	initCmd.Flags().StringVarP(&setupSourceRepo, "setup", "s", "", "url of the setup repository")
	initCmd.Flags().StringVarP(&containerDriver, "container", "d", "docker", "container provider docker | podman. Defaults to docker")

	rootCmd.AddCommand(initCmd)

	cmdsWithoutConfig = append(cmdsWithoutConfig, CmdNameInit)
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

const CmdNameInit = "init"

var initCmd = &cobra.Command{
	Use:   CmdNameInit,
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

		exists, err := utils.FileOrDirExists(config.RepoRoot)
		if err != nil {
			return err
		}

		if !exists {
			err = utils.CreateDirs(config.RepoRoot)
			if err != nil {
				return err
			}
		} else {
			log.Println("skipped creation of root directory")
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

		initCfg.Module = config.Module
		if config.CtlDir != "" {
			initCfg.CtlDir = config.CtlDir
		}

		outYml, err := yaml.Marshal(initCfg)
		if err != nil {
			return err
		}

		scFile := filepath.Join(scPath, "sc.yaml")
		err = utils.WriteFile(scFile, string(outYml))
		if err != nil {
			return err
		}

		// set merged config to be used
		initCfg.RepoRoot = config.RepoRoot // the repo root should not be stored in sc.yaml
		config = initCfg

		ctlPath := filepath.Join(config.RepoRoot, config.CtlDir)
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
		umlContent, err := templates.ExecTemplate(umlFile, string(umlBytes), templateInputInitFiles{
			Cfg:       config,
			InitCtl:   initCtl,
			Container: containerDriver,
		}, nil)
		if err != nil {
			return err
		}

		err = utils.WriteFile(umlFile, umlContent)
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

	confirm, skipped := utils.UserConfirm("Execute '" + execCmd + "'? (y/s/n) ")

	if skipped {
		log.Println("skipped action")
		return nil
	}

	if !confirm {
		return fmt.Errorf("user stopped the init script")
	}

	parts := strings.Fields(execCmd)

	cmdOut, err := utils.ExecuteCommand(parts[0], parts[1:], nil, config.RepoRoot)
	if err != nil {
		log.Printf("Executing the command=%s failed with=%s and output %s", execCmd, err, cmdOut)
		return err
	}

	return nil
}
