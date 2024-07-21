package parseuml_test

import (
	"strings"

	"github.com/SoenkeD/sc/src/generator/parseuml"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseStateActionFromCode()", func() {
	It("No params", func() {
		name, params, err := parseuml.ParseStateActionFromCode(parseuml.DO_SLASH_LINE + "Name")
		Expect(err).ToNot(HaveOccurred())
		Expect(name).To(Equal("Name"))
		Expect(params).To(HaveLen(0))
	})

	It("With params", func() {
		name, params, err := parseuml.ParseStateActionFromCode(parseuml.DO_SLASH_LINE + "Name(Param1, Param2)")
		Expect(err).ToNot(HaveOccurred())
		Expect(name).To(Equal("Name"))
		Expect(params).To(HaveLen(2))
		Expect(params).To(ContainElements("Param1", "Param2"))
	})

	It("Missing prefix", func() {
		_, _, err := parseuml.ParseStateActionFromCode("Name")
		Expect(err).To(HaveOccurred())
	})

	It("Multiple do slash", func() {
		_, _, err := parseuml.ParseStateActionFromCode(strings.Repeat(parseuml.DO_SLASH_LINE, 2) + "Name")
		Expect(err).To(HaveOccurred())
	})
})
