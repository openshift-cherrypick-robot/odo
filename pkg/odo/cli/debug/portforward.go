package debug

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/redhat-developer/odo/pkg/debug"
	"github.com/redhat-developer/odo/pkg/devfile/location"
	"github.com/redhat-developer/odo/pkg/log"
	"github.com/redhat-developer/odo/pkg/odo/cmdline"
	"github.com/redhat-developer/odo/pkg/odo/genericclioptions"
	odoutil "github.com/redhat-developer/odo/pkg/odo/util"
	"github.com/redhat-developer/odo/pkg/util"

	"github.com/spf13/cobra"

	k8sgenclioptions "k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/util/templates"
)

const (
	// DefaultDebugPort is the default port used for debugging on remote pod
	DefaultDebugPort = 5858
)

// PortForwardOptions contains all the options for running the port-forward cli command.
type PortForwardOptions struct {
	// Context
	*genericclioptions.Context

	// Flags
	contextFlag   string
	localPortFlag int

	// PortPair is the combination of local and remote port in the format "local:remote"
	PortPair string

	// Port forwarder backend
	PortForwarder *debug.DefaultPortForwarder

	// StopChannel is used to stop port forwarding
	StopChannel chan struct{}

	// ReadChannel is used to receive status of port forwarding ( ready or not ready )
	ReadyChannel chan struct{}
}

var (
	portforwardLong = templates.LongDesc(`Forward a local port to a remote port on the pod where the application is listening for a debugger. By default the local port and the remote port will be same. To change the local port you can use --local-port argument and to change the remote port use "odo env set DebugPort <port>"   		  
	`)

	portforwardExample = templates.Examples(`
		# Listen on default port and forwarding to the default port in the pod
		odo debug port-forward 

		# Listen on the 5000 port locally, forwarding to default port in the pod
		odo debug port-forward --local-port 5000
		
		`)
)

const (
	portforwardCommandName = "port-forward"
)

// NewPortForwardOptions returns the PortForwardOptions struct
func NewPortForwardOptions() *PortForwardOptions {
	return &PortForwardOptions{}
}

// Complete completes all the required options for port-forward cmd.
func (o *PortForwardOptions) Complete(cmdline cmdline.Cmdline, args []string) (err error) {
	o.Context, err = genericclioptions.New(genericclioptions.NewCreateParameters(cmdline))
	if err != nil {
		return err
	}

	remotePort := o.Context.EnvSpecificInfo.GetDebugPort()

	// try to listen on the given local port and check if the port is free or not
	addressLook := "localhost:" + strconv.Itoa(o.localPortFlag)
	listener, err := net.Listen("tcp", addressLook)
	if err != nil {
		// if the local-port flag is set by the user, return the error and stop execution
		if cmdline.IsFlagSet("local-port") {
			return err
		}
		// else display a error message and auto select a new free port
		log.Errorf("the local debug port %v is not free, cause: %v", o.localPortFlag, err)
		o.localPortFlag, err = util.HTTPGetFreePort()
		if err != nil {
			return err
		}
		log.Infof("The local port %v is auto selected", o.localPortFlag)
	} else {
		err = listener.Close()
		if err != nil {
			return err
		}
	}

	o.PortPair = fmt.Sprintf("%d:%d", o.localPortFlag, remotePort)

	// Using Discard streams because nothing important is logged
	o.PortForwarder = debug.NewDefaultPortForwarder(o.Context.EnvSpecificInfo.GetName(), o.Context.GetApplication(), o.Context.EnvSpecificInfo.GetNamespace(), o.KClient, k8sgenclioptions.NewTestIOStreamsDiscard())

	o.StopChannel = make(chan struct{}, 1)
	o.ReadyChannel = make(chan struct{})
	return nil
}

// Validate validates all the required options for port-forward cmd.
func (o PortForwardOptions) Validate() error {
	if len(o.PortPair) < 1 {
		return fmt.Errorf("ports cannot be empty")
	}
	return nil
}

// Run implements all the necessary functionality for port-forward cmd.
func (o PortForwardOptions) Run() error {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer signal.Stop(signals)
	defer os.RemoveAll(debug.GetDebugInfoFilePath(o.Context.EnvSpecificInfo.GetName(), o.Context.GetApplication(), o.Context.EnvSpecificInfo.GetNamespace()))

	go func() {
		<-signals
		if o.StopChannel != nil {
			close(o.StopChannel)
		}
	}()

	err := debug.CreateDebugInfoFile(o.PortForwarder, o.PortPair)
	if err != nil {
		return err
	}

	devfilePath := location.DevfileLocation(o.contextFlag)
	return o.PortForwarder.ForwardPorts(o.PortPair, o.StopChannel, o.ReadyChannel, util.CheckPathExists(devfilePath))
}

// NewCmdPortForward implements the port-forward odo command
func NewCmdPortForward(name, fullName string) *cobra.Command {

	opts := NewPortForwardOptions()
	cmd := &cobra.Command{
		Use:     name,
		Short:   "Forward one or more local ports to a pod",
		Long:    portforwardLong,
		Example: portforwardExample,
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(opts, cmd, args)
		},
	}

	odoutil.AddContextFlag(cmd, &opts.contextFlag)
	cmd.Flags().IntVarP(&opts.localPortFlag, "local-port", "l", DefaultDebugPort, "Set the local port")

	return cmd
}
