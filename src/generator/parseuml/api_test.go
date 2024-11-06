package parseuml_test

import (
	"github.com/SoenkeD/sc/src/generator/parseuml"
	"github.com/SoenkeD/sc/src/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const TEST_UML = `
@startuml TestStage1
[*] -[bold]-> Initialising

state Initialising {
	[*] -[bold]-> SettingUp

	SettingUp: do / Setup(cfg.yaml, prod)
	SettingUp: do / ReadConfig
	SettingUp --> RaisingError: [ Failed ] / ClearError(SetupError)
	SettingUp -[bold]-> [*]

	RaisingError: do / RaiseError
	RaisingError -[bold]-> [*]
}
Initialising -[bold]-> Closing

state Closing {
	[*] --> PrintingError : [ !HasError(ReadConfigError) ]
	[*] -[bold]-> PrintingSuccess

	PrintingError: do / PrintError
	PrintingError -[bold]-> [*]

	PrintingSuccess: do / PrintSuccess
	PrintingSuccess -[bold]-> [*]
}
Closing -[bold]-> [*]
`

var _ = Describe("Stage One", Ordered, func() {

	var out parseuml.ParseUmlStage1

	It("Parse UML", func() {
		var err error
		out, err = parseuml.GenerateFromUml(TEST_UML)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Check UML name", func() {
		Expect(out.Name).To(Equal("TestStage1"))
	})

	It("Check start state transition", func() {
		Expect(out.Transitions).To(HaveKeyWithValue("happy[*]Initialising", parseuml.ParsedTransition{
			Type:    types.TransitionTypeHappy,
			Start:   "[*]",
			Target:  "Initialising",
			Options: []string{"bold"},
		}))
	})

	It("Check end state transition", func() {
		Expect(out.Transitions).To(HaveKeyWithValue("happyInitialisingClosing", parseuml.ParsedTransition{
			Type:    types.TransitionTypeHappy,
			Start:   "Initialising",
			Target:  "Closing",
			Options: []string{"bold"},
		}))
	})

	It("Check Initialising state group", func() {
		Expect(out.StateGroups).To(HaveKeyWithValue("Initialising", "Initialising"))

		By("checking SettingUp")
		Expect(out.StateActions).To(HaveKeyWithValue("SettingUpSetupcfg.yamlprod", parseuml.ParseStateAction{
			State:        "SettingUp",
			Action:       "Setup",
			ActionParams: []string{"cfg.yaml", "prod"},
		}))

		Expect(out.StateActions).To(HaveKeyWithValue("SettingUpReadConfig", parseuml.ParseStateAction{
			State:        "SettingUp",
			Action:       "ReadConfig",
			ActionParams: nil,
		}))

		Expect(out.Transitions).To(HaveKeyWithValue("happy[*]SettingUp", parseuml.ParsedTransition{
			Type:    types.TransitionTypeHappy,
			Start:   "[*]",
			Target:  "SettingUp",
			Options: []string{"bold"},
		}))

		Expect(out.Transitions).To(HaveKeyWithValue("normalSettingUpRaisingErrorFailedClearErrorSetupError", parseuml.ParsedTransition{
			Type:         "normal",
			Start:        "SettingUp",
			Target:       "RaisingError",
			Guard:        "Failed",
			Action:       "ClearError",
			ActionParams: []string{"SetupError"},
		}))

		Expect(out.Transitions).To(HaveKeyWithValue("happySettingUp[*]", parseuml.ParsedTransition{
			Type:    types.TransitionTypeHappy,
			Start:   "SettingUp",
			Target:  "[*]",
			Options: []string{"bold"},
		}))

		By("checking RasingErr")
		Expect(out.StateActions).To(HaveKeyWithValue("RaisingErrorRaiseError", parseuml.ParseStateAction{
			State:        "RaisingError",
			Action:       "RaiseError",
			ActionParams: nil,
		}))

		Expect(out.Transitions).To(HaveKeyWithValue("happyRaisingError[*]", parseuml.ParsedTransition{
			Type:    types.TransitionTypeHappy,
			Start:   "RaisingError",
			Target:  "[*]",
			Options: []string{"bold"},
		}))
	})

	It("Check Closing state group", func() {
		Expect(out.StateGroups).To(HaveKeyWithValue("Closing", "Closing"))

		Expect(out.StateActions).To(HaveKeyWithValue("PrintingErrorPrintError", parseuml.ParseStateAction{
			State:        "PrintingError",
			Action:       "PrintError",
			ActionParams: nil,
		}))

		Expect(out.Transitions).To(HaveKeyWithValue("happyPrintingError[*]", parseuml.ParsedTransition{
			Type:    types.TransitionTypeHappy,
			Start:   "PrintingError",
			Target:  "[*]",
			Options: []string{"bold"},
		}))

		Expect(out.StateActions).To(HaveKeyWithValue("PrintingSuccessPrintSuccess", parseuml.ParseStateAction{
			State:        "PrintingSuccess",
			Action:       "PrintSuccess",
			ActionParams: nil,
		}))

		Expect(out.Transitions).To(HaveKeyWithValue("happyPrintingSuccess[*]", parseuml.ParsedTransition{
			Type:    types.TransitionTypeHappy,
			Start:   "PrintingSuccess",
			Target:  "[*]",
			Options: []string{"bold"},
		}))

		Expect(out.Transitions).To(HaveKeyWithValue("normal[*]PrintingErrorHasErrorReadConfigError", parseuml.ParsedTransition{
			Type:        "normal",
			Start:       "[*]",
			Target:      "PrintingError",
			Guard:       "HasError",
			GuardParams: []string{"ReadConfigError"},
			Action:      "",
			Negation:    true,
		}))

		Expect(out.Transitions).To(HaveKeyWithValue("happy[*]PrintingSuccess", parseuml.ParsedTransition{
			Type:    types.TransitionTypeHappy,
			Start:   "[*]",
			Target:  "PrintingSuccess",
			Guard:   "",
			Action:  "",
			Options: []string{"bold"},
		}))
	})

})
