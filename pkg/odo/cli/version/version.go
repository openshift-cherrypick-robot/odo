package version

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/redhat-developer/odo/pkg/kclient"
	"github.com/redhat-developer/odo/pkg/odo/cmdline"
	"github.com/redhat-developer/odo/pkg/odo/genericclioptions"
	"github.com/redhat-developer/odo/pkg/preference"
	odoversion "github.com/redhat-developer/odo/pkg/version"

	"github.com/redhat-developer/odo/pkg/notify"
	"github.com/redhat-developer/odo/pkg/odo/util"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

// RecommendedCommandName is the recommended version command name
const RecommendedCommandName = "version"

// OdoReleasesPage is the GitHub page where we do all our releases
const OdoReleasesPage = "https://github.com/redhat-developer/odo/releases"

var versionLongDesc = ktemplates.LongDesc("Print the client version information")

var versionExample = ktemplates.Examples(`
# Print the client version of odo
%[1]s`,
)

// VersionOptions encapsulates all options for odo version command
type VersionOptions struct {
	// Flags
	clientFlag bool

	// serverInfo contains the remote server information if the user asked for it, nil otherwise
	serverInfo *kclient.ServerInfo

	prefClient preference.Client
}

// NewVersionOptions creates a new VersionOptions instance
func NewVersionOptions(prefClient preference.Client) *VersionOptions {
	return &VersionOptions{
		prefClient: prefClient,
	}
}

// Complete completes VersionOptions after they have been created
func (o *VersionOptions) Complete(cmdline cmdline.Cmdline, args []string) (err error) {
	if !o.clientFlag {
		// Let's fetch the info about the server, ignoring errors
		client, err := kclient.New()

		if err == nil {
			// checking the value of timeout in preference
			var timeout time.Duration
			if o.prefClient != nil {
				timeout = time.Duration(o.prefClient.GetTimeout()) * time.Second
			} else {
				// the default timeout will be used
				// when the value is not readable from preference
				timeout = preference.DefaultTimeout * time.Second
			}
			o.serverInfo, _ = client.GetServerVersion(timeout)
		}
	}
	return nil
}

// Validate validates the VersionOptions based on completed values
func (o *VersionOptions) Validate() (err error) {
	return nil
}

// Run contains the logic for the odo service create command
func (o *VersionOptions) Run() (err error) {
	// If verbose mode is enabled, dump all KUBECLT_* env variables
	// this is usefull for debuging oc plugin integration
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "KUBECTL_") {
			klog.V(4).Info(v)
		}
	}

	fmt.Println("odo " + odoversion.VERSION + " (" + odoversion.GITCOMMIT + ")")

	if !o.clientFlag && o.serverInfo != nil {
		// make sure we only include OpenShift info if we actually have it
		openshiftStr := ""
		if len(o.serverInfo.OpenShiftVersion) > 0 {
			openshiftStr = fmt.Sprintf("OpenShift: %v\n", o.serverInfo.OpenShiftVersion)
		}
		fmt.Printf("\n"+
			"Server: %v\n"+
			"%v"+
			"Kubernetes: %v\n",
			o.serverInfo.Address,
			openshiftStr,
			o.serverInfo.KubernetesVersion)
	}

	return nil
}

// NewCmdVersion implements the version odo command
func NewCmdVersion(name, fullName string) *cobra.Command {
	prefClient, err := preference.NewClient()
	if err != nil {
		klog.V(3).Info(errors.Wrap(err, "unable to read preference file"))
	}
	o := NewVersionOptions(prefClient)
	// versionCmd represents the version command
	var versionCmd = &cobra.Command{
		Use:     name,
		Short:   versionLongDesc,
		Long:    versionLongDesc,
		Example: fmt.Sprintf(versionExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	// Add a defined annotation in order to appear in the help menu
	versionCmd.Annotations = map[string]string{"command": "utility"}
	versionCmd.SetUsageTemplate(util.CmdUsageTemplate)
	versionCmd.Flags().BoolVar(&o.clientFlag, "client", false, "Client version only (no server required).")

	return versionCmd
}

// GetLatestReleaseInfo Gets information about the latest release
func GetLatestReleaseInfo(info chan<- string) {
	newTag, err := notify.CheckLatestReleaseTag(odoversion.VERSION)
	if err != nil {
		// The error is intentionally not being handled because we don't want
		// to stop the execution of the program because of this failure
		klog.V(4).Infof("Error checking if newer odo release is available: %v", err)
	}
	if len(newTag) > 0 {
		info <- fmt.Sprintf(`
---
A newer version of odo (%s) is available,
visit %s to update.
If you wish to disable this notification, run:
odo preference set UpdateNotification false
---`, fmt.Sprint(newTag), OdoReleasesPage)

	}
}
