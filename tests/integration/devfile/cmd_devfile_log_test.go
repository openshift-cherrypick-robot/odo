package devfile

import (
	"path/filepath"

	"github.com/redhat-developer/odo/tests/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("odo devfile log command tests", func() {
	var cmpName string
	var commonVar helper.CommonVar

	// This is run before every Spec (It)
	var _ = BeforeEach(func() {
		commonVar = helper.CommonBeforeEach()
		cmpName = helper.RandString(6)
	})

	// This is run after every Spec (It)
	var _ = AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})

	When("When a springboot component is created and pushed", func() {
		BeforeEach(func() {
			helper.Cmd("odo", "create", "--project", commonVar.Project, cmpName, "--context", commonVar.Context, "--devfile", helper.GetExamplePath("source", "devfiles", "springboot", "devfile.yaml")).ShouldPass()
			helper.CopyExample(filepath.Join("source", "devfiles", "springboot", "project"), commonVar.Context)
			helper.Cmd("odo", "push", "--context", commonVar.Context).ShouldPass()
		})
		It("should log run command output and fail for debug command", func() {
			output := helper.Cmd("odo", "log", "--context", commonVar.Context).ShouldPass().Out()
			Expect(output).To(ContainSubstring("ODO_COMMAND_RUN"))

			// It should fail for debug command as no debug command in devfile
			helper.Cmd("odo", "log", "--debug").ShouldFail()

			/*
			   Flaky Test odo log -f, see issue https://github.com/redhat-developer/odo/issues/3809 *** Issue got closed due to inactivity, but is not resolved ***
			   match, err := helper.RunCmdWithMatchOutputFromBuffer(30*time.Second, "program=devrun", "odo", "log", "-f")
			   Expect(err).To(BeNil())
			   Expect(match).To(BeTrue())
			*/

		})
	})

	When("a component is created but not pushed", func() {
		BeforeEach(func() {
			helper.Cmd("odo", "create", "--project", commonVar.Project, cmpName, "--context", commonVar.Context, "--devfile", helper.GetExamplePath("source", "devfiles", "nodejs", "devfile-registry.yaml")).ShouldPass()
		})
		It("Should error out", func() {
			helper.Cmd("odo", "log").ShouldFail()
		})
		When("a command is created and pushed with debugrun", func() {
			var projectDir string
			BeforeEach(func() {
				projectDir = filepath.Join(commonVar.Context, "projectDir")
				helper.CopyExample(filepath.Join("source", "web-nodejs-sample"), projectDir)
				helper.Cmd("odo", "create", "--project", commonVar.Project, cmpName, "--context", projectDir, "--devfile", helper.GetExamplePath("source", "devfiles", "nodejs", "devfile-with-debugrun.yaml")).ShouldPass()
				helper.Cmd("odo", "push", "--debug", "--context", projectDir).ShouldPass()
			})
			It("should log debug command output", func() {
				output := helper.Cmd("odo", "log", "--debug", "--context", projectDir).ShouldPass().Out()
				Expect(output).To(ContainSubstring("ODO_COMMAND_DEBUG"))

				/*
				   Flaky Test odo log -f, see issue https://github.com/redhat-developer/odo/issues/3809 *** Issue got closed due to inactivity, but is not resolved ***
				   match, err := helper.RunCmdWithMatchOutputFromBuffer(30*time.Second, "program=debugrun", "odo", "log", "-f")
				   Expect(err).To(BeNil())
				   Expect(match).To(BeTrue())
				*/
			})
		})

	})
})
