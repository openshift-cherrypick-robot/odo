package search

import (
	"fmt"

	olm "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-developer/odo/pkg/odo/cli/catalog/util"
	"github.com/redhat-developer/odo/pkg/odo/cmdline"
	"github.com/redhat-developer/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"
)

const serviceRecommendedCommandName = "service"

var serviceExample = `  # Search for a service
  %[1]s mysql`

// SearchServiceOptions encapsulates the options for the odo catalog describe service command
type SearchServiceOptions struct {
	// Context
	*genericclioptions.Context

	// Parameters
	searchTerm string

	// list of clusterserviceversions (installed by Operators)
	csvs *olm.ClusterServiceVersionList
}

// NewSearchServiceOptions creates a new SearchServiceOptions instance
func NewSearchServiceOptions() *SearchServiceOptions {
	return &SearchServiceOptions{}
}

// Complete completes SearchServiceOptions after they've been created
func (o *SearchServiceOptions) Complete(cmdline cmdline.Cmdline, args []string) (err error) {
	o.Context, err = genericclioptions.New(genericclioptions.NewCreateParameters(cmdline))
	if err != nil {
		return err
	}

	o.searchTerm = args[0]

	o.csvs, err = o.KClient.SearchClusterServiceVersionList(o.searchTerm)
	if err != nil {
		return fmt.Errorf("unable to list services because Operator Hub is not enabled in your cluster: %v", err)
	}

	return nil
}

// Validate validates the SearchServiceOptions based on completed values
func (o *SearchServiceOptions) Validate() error {
	if len(o.csvs.Items) == 0 {
		return fmt.Errorf("no service matched the query: %s", o.searchTerm)
	}
	return nil
}

// Run contains the logic for the command associated with SearchServiceOptions
func (o *SearchServiceOptions) Run() error {
	if len(o.csvs.Items) > 0 {
		util.DisplayClusterServiceVersions(o.csvs)
	}
	return nil
}

// NewCmdCatalogSearchService implements the odo catalog search service command
func NewCmdCatalogSearchService(name, fullName string) *cobra.Command {
	o := NewSearchServiceOptions()
	return &cobra.Command{
		Use:   name,
		Short: "Search service type in catalog",
		Long: `Search service type in catalog.

This searches for a partial match for the given search term in all the available
services from operator hub services.
`,
		Example: fmt.Sprintf(serviceExample, fullName),
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

}
