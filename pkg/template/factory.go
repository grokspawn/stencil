package template

import (
	"bytes"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/util/yaml"
)

const BasicSchema string = "olm.template.basic"
const SemverSchema string = "olm.semver"

func NewExpander(to TemplateOptions) (TemplateExpanderInterface, error) {
	b, err := io.ReadAll(to.Input)
	if err != nil {
		return nil, err
	}
	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(b), 4096)

	var in Template
	var out TemplateExpanderInterface
	if err := dec.Decode(&in); err != nil {
		return nil, err
	}
	to.Input = bytes.NewReader(b)
	switch in.Schema {
	case BasicSchema:
		out = &BasicOptions{TemplateOptions: to}
	case SemverSchema:
		out = &SemverOptions{TemplateOptions: to}
	default:
		return nil, fmt.Errorf("unknown schema %q", in.Schema)
	}

	return out, nil
}
