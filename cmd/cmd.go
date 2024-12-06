package cmdroot

import (
	"github.com/grokspawn/stencil/cmd/convert"
	"github.com/grokspawn/stencil/cmd/expand"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "stencil",
		Short: "FBC catalog template manipulation tool",
		Args:  cobra.NoArgs,
		Run:   func(_ *cobra.Command, _ []string) {}, // adding an empty function here to preserve non-zero exit status for misstated subcommands/flags for the command hierarchy
	}

	cmd.AddCommand(expand.NewCmd())
	cmd.AddCommand(convert.NewCmd())
	cmd.PersistentFlags().StringVarP(&output, "output", "o", "json", "Output format (json|yaml)")
	cmd.PersistentFlags().Bool("skip-tls", false, "skip TLS certificate verification for container image registries while pulling bundles or index")
	cmd.PersistentFlags().Bool("skip-tls-verify", false, "skip TLS certificate verification for container image registries while pulling bundles")
	cmd.PersistentFlags().Bool("use-http", false, "use plain HTTP for container image registries while pulling bundles")
	if err := cmd.PersistentFlags().MarkDeprecated("skip-tls", "use --use-http and --skip-tls-verify instead"); err != nil {
		logrus.Panic(err.Error())
	}

	return cmd
}
