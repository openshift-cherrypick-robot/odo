package service

import (
	"fmt"

	"github.com/redhat-developer/odo/pkg/odo/cmdline"
	"github.com/redhat-developer/odo/pkg/odo/genericclioptions"
	odoutil "github.com/redhat-developer/odo/pkg/odo/util"
	"github.com/redhat-developer/odo/pkg/service"
	svc "github.com/redhat-developer/odo/pkg/service"
	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const describeRecommendedCommandName = "describe"

var (
	describeExample = ktemplates.Examples(`
    # Describe the service named 'mysql-persistent'
    %[1]s mysql-persistent`)

	describeLongDesc = ktemplates.LongDesc(`
	Describe an existing service, either defined locally or deployed to the cluster`)
)

// DescribeOptions encapsulates the options for the odo service describe command
type DescribeOptions struct {
	// Context
	*genericclioptions.Context

	// Parameters
	serviceName string

	// Flags
	contextFlag string

	// Backend is the service provider backend that was used to create the service
	Backend ServiceProviderBackend
}

// NewDescribeOptions creates a new DescribeOptions instance
func NewDescribeOptions() *DescribeOptions {
	return &DescribeOptions{}
}

func (o *DescribeOptions) Complete(cmdline cmdline.Cmdline, args []string) (err error) {
	o.Context, err = genericclioptions.New(genericclioptions.NewCreateParameters(cmdline).NeedDevfile(o.contextFlag))
	if err != nil {
		return err
	}

	err = validDevfileDirectory(o.contextFlag)
	if err != nil {
		return err
	}

	o.serviceName = args[0]
	_, _, err = service.SplitServiceKindName(o.serviceName)
	if err != nil {
		return fmt.Errorf("invalid service name")
	}
	o.Backend = NewOperatorBackend()

	return nil
}

// Validate validates the DescribeOptions based on completed values
func (o *DescribeOptions) Validate() error {
	svcDefined, err := o.Backend.ServiceDefined(o.Context, o.serviceName)
	if err != nil {
		return err
	}

	svcDeployed, err := svc.OperatorSvcExists(o.KClient, o.serviceName)
	if err != nil {
		return err
	}

	if !svcDefined && !svcDeployed {
		return fmt.Errorf("couldn't find service named %q. Refer %q to see list of defined services", o.serviceName, "odo service list")
	}
	return nil
}

// Run contains the logic for the odo service describe command
func (o *DescribeOptions) Run() error {
	return o.Backend.DescribeService(o, o.serviceName, o.GetApplication())
}

// NewCmdDescribe implements the describe odo command
func NewCmdServiceDescribe(name, fullName string) *cobra.Command {
	do := NewDescribeOptions()

	var describeCmd = &cobra.Command{
		Use:         fmt.Sprintf("%s [service_name]", name),
		Short:       "Describe an existing service",
		Long:        describeLongDesc,
		Example:     fmt.Sprintf(describeExample, fullName),
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"machineoutput": "json"},
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(do, cmd, args)
		},
	}

	odoutil.AddContextFlag(describeCmd, &do.contextFlag)
	return describeCmd
}
