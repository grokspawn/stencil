package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/grokspawn/stencil/pkg/template"
	"github.com/operator-framework/operator-registry/pkg/image"
	"sigs.k8s.io/yaml"
)

type Converter struct {
	FbcReader    io.Reader
	OutputFormat string
	Registry     image.Registry
}

func (c *Converter) Convert() error {
	bt, err := template.FromReader(c.FbcReader)
	if err != nil {
		return err
	}

	b, _ := json.MarshalIndent(bt, "", "    ")
	if c.OutputFormat == "json" {
		fmt.Fprintln(os.Stdout, string(b))
	} else {
		y, err := yaml.JSONToYAML(b)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(y))
	}

	return nil
}
