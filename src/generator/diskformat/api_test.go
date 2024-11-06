package diskformat_test

import (
	"github.com/SoenkeD/sc/src/types"

	"github.com/SoenkeD/sc/src/generator/diskformat"
	"github.com/SoenkeD/sc/src/generator/templates"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transform2DiskFormat", Ordered, func() {

	var gen diskformat.Generation

	BeforeAll(func() {

		input := diskformat.GenerationInput{
			CtlName:         "test",
			RepoRoot:        "/code/test",
			RelativeCtlRoot: "src",
			Module:          "github.com/User/Test",
			Actions: map[string]string{
				"DemoAction":         "some_code",
				"DemoTemplateAction": "some_template_code",
			},
			ActionTests: map[string]string{
				"DemoAction":         "some_test_code",
				"DemoTemplateAction": "some_default_test_code",
			},
			TemplatedActions: map[string]string{
				"DemoTemplateAction": "some_template_code",
			},
			Guards: map[string]string{
				"DemoGuard":         "some_code",
				"DemoTemplateGuard": "some_template_code",
			},
			GuardTests: map[string]string{
				"DemoGuard":         "some_test_code",
				"DemoTemplateGuard": "some_default_test_code",
			},
			TemplatedGuards: map[string]string{
				"DemoTemplateGuard": "some_template_code",
			},
			States: []string{
				"demo_state_code",
			},
		}

		cfg := types.Config{
			ImportPathSeparator: "/",
		}

		var err error
		gen, err = diskformat.Transform2DiskFormat(input, templates.GenerateTemplatesInput{}, cfg)
		Expect(err).ToNot(HaveOccurred())

	})

	It("checks base path", func() {
		Expect(gen.BasePath).To(Equal("/code/test/src/test"))
	})

	It("check dirs", func() {
		Expect(gen.Dirs).To(ConsistOf(
			"/code/test/src/test",
			"/code/test/src/test/controller",
			"/code/test/src/test/state",
			"/code/test/src/test/actions",
			"/code/test/src/test/guards",
		))
	})

	It("checks files", func() {
		Expect(gen.Files).To(ConsistOf(

			SatisfyAll(
				HaveField("Type", Equal("ctl")), // TODO type multiple times not matching
				HaveField("Name", Equal("initctl")),
				HaveField("Path", Equal("")),
			),
			SatisfyAll(
				HaveField("Type", Equal("ctl")),
				HaveField("Name", Equal("sm")),
				HaveField("Path", Equal("")),
			),
			SatisfyAll(
				HaveField("Type", Equal("state")),
				HaveField("Name", Equal("ExtendedState")),
				HaveField("Path", Equal("state/")),
			),
			SatisfyAll(
				HaveField("Type", Equal("context")),
				HaveField("Name", Equal("Ctx")),
				HaveField("Path", Equal("state/")),
			),
			SatisfyAll(
				HaveField("Type", Equal("ctl")),
				HaveField("Name", Equal("ctl")),
				HaveField("Path", Equal("controller/")),
			),
			SatisfyAll(
				HaveField("Type", Equal("actions")),
				HaveField("Name", Equal("actions")),
				HaveField("Path", Equal("actions/")),
			),
			SatisfyAll(
				HaveField("Type", Equal("actions_test")),
				HaveField("Name", Equal("actions_test")),
				HaveField("Path", Equal("actions/")),
			),
			SatisfyAll(
				HaveField("Type", Equal("guards")),
				HaveField("Name", Equal("guards")),
				HaveField("Path", Equal("guards/")),
			),
			SatisfyAll(
				HaveField("Type", Equal("guards_test")),
				HaveField("Name", Equal("guards_test")),
				HaveField("Path", Equal("guards/")),
			),

			// actions
			SatisfyAll(
				HaveField("Type", Equal("action")),
				HaveField("Name", Equal("DemoAction")),
				HaveField("Path", Equal("actions/")),
			),
			SatisfyAll(
				HaveField("Type", Equal("action_test")),
				HaveField("Name", Equal("DemoAction_test")),
				HaveField("Path", Equal("actions/")),
			),
			SatisfyAll(
				HaveField("Type", Equal("action")),
				HaveField("Name", Equal("DemoTemplateAction")),
				HaveField("Path", Equal("actions/")),
			),
			// no template test
			// guards
			SatisfyAll(
				HaveField("Type", Equal("guard")),
				HaveField("Name", Equal("DemoGuard")),
				HaveField("Path", Equal("guards/")),
			),
			SatisfyAll(
				HaveField("Type", Equal("guard_test")),
				HaveField("Name", Equal("DemoGuard_test")),
				HaveField("Path", Equal("guards/")),
			),
			SatisfyAll(
				HaveField("Type", Equal("guard")),
				HaveField("Name", Equal("DemoTemplateGuard")),
				HaveField("Path", Equal("guards/")),
			),
			// no template test

		))
	})
})
