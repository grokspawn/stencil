# stencil
Stencil is a re-imagining of [catalog templates](https://olm.operatorframework.io/docs/reference/catalog-templates/) operations for [file-based catalogs](https://olm.operatorframework.io/docs/reference/file-based-catalogs/) (FBC) as a separate binary than [opm](https://github.com/operator-framework/operator-registry/releases).

`opm` is a binary used by FBC projects for:
- bundle creation and validation
- FBC creation and validation
- migration to FBC from legacy catalog formats
- expansion of catalog templates to FBC

These roles occur throughout the bundle/catalog publishing lifecycle, so it historically has been a real challenge ensuring version/feature parity by all components of conglomerated toolchains.

The advantage of a discrete binary is that it will allow those concerns to be separated from other `opm` roles, allowing for a separate statement of support and versioning axis.

## Catalog Template API updates
All catalog templates share two components:
- an initializing `TemplateOptions` instance, which provides input data and customizable bundle rendering behavior
```golang
type TemplateOptions struct {
	Input        io.Reader
	RenderBundle func(context.Context, string) (*declcfg.DeclarativeConfig, error)
}
```
- an intermediate type for interpreting the catalog template data which implements a base `Template` type
```golang
type Template struct {
	Schema string `json:"schema"`
}
```

## Commands
### expand
`stencil` can expand any supported catalog template format into generated FBC, and will automagically identify the template format to assure correct expansion.

```sh
./bin/stencil expand --help
Generate a file-based catalog from a catalog template file
When FILE is '-' or not provided, the template is read from standard input

Usage:
  stencil expand [FILE] [flags]

Flags:
  -h, --help                   help for expand
      --migrate-level string   Name of the last migration to run (default: none)

                               The migrator will run all migrations up to and including the selected level.

                               Available migrators:
                                 - none                          : do nothing
                                 - bundle-object-to-csv-metadata : migrates bundles' "olm.bundle.object" to "olm.csv.metadata"

  -o, --output string          Output format (json|yaml) (default "json")

Global Flags:
      --skip-tls-verify   skip TLS certificate verification for container image registries while pulling bundles
      --use-http          use plain HTTP for container image registries while pulling bundles
```

### convert
`stencil` can convert existing FBC to a basic catalog template instance using the subcommand `basic`

#### basic
```sh
./bin/stencil convert basic --help
Generate a basic template from existing FBC.

This command outputs a basic catalog template to STDOUT from input FBC.
If no argument is specified or is '-' input is assumed from STDIN.

Usage:
  stencil convert basic [<fbc-file> | -] [flags]

Flags:
  -h, --help            help for basic
  -o, --output string   Output format (json|yaml) (default "json")

Global Flags:
      --skip-tls-verify   skip TLS certificate verification for container image registries while pulling bundles
      --use-http          use plain HTTP for container image registries while pulling bundles
```
