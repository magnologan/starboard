package cmd

import (
	"fmt"
	"os"

	vulnsCrd "github.com/aquasecurity/starboard/pkg/find/vulnerabilities/crd"
	clientset "github.com/aquasecurity/starboard/pkg/generated/clientset/versioned"
	configAuditCrd "github.com/aquasecurity/starboard/pkg/polaris/crd"
	"github.com/aquasecurity/starboard/pkg/report"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewGetReportCmd(cf *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report (NAME | TYPE/NAME)",
		Short: "Get a full html security report for a specified workload",
		Long: `Generates a report that contains vulnerabilities and config audits found for the specified workload

TYPE is a Kubernetes workload. Shortcuts and API groups will be resolved, e.g. 'po' or 'deployments.apps'.
NAME is the name of a particular Kubernetes workload.			
`,
		Example: fmt.Sprintf(`  # Save report to a file
  %[1]s get report deploy/nginx > report.html`, "starboard"),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			config, err := cf.ToRESTConfig()
			if err != nil {
				return
			}
			starboardClientset, err := clientset.NewForConfig(config)
			if err != nil {
				return
			}
			ns, _, err := cf.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				return
			}
			workload, err := WorkloadFromArgs(ns, args)
			if err != nil {
				return
			}

			caReader := configAuditCrd.NewReadWriter(starboardClientset)
			vulnsReader := vulnsCrd.NewReadWriter(starboardClientset)

			reporter := report.NewHTMLReporter(caReader, vulnsReader, workload)
			err = reporter.GenerateReport(os.Stdout)

			return
		},
	}

	return cmd
}
