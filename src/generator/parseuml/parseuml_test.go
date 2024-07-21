package parseuml_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStage1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stage One Unit Test Suite")
}
