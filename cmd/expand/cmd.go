package expand

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var output string

	runCmd := &cobra.Command{
		Use:   "expand",
		Short: "Produce a catalog from a template",
		Args:  cobra.NoArgs,
	}

	bc := newBasicTemplateCmd()
	// bc.Hidden = true
	runCmd.AddCommand(bc)

	sc := newSemverTemplateCmd()
	// sc.Hidden = true
	runCmd.AddCommand(sc)

	runCmd.PersistentFlags().StringVarP(&output, "output", "o", "json", "Output format (json|yaml)")

	return runCmd
}
