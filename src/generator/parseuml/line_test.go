package parseuml_test

import (
	"github.com/SoenkeD/sc/src/generator/parseuml"
	"github.com/SoenkeD/sc/src/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseLine()", func() {

	var uml *parseuml.ParseUmlStage1

	BeforeEach(func() {
		uml = &parseuml.ParseUmlStage1{}
		uml.Init()
	})

	It("Filter empty line", func() {
		area, err := parseuml.ParseLine("", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		expectEmptyUml(uml)
	})

	It("Filter empty space line", func() {
		area, err := parseuml.ParseLine("    ", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		expectEmptyUml(uml)
	})

	It("Filter empty comment line", func() {
		area, err := parseuml.ParseLine("'Some Comment", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		expectEmptyUml(uml)
	})

	It("Filter empty space comment line", func() {
		area, err := parseuml.ParseLine(" 'Some Comment", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		expectEmptyUml(uml)
	})

	It("Headline", func() {
		area, err := parseuml.ParseLine(parseuml.UmlPrefix+" UmlName", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		Expect(uml.Name).To(Equal("UmlName"))
		Expect(uml.LinesInOrder).To(HaveLen(0))
	})

	It("Headline with comment", func() {
		area, err := parseuml.ParseLine(parseuml.UmlPrefix+" UmlName 'UmlName", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		Expect(uml.Name).To(Equal("UmlName"))
		Expect(uml.LinesInOrder).To(HaveLen(0))
	})

	It("Headline without name", func() {
		area, err := parseuml.ParseLine(parseuml.UmlPrefix+"MissingSpace", uml)
		Expect(err).To(HaveOccurred())
		Expect(area).To(Equal("head"))
		Expect(uml.LinesInOrder).To(HaveLen(0))
	})

	It("Uml Suffix", func() {
		_, err := parseuml.ParseLine(parseuml.UmlSuffix, uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(uml.LinesInOrder).To(HaveLen(0))
	})

	It("State group closing", func() {
		area, err := parseuml.ParseLine("}", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		Expect(uml.LinesInOrder).To(HaveLen(1))
		Expect(uml.LinesInOrder).To(ContainElements(parseuml.LINE_STATE_GROUP_CLOSING))
	})

	It("State group opening", func() {
		area, err := parseuml.ParseLine("state Test {", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		Expect(uml.LinesInOrder).To(HaveLen(1))
		Expect(uml.LinesInOrder).To(ContainElements(parseuml.LINE_STATE_GROUP_OPENING + "Test"))
	})

	It("Fail parsing content line", func() {
		area, err := parseuml.ParseLine("::", uml)
		Expect(err).To(HaveOccurred())
		Expect(area).To(Equal("content_line"))
		Expect(uml.LinesInOrder).To(HaveLen(0))
	})

	It("State action", func() {
		area, err := parseuml.ParseLine("State: "+parseuml.DO_SLASH_LINE+"Action", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		Expect(uml.LinesInOrder).To(HaveLen(1))
		Expect(uml.LinesInOrder).To(ContainElements(parseuml.LINE_ACTION + "State" + "Action"))
	})

	It("Transition without content part", func() {
		area, err := parseuml.ParseLine("Start --> End", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		Expect(uml.LinesInOrder).To(HaveLen(1))
		lineStr := parseuml.LINE_TRANSITION + string(types.TransitionTypeNormal) + "Start" + "End"
		Expect(uml.LinesInOrder).To(ContainElements(lineStr))
	})

	It("Transition with content part", func() {
		area, err := parseuml.ParseLine("Start --> End: [ Has ] / Do", uml)
		Expect(err).ToNot(HaveOccurred())
		Expect(area).To(BeEmpty())
		Expect(uml.LinesInOrder).To(HaveLen(1))
		lineStr := parseuml.LINE_TRANSITION + string(types.TransitionTypeNormal) + "Start" + "End" + "Has" + "Do"
		Expect(uml.LinesInOrder[0]).To(Equal(lineStr))
		Expect(uml.LinesInOrder).To(ContainElements(lineStr))
	})
})

func expectEmptyUml(uml *parseuml.ParseUmlStage1) {
	Expect(uml.Name).To(BeEmpty())
	Expect(uml.LinesInOrder).To(HaveLen(0))
	Expect(uml.StateActions).To(HaveLen(0))
	Expect(uml.StateGroups).To(HaveLen(0))
	Expect(uml.Transitions).To(HaveLen(0))
}
