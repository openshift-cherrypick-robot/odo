package component

import (
	"fmt"

	servicebinding "github.com/redhat-developer/service-binding-operator/apis/binding/v1alpha1"
	"github.com/spf13/cobra"

	"github.com/redhat-developer/odo/pkg/devfile/location"
	"github.com/redhat-developer/odo/pkg/odo/cmdline"
	"github.com/redhat-developer/odo/pkg/odo/genericclioptions"
	odoutil "github.com/redhat-developer/odo/pkg/odo/util"

	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

// LinkRecommendedCommandName is the recommended link command name
const LinkRecommendedCommandName = "link"

var (
	linkExample = ktemplates.Examples(`# Link the current component to the 'EtcdCluster' named 'myetcd'
%[1]s EtcdCluster/myetcd

# Link current component to the 'backend' component (backend must have a single exposed port)
%[1]s backend

# Link current component to the 'backend' component and puts the link definition in the devfile instead of a separate file
%[1]s backend --inlined

# Link component 'nodejs' to the 'backend' component
%[1]s backend --component nodejs

# Link current component to port 8080 of the 'backend' component (backend must have port 8080 exposed) 
%[1]s backend --port 8080

# Link the current component to the 'EtcdCluster' named 'myetcd'
# and make the secrets accessible as files in the '/bindings/etcd/' directory
%[1]s EtcdCluster/myetcd  --bind-as-files --name etcd`)

	linkLongDesc = `Link current or provided component to a service (backed by an Operator) or another component

The appropriate secret will be added to the environment of the source component as environment variables by 
default.

For example:

Let us say we have created a nodejs application called 'frontend' which we link to an another component called
'backend' which exposes port 8080, then linking the 2 using:
odo link backend --component frontend

The frontend has 2 ENV variables it can use:
SERVICE_BACKEND_IP=10.217.4.194
SERVICE_BACKEND_PORT_PORT-8080=8080

Using the '--bind-as-files' flag, secrets will be accessible as files instead of environment variables.
The value of the '--name' flag indicates the name of the directory under '/bindings/' containing the secrets files.
`
)

// LinkOptions encapsulates the options for the odo link command
type LinkOptions struct {
	// Common link/unlink context
	*commonLinkOptions

	// Flags
	contextFlag string
}

// NewLinkOptions creates a new LinkOptions instance
func NewLinkOptions() *LinkOptions {
	options := LinkOptions{}
	options.commonLinkOptions = newCommonLinkOptions()
	options.commonLinkOptions.serviceBinding = &servicebinding.ServiceBinding{}
	return &options
}

// Complete completes LinkOptions after they've been created
func (o *LinkOptions) Complete(cmdline cmdline.Cmdline, args []string) (err error) {
	o.commonLinkOptions.devfilePath = location.DevfileLocation(o.contextFlag)

	err = o.complete(cmdline, args, o.contextFlag)
	if err != nil {
		return err
	}

	if o.csvSupport {
		o.operation = o.KClient.LinkSecret
	}
	return err
}

// Validate validates the LinkOptions based on completed values
func (o *LinkOptions) Validate() (err error) {
	return o.validate()
}

// Run contains the logic for the odo link command
func (o *LinkOptions) Run() (err error) {
	return o.run()
}

// NewCmdLink implements the link odo command
func NewCmdLink(name, fullName string) *cobra.Command {
	o := NewLinkOptions()

	linkCmd := &cobra.Command{
		Use:         fmt.Sprintf("%s <operator-service-type>/<service-name> OR %s <operator-service-type>/<service-name> --component [component] OR %s <component> --component [component]", name, name, name),
		Short:       "Link component to a service or component",
		Long:        linkLongDesc,
		Example:     fmt.Sprintf(linkExample, fullName),
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"command": "component"},
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	linkCmd.PersistentFlags().BoolVarP(&o.inlined, "inlined", "", false, "Puts the link definition in the devfile instead of a separate file")
	linkCmd.PersistentFlags().StringVar(&o.name, "name", "", "Name of the created ServiceBinding resource")
	linkCmd.PersistentFlags().BoolVar(&o.bindAsFiles, "bind-as-files", false, "If enabled, configuration values will be mounted as files, instead of declared as environment variables")
	linkCmd.PersistentFlags().StringArrayVarP(&o.mappings, "map", "", []string{}, "Mappings (custom binding data) to be added to the component; each map should be specified as <key>=<value>")
	linkCmd.SetUsageTemplate(odoutil.CmdUsageTemplate)

	//Adding `--component` flag
	AddComponentFlag(linkCmd)

	//Adding context flag
	odoutil.AddContextFlag(linkCmd, &o.contextFlag)

	return linkCmd
}
