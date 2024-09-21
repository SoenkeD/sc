package generator

import (
	"log"
	"os"
	"path/filepath"

	"github.com/SoenkeD/sc/src/generator/aftercompilation"
	"github.com/SoenkeD/sc/src/generator/diskformat"
	"github.com/SoenkeD/sc/src/generator/parseuml"
	"github.com/SoenkeD/sc/src/generator/stage2"
	"github.com/SoenkeD/sc/src/generator/templates"
	"github.com/SoenkeD/sc/src/types"
	"github.com/SoenkeD/sc/src/utils"
)

func ParseUmlFile(umlPath string) (st2 stage2.ParseUmlStage2, err error) {
	umlBytes, err := os.ReadFile(umlPath)
	if err != nil {
		return
	}
	umlStr := string(umlBytes)

	st1, err := parseuml.GenerateFromUml(umlStr)
	if err != nil {
		return
	}

	st2, err = stage2.Stage2(st1)
	if err != nil {
		return
	}

	err = stage2.CheckForHappyPath(st2.States)
	if err != nil {
		return
	}

	return
}

func Generate(cfg types.Config, ctlName string, tplIn templates.GenerateTemplatesInput, clearUnnecessary bool) error {

	umlFilePath := filepath.Join(cfg.RepoRoot, cfg.CtlDir, ctlName, ctlName+".plantuml")
	st2, err := ParseUmlFile(umlFilePath)
	if err != nil {
		return err
	}

	actionsDir := filepath.Join(cfg.RepoRoot, cfg.CtlDir, ctlName, "actions")
	guardsDir := filepath.Join(cfg.RepoRoot, cfg.CtlDir, ctlName, "guards")
	templateCollection, err := templates.ExecuteTemplates(
		st2,
		tplIn,
		cfg.RepoRoot,
		cfg.Module,
		cfg.CtlDir,
		ctlName,
		actionsDir,
		guardsDir,
		cfg.Language,
	)
	if err != nil {
		return err
	}

	input := diskformat.GenerationInput{
		CtlName:                  ctlName,
		RepoRoot:                 cfg.RepoRoot,
		Module:                   cfg.Module,
		Actions:                  templateCollection.Actions,
		TemplatedActions:         templateCollection.ActionTemplates,
		ActionTests:              templateCollection.ActionTests,
		Guards:                   templateCollection.Guards,
		TemplatedGuards:          templateCollection.GuardTemplates,
		GuardTests:               templateCollection.GuardTests,
		TemplatedStateExtensions: tplIn.TemplatedStateExtensions,
		States:                   templateCollection.States,
		RelativeCtlRoot:          cfg.CtlDir,
		HasActions:               len(st2.Actions) > 0,
		HasGuards:                len(st2.Guards) > 0,
	}

	gen, err := diskformat.Transform2DiskFormat(input, tplIn, cfg)
	if err != nil {
		return err
	}

	unnecessary, err := aftercompilation.WarnUnnecessaryFiles(actionsDir, st2.Actions, guardsDir, st2.Guards, cfg.Language)
	if err != nil {
		return err
	}

	if clearUnnecessary {
		log.Println("clearing no longer need files")
		err = utils.RemoveFiles(unnecessary)
		if err != nil {
			return err
		}
	}

	err = aftercompilation.WriteToDisk(gen, cfg.Language, cfg.EnableGeneratedFilePrefix, cfg.EnableFileCapitalization)
	if err != nil {
		return err
	}

	return nil
}
