package template

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/operator-framework/operator-registry/alpha/declcfg"
)

type BasicOptions struct {
	TemplateOptions
}

type BasicTemplate struct {
	Template
	Entries []*declcfg.Meta `json:"entries"`
}

func parseSpec(reader io.Reader) (*BasicTemplate, error) {
	bt := &BasicTemplate{}
	btDoc := json.RawMessage{}
	btDecoder := yaml.NewYAMLOrJSONDecoder(reader, 4096)
	err := btDecoder.Decode(&btDoc)
	if err != nil {
		return nil, fmt.Errorf("decoding template schema: %v", err)
	}
	err = json.Unmarshal(btDoc, bt)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling template: %v", err)
	}

	return bt, nil
}

func (t BasicOptions) Render(ctx context.Context) (*declcfg.DeclarativeConfig, error) {
	bt, err := parseSpec(t.Input)
	if err != nil {
		return nil, err
	}
	cfg, err := declcfg.LoadSlice(bt.Entries)
	if err != nil {
		return cfg, err
	}

	outb := cfg.Bundles[:0]
	for _, b := range cfg.Bundles {
		if !isBundleTemplate(&b) {
			return nil, fmt.Errorf("unexpected fields present in basic template bundle")
		}
		contributor, err := t.TemplateOptions.RenderBundle(ctx, b.Image)
		if err != nil {
			return nil, err
		}
		outb = append(outb, contributor.Bundles...)
	}

	cfg.Bundles = outb
	return cfg, nil
}

// isBundleTemplate identifies a Bundle template source any valid FBC entity as having a Schema and Image defined
// but no Properties, RelatedImages or Package defined
func isBundleTemplate(b *declcfg.Bundle) bool {
	return b.Schema != "" && b.Image != "" && b.Package == "" && len(b.Properties) == 0 && len(b.RelatedImages) == 0
}

// FromReader reads FBC from a reader and generates a BasicTemplate from it
func FromReader(r io.Reader) (*BasicTemplate, error) {
	var entries []*declcfg.Meta
	if err := declcfg.WalkMetasReader(r, func(meta *declcfg.Meta, err error) error {
		if err != nil {
			return err
		}
		if meta.Schema == declcfg.SchemaBundle {
			var b declcfg.Bundle
			if err := json.Unmarshal(meta.Blob, &b); err != nil {
				return fmt.Errorf("parse bundle: %v", err)
			}
			b2 := declcfg.Bundle{
				Schema: b.Schema,
				Image:  b.Image,
			}
			meta.Blob, err = json.Marshal(b2)
			if err != nil {
				return fmt.Errorf("re-serialize bundle: %v", err)
			}
		}
		entries = append(entries, meta)
		return nil
	}); err != nil {
		return nil, err
	}

	bt := &BasicTemplate{
		Entries: entries,
	}

	return bt, nil
}
