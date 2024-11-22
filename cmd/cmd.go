package cmdroot

import (
	"github.com/grokspawn/stencil/cmd/convert"
	"github.com/grokspawn/stencil/cmd/expand"
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

	return cmd
}
