package template

import (
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/util/yaml"
)

const BasicSchema string = "olm.template.basic"
const SemverSchema string = "olm.semver"

func NewTemplate(r io.Reader) (TemplateOptionInterface, error) {
	dec := yaml.NewYAMLOrJSONDecoder(r, 4096)

	var in Template
	var out TemplateOptionInterface
	if err := dec.Decode(&in); err != nil {
		return nil, err
	}
	switch in.Schema {
	case BasicSchema:
		out = &BasicOptions{}
	case SemverSchema:
		out = &SemverOptions{}
	default:
		return nil, fmt.Errorf("unknown schema %q", in.Schema)
	}

	return out, nil
}
