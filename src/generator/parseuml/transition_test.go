package parseuml_test

import (
	"github.com/SoenkeD/sc/src/generator/parseuml"
	"github.com/SoenkeD/sc/src/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseTransitionType()", func() {
	It("Normal transition", func() {
		ta, err := parseuml.ParseTransitionType(parseuml.TransitionNormal)
		Expect(err).ToNot(HaveOccurred())
		Expect(ta).To(Equal(types.TransitionTypeNormal))
	})

	It("Happy transition", func() {
		ta, err := parseuml.ParseTransitionType(parseuml.TransitionHappy)
		Expect(err).ToNot(HaveOccurred())
		Expect(ta).To(Equal(types.TransitionTypeHappy))
	})

	It("Error transition", func() {
		ta, err := parseuml.ParseTransitionType(parseuml.TransitionError)
		Expect(err).ToNot(HaveOccurred())
		Expect(ta).To(Equal(types.TransitionTypeError))
	})

	It("Unkown transition", func() {
		_, err := parseuml.ParseTransitionType("-" + parseuml.TransitionNormal)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("ParseTransition()", func() {

	It("To few tokens given", func() {
		tokens := []string{"", ""}
		_, err := parseuml.ParseTransition(tokens, "")
		Expect(err).To(HaveOccurred())
	})

	It("To many tokens given", func() {
		tokens := []string{"", "", "", ""}
		_, err := parseuml.ParseTransition(tokens, "")
		Expect(err).To(HaveOccurred())
	})

	It("No content line given", func() {
		tokens := []string{"Start", parseuml.TransitionNormal, "End"}
		ta, err := parseuml.ParseTransition(tokens, "")
		Expect(err).ToNot(HaveOccurred())

		Expect(ta.Start).To(Equal("Start"))
		Expect(ta.Type).To(Equal(types.TransitionTypeNormal))
		Expect(ta.Target).To(Equal("End"))

		Expect(ta.Action).To(Equal(""))
		Expect(ta.ActionParams).To(HaveLen(0))

		Expect(ta.Negation).To(BeFalse())
		Expect(ta.Guard).To(Equal(""))
		Expect(ta.GuardParams).To(HaveLen(0))
	})

	It("Content line contains guard only", func() {
		tokens := []string{"Start", parseuml.TransitionNormal, "End"}
		ta, err := parseuml.ParseTransition(tokens, " [ SomeGuard(Param) ] ")
		Expect(err).ToNot(HaveOccurred())

		Expect(ta.Start).To(Equal("Start"))
		Expect(ta.Type).To(Equal(types.TransitionTypeNormal))
		Expect(ta.Target).To(Equal("End"))

		Expect(ta.Action).To(Equal(""))
		Expect(ta.ActionParams).To(HaveLen(0))

		Expect(ta.Negation).To(BeFalse())
		Expect(ta.Guard).To(Equal("SomeGuard"))
		Expect(ta.GuardParams).To(HaveLen(1))
		Expect(ta.GuardParams).To(ContainElements("Param"))
	})

	It("Content line contains action only", func() {
		tokens := []string{"Start", parseuml.TransitionNormal, "End"}
		ta, err := parseuml.ParseTransition(tokens, " / SomeAction(Param) ")
		Expect(err).ToNot(HaveOccurred())

		Expect(ta.Start).To(Equal("Start"))
		Expect(ta.Type).To(Equal(types.TransitionTypeNormal))
		Expect(ta.Target).To(Equal("End"))

		Expect(ta.Action).To(Equal("SomeAction"))
		Expect(ta.ActionParams).To(HaveLen(1))
		Expect(ta.ActionParams).To(ContainElements("Param"))

		Expect(ta.Negation).To(BeFalse())
		Expect(ta.Guard).To(Equal(""))
		Expect(ta.GuardParams).To(HaveLen(0))
	})

	It("Content line contains boith action & guard", func() {
		tokens := []string{"Start", parseuml.TransitionNormal, "End"}
		ta, err := parseuml.ParseTransition(tokens, " [ SomeGuard(Param) ] / SomeAction(Param) ")
		Expect(err).ToNot(HaveOccurred())

		Expect(ta.Start).To(Equal("Start"))
		Expect(ta.Type).To(Equal(types.TransitionTypeNormal))
		Expect(ta.Target).To(Equal("End"))

		Expect(ta.Action).To(Equal("SomeAction"))
		Expect(ta.ActionParams).To(HaveLen(1))
		Expect(ta.ActionParams).To(ContainElements("Param"))

		Expect(ta.Negation).To(BeFalse())
		Expect(ta.Guard).To(Equal("SomeGuard"))
		Expect(ta.GuardParams).To(HaveLen(1))
		Expect(ta.GuardParams).To(ContainElements("Param"))
	})
})
