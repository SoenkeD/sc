package diskformat

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStage1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Generator DiskFormat Unit Test Suite")
}
