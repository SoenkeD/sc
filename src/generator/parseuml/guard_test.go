package parseuml_test

import (
	"github.com/SoenkeD/sc/src/generator/parseuml"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseGuard()", func() {

	It("Empty guard part given", func() {
		guard, guardParams, negation, err := parseuml.ParseGuard("  ")
		Expect(err).ToNot(HaveOccurred())
		Expect(guard).To(BeEmpty())
		Expect(guardParams).To(HaveLen(0))
		Expect(negation).To(BeFalse())
	})

	It("Not empty guard part but without []", func() {
		_, _, _, err := parseuml.ParseGuard(" SomeFakeGuard ")
		Expect(err).To(HaveOccurred())
	})

	It("Guard without negation", func() {
		guard, guardParams, negation, err := parseuml.ParseGuard(" [ SomeGuard(Param1, Param2) ] ")
		Expect(err).ToNot(HaveOccurred())
		Expect(guard).To(Equal("SomeGuard"))
		Expect(guardParams).To(HaveLen(2))
		Expect(guardParams).To(ContainElements("Param1", "Param2"))
		Expect(negation).To(BeFalse())
	})

	It("Guard with negation", func() {
		guard, guardParams, negation, err := parseuml.ParseGuard(" [ ! SomeGuard(Param1, Param2) ] ")
		Expect(err).ToNot(HaveOccurred())
		Expect(guard).To(Equal("SomeGuard"))
		Expect(guardParams).To(HaveLen(2))
		Expect(guardParams).To(ContainElements("Param1", "Param2"))
		Expect(negation).To(BeTrue())
	})
})
