package template

import (
	"context"
	"fmt"
	"io"

	"github.com/operator-framework/operator-registry/alpha/declcfg"
)

type TemplateOptionInterface interface {
	Render(context.Context) (*declcfg.DeclarativeConfig, error)
}

type TemplateOptions struct {
	RenderBundle func(context.Context, string) (*declcfg.DeclarativeConfig, error)
	Input        io.Reader
	Output       io.Writer
}

type Template struct {
	Schema string `json:"schema"`
}

func (t TemplateOptions) Render(ctx context.Context) (*declcfg.DeclarativeConfig, error) {
	return nil, fmt.Errorf("not implemented")
}