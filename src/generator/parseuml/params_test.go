package parseuml_test

import (
	"github.com/SoenkeD/sc/src/generator/parseuml"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseParams()", func() {
	It("No params", func() {
		name, params, err := parseuml.ParseParams("Name")
		Expect(err).ToNot(HaveOccurred())
		Expect(name).To(Equal("Name"))
		Expect(params).To(HaveLen(0))
	})

	It("No params, trim spaces", func() {
		name, params, err := parseuml.ParseParams("  Name ")
		Expect(err).ToNot(HaveOccurred())
		Expect(name).To(Equal("Name"))
		Expect(params).To(HaveLen(0))
	})

	It("One params", func() {
		name, params, err := parseuml.ParseParams("Name(Param)")
		Expect(err).ToNot(HaveOccurred())
		Expect(name).To(Equal("Name"))
		Expect(params).To(HaveLen(1))
		Expect(params).To(ContainElements("Param"))
	})

	It("Multiple params", func() {
		name, params, err := parseuml.ParseParams("Name(Param1, Param2)")
		Expect(err).ToNot(HaveOccurred())
		Expect(name).To(Equal("Name"))
		Expect(params).To(HaveLen(2))
		Expect(params).To(ContainElements("Param1", "Param2"))
	})

	It("To many (", func() {
		_, _, err := parseuml.ParseParams("Name(())")
		Expect(err).To(HaveOccurred())
	})
})
