package expand

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/grokspawn/stencil/internal/util"
	"github.com/grokspawn/stencil/pkg/template"
	"github.com/operator-framework/operator-registry/alpha/action"
	"github.com/operator-framework/operator-registry/alpha/action/migrations"
	"github.com/operator-framework/operator-registry/alpha/declcfg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var (
		output       string
		migrateLevel string
	)

	cmd := &cobra.Command{
		Use: "expand [FILE]",
		Short: `Generate a file-based catalog from a catalog template file
When FILE is '-' or not provided, the template is read from standard input`,
		Long: `Generate a file-based catalog from a catalog template file
When FILE is '-' or not provided, the template is read from standard input`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle different input argument types
			// When no arguments or "-" is passed to the command,
			// assume input is coming from stdin
			// Otherwise open the file passed to the command
			data, source, err := util.OpenFileOrStdin(cmd, args)
			if err != nil {
				return err
			}
			defer data.Close()

			var write func(declcfg.DeclarativeConfig, io.Writer) error
			switch output {
			case "json":
				write = declcfg.WriteJSON
			case "yaml":
				write = declcfg.WriteYAML
			case "mermaid":
				write = func(cfg declcfg.DeclarativeConfig, writer io.Writer) error {
					mermaidWriter := declcfg.NewMermaidWriter()
					return mermaidWriter.WriteChannels(cfg, writer)
				}
			default:
				return fmt.Errorf("invalid output format %q", output)
			}

			// The bundle loading impl is somewhat verbose, even on the happy path,
			// so discard all logrus default logger logs. Any important failures will be
			// returned from template.Render and logged as fatal errors.
			logrus.SetOutput(io.Discard)

			reg, err := util.CreateCLIRegistry(cmd)
			if err != nil {
				log.Fatalf("creating containerd registry: %v", err)
			}
			defer reg.Destroy()

			var m *migrations.Migrations
			if migrateLevel != "" {
				m, err = migrations.NewMigrations(migrateLevel)
				if err != nil {
					log.Fatal(err)
				}
			}

			options, err := template.NewExpander(template.TemplateOptions{
				Input: data,
				RenderBundle: func(ctx context.Context, ref string) (*declcfg.DeclarativeConfig, error) {
					renderer := action.Render{
						Refs:           []string{ref},
						Registry:       reg,
						AllowedRefMask: action.RefBundleImage,
						Migrations:     m,
					}
					return renderer.Run(ctx)
				}})
			if err != nil {
				log.Fatalf("detecting template type for %q: %v", source, err)
			}

			cfg, err := options.Expand(cmd.Context())
			if err != nil {
				log.Fatalf("rendering %q: %v", source, err)
			}

			if cfg != nil {
				if err := write(*cfg, os.Stdout); err != nil {
					log.Fatal(err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "json", "Output format (json|yaml)")
	cmd.Flags().StringVar(&migrateLevel, "migrate-level", "", "Name of the last migration to run (default: none)\n"+migrations.HelpText())

	return cmd
}
