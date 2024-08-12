package stage2_test

import (
	"testing"

	"github.com/SoenkeD/sc/src/generator/parseuml"
	"github.com/SoenkeD/sc/src/generator/stage2"
	"github.com/SoenkeD/sc/src/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStage2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stage Two Suite")
}

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
	[*] --> PrintingError : [ ! HasError(ReadConfigError) ]
	[*] -[bold]-> PrintingSuccess

	PrintingError: do / PrintError
	PrintingError -[bold]-> [*]

	PrintingSuccess: do / PrintSuccess
	PrintingSuccess -[bold]-> [*]
}
Closing -[bold]-> [*]
`

var _ = Describe("Stage Two", Ordered, func() {

	var input parseuml.ParseUmlStage1
	var out stage2.ParseUmlStage2

	It("Get stage one", func() {
		var err error
		input, err = parseuml.GenerateFromUml(TEST_UML)
		Expect(err).ToNot(HaveOccurred())

		out, err = stage2.Stage2(input)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Check UML name", func() {
		Expect(out.Name).To(Equal("TestStage1"))
	})

	It("Check guards", func() {
		Expect(out.Guards).To(ContainElements(
			"Failed",
			"HasError",
		))
	})

	It("Check actions", func() {
		Expect(out.Actions).To(ContainElements(
			"Setup",
			"ReadConfig",
			"ClearError",
			"RaiseError",
			"PrintError",
			"PrintSuccess",
		))
	})

	It("Check state existence", func() {
		expectedStates := []string{
			"/Start",
			"/End",
			"/Initialising/Start",
			"/Initialising/SettingUp",
			"/Initialising/RaisingError",
			"/Initialising/End",
			"/Closing/Start",
			"/Closing/PrintingError",
			"/Closing/PrintingSuccess",
			"/Closing/End",
		}

		var states []string
		for stateId := range out.States {
			states = append(states, stateId)
		}

		Expect(states).To(ContainElements(expectedStates))
		Expect(expectedStates).To(ContainElements(states))

	})

	It("Check entry state", func() {
		Expect(out.States).To(HaveKeyWithValue("/Start", stage2.ParsedState{
			Transitions: []stage2.ParseTransitionStage2{
				{
					Type:   types.TransitionTypeHappy,
					Target: "/Initialising/Start",
				},
			},
		}))
	})

	It("Check exit state", func() {
		Expect(out.States).To(HaveKeyWithValue("/End", stage2.ParsedState{}))
	})

	It("Check entry of state Initialising", func() {
		Expect(out.States).To(HaveKeyWithValue("/Initialising/Start", stage2.ParsedState{
			Transitions: []stage2.ParseTransitionStage2{
				{
					Type:   types.TransitionTypeHappy,
					Target: "/Initialising/SettingUp",
				},
			},
		}))
	})

})

var _ = Describe("Extract Route Information", func() {

	route := []string{
		"StartState",
		"StartState/////happy/ProcessingStart/false",
		"ProcessingStartState",
		"ProcessingStartState///CheckAlwaysTrue//normal/ProcessingDemoing/false",
		"ProcessingDemoingState",
		"AddMsgAction",
		"AddMsgAction",
		"ProcessingDemoingState/Print/Go to BurnState/CheckAlwaysTrue//normal/ProcessingBurning/false",
		"PrintAction",
		"ProcessingBurningState",
		"PrintAction",
		"PrintMsgsAction",
		"ProcessingBurningState/////happy/ProcessingEnd/false",
		"ProcessingEndState",
		"ProcessingEndState/////happy/End/false",
		"EndState",
	}

	states, trans := stage2.ExtractVisitedTransactions(route)

	It("Extract states from route", func() {
		Expect(states).To(ConsistOf(
			"Start",
			"ProcessingStart",
			"ProcessingDemoing",
			"ProcessingBurning",
			"ProcessingEnd",
			"End",
		))
	})

	It("Extract transitions from route", func() {
		Expect(trans).To(ConsistOf(
			"/Start/////happy/ProcessingStart/false",
			"/ProcessingStart///CheckAlwaysTrue//normal/ProcessingDemoing/false",
			"/ProcessingDemoing/Print/Go to BurnState/CheckAlwaysTrue//normal/ProcessingBurning/false",
			"/ProcessingBurning/////happy/ProcessingEnd/false",
			"/ProcessingEnd/////happy/End/false",
		))
	})
})

var _ = Describe("PrintTransitionType", func() {
	It("should return a transition with bold when taType is TransitionTypeHappy", func() {
		result := stage2.PrintTransitionType(types.TransitionTypeHappy, "")
		Expect(result).To(Equal("-[bold]->"))
	})

	It("should return a transition with dotted when taType is TransitionTypeError", func() {
		result := stage2.PrintTransitionType(types.TransitionTypeError, "")
		Expect(result).To(Equal("-[dotted]->"))
	})

	It("should return a transition with color when color is provided", func() {
		result := stage2.PrintTransitionType(types.TransitionTypeHappy, "red")
		Expect(result).To(Equal("-[bold,red]->"))
	})

	It("should return a transition with dotted and color when taType is TransitionTypeError and color is provided", func() {
		result := stage2.PrintTransitionType(types.TransitionTypeError, "blue")
		Expect(result).To(Equal("-[dotted,blue]->"))
	})
})
