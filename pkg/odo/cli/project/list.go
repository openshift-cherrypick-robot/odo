package project

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/redhat-developer/odo/pkg/kclient"
	"github.com/redhat-developer/odo/pkg/log"
	"github.com/redhat-developer/odo/pkg/machineoutput"
	"github.com/redhat-developer/odo/pkg/odo/cmdline"
	"github.com/redhat-developer/odo/pkg/odo/genericclioptions"
	"github.com/redhat-developer/odo/pkg/project"
	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const listRecommendedCommandName = "list"

var (
	listExample = ktemplates.Examples(`
	# List all the projects
    %[1]s`)
	listLongDesc = ktemplates.LongDesc(`
	List all the projects
`)
)

// ProjectListOptions encapsulates the options for the odo project list command
type ProjectListOptions struct {
	// Context
	*genericclioptions.Context

	// Clients
	prjClient project.Client
}

// NewProjectListOptions creates a new ProjectListOptions instance
func NewProjectListOptions(prjClient project.Client) *ProjectListOptions {
	return &ProjectListOptions{
		prjClient: prjClient,
	}
}

// Complete completes ProjectListOptions after they've been created
func (plo *ProjectListOptions) Complete(cmdline cmdline.Cmdline, args []string) (err error) {
	plo.Context, err = genericclioptions.New(genericclioptions.NewCreateParameters(cmdline))
	return err
}

// Validate validates the ProjectListOptions based on completed values
func (plo *ProjectListOptions) Validate() (err error) {
	return nil
}

// Run contains the logic for the odo project list command
func (plo *ProjectListOptions) Run() error {
	projects, err := plo.prjClient.List()
	if err != nil {
		return err
	}

	if log.IsJSON() {
		machineoutput.OutputSuccess(projects)
	} else {
		err = HumanReadableOutput(os.Stdout, projects)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewCmdProjectList implements the odo project list command.
func NewCmdProjectList(name, fullName string) *cobra.Command {
	// The error is not handled at this point, it will be handled during Context creation
	kubclient, _ := kclient.New()
	o := NewProjectListOptions(project.NewClient(kubclient))
	projectListCmd := &cobra.Command{
		Use:         name,
		Short:       listLongDesc,
		Long:        listLongDesc,
		Example:     fmt.Sprintf(listExample, fullName),
		Args:        cobra.ExactArgs(0),
		Annotations: map[string]string{"machineoutput": "json"},
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}
	return projectListCmd
}

// HumanReadableOutput outputs the list of projects in a human readable format
func HumanReadableOutput(w io.Writer, o project.ProjectList) error {
	if len(o.Items) == 0 {
		return fmt.Errorf("you are not a member of any projects. You can request a project to be created using the `odo project create <project_name>` command")
	}
	wr := tabwriter.NewWriter(w, 5, 2, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(wr, "ACTIVE", "\t", "NAME")
	for _, project := range o.Items {
		activeMark := " "
		if project.Status.Active {
			activeMark = "*"
		}
		fmt.Fprintln(wr, activeMark, "\t", project.Name)
	}
	wr.Flush()
	return nil
}
