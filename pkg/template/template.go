package template

import (
	"context"
	"fmt"
	"io"

	"github.com/operator-framework/operator-registry/alpha/declcfg"
)

type TemplateExpanderInterface interface {
	Expand(context.Context) (*declcfg.DeclarativeConfig, error)
}

type TemplateOptions struct {
	Input        io.Reader
	RenderBundle func(context.Context, string) (*declcfg.DeclarativeConfig, error)
}

type Template struct {
	Schema string `json:"schema"`
}

func (t TemplateOptions) Expand(ctx context.Context) (*declcfg.DeclarativeConfig, error) {
	return nil, fmt.Errorf("not implemented")
}
