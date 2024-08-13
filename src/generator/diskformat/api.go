package diskformat

import (
	"path/filepath"

	"github.com/SoenkeD/sc/src/generator/templates"
	"github.com/SoenkeD/sc/src/types"
)

func Transform2DiskFormat(input GenerationInput, tplIn templates.GenerateTemplatesInput, cfg types.Config) (Generation, error) {
	baseDir := filepath.Join(input.RepoRoot, input.RelativeCtlRoot, input.CtlName)
	importRoot := filepath.Join(
		input.Module,
		input.RelativeCtlRoot,
		input.CtlName,
	)
	gen := Generation{}
	gen.BasePath = baseDir
	ctlDirName := "controller/"
	ctlDir := filepath.Join(baseDir, ctlDirName)
	stateDir := filepath.Join(baseDir, "state")
	actionsDir := filepath.Join(baseDir, "actions")
	guardsDir := filepath.Join(baseDir, "guards")

	gen.Dirs = append(gen.Dirs,
		baseDir,
		ctlDir,
		stateDir,
		actionsDir,
		guardsDir,
	)

	ctlFile, err := generateCtl(input.CtlName, importRoot, tplIn.TemplatedBaseFiles["initctl"])
	if err != nil {
		return Generation{}, err
	}

	baseActions, err := generateBaseActions(importRoot, tplIn.TemplatedBaseFiles["actions"], input.Actions, cfg.ForceUnitSetupRegeneration)
	if err != nil {
		return Generation{}, err
	}

	baseGuards, err := generateBaseGuards(importRoot, tplIn.TemplatedBaseFiles["guards"], input.Guards, cfg.ForceUnitSetupRegeneration)
	if err != nil {
		return Generation{}, err
	}

	stateMachine, err := generateStateMachine(input.CtlName, importRoot, tplIn.TemplatedBaseFiles["sm"], input.States, input.HasActions, input.HasGuards)
	if err != nil {
		return Generation{}, err
	}

	ctx, err := generateCtx(tplIn.TemplatedBaseFiles["context"], importRoot)
	if err != nil {
		return Generation{}, err
	}

	state, err := generateState(tplIn.TemplatedBaseFiles["state"], importRoot)
	if err != nil {
		return Generation{}, err
	}

	gen.Files = append(gen.Files,
		ctlFile,
		stateMachine,
		state,
		baseActions,
		generateBaseActionsTest(tplIn.TemplatedBaseFiles["actions_test"]),
		baseGuards,
		generateBaseGuardsTest(tplIn.TemplatedBaseFiles["guards_test"]),
		ctx,
	)

	ctlFiles, err := generateCtlFiles(ctlDirName, importRoot, tplIn, cfg.ImportPathSeparator)
	if err != nil {
		return Generation{}, err
	}
	gen.Files = append(gen.Files, ctlFiles...)

	actions, err := generateActions(importRoot, tplIn.TemplatedBaseFiles["action"], input.Actions, input.TemplatedActions)
	if err != nil {
		return Generation{}, err
	}
	gen.Files = append(gen.Files, actions...)

	actionTests, err := generateActionTests(importRoot, tplIn.TemplatedBaseFiles["action_test"], input.ActionTests, input.TemplatedActions)
	if err != nil {
		return Generation{}, err
	}
	gen.Files = append(gen.Files, actionTests...)

	guards, err := generateGuards(importRoot, tplIn.TemplatedBaseFiles["guard"], input.Guards, input.TemplatedGuards)
	if err != nil {
		return Generation{}, err
	}
	gen.Files = append(gen.Files, guards...)

	guardTests, err := generateGuardTests(importRoot, tplIn.TemplatedBaseFiles["guard_test"], input.GuardTests, input.TemplatedGuards)
	if err != nil {
		return Generation{}, err
	}
	gen.Files = append(gen.Files, guardTests...)

	gen.Files = append(gen.Files, generateStateExtensions(input.TemplatedStateExtensions)...)

	return gen, nil
}
