package bootstrap

import (
	"encoding/json"

	"github.com/faber-numeris/foundation/beholder/api"
	"gopkg.in/yaml.v3"
)

// openAPISpecGenerator implements spec-ui's SpecGenerator interface,
// serving the spec embedded in the generated api package instead of
// requiring a spec file on disk.
type openAPISpecGenerator struct{}

func (openAPISpecGenerator) MarshalJSON() ([]byte, error) {
	return api.GetSpecJSON()
}

func (openAPISpecGenerator) MarshalYAML() ([]byte, error) {
	raw, err := api.GetSpecJSON()
	if err != nil {
		return nil, err
	}

	var spec any
	if err := json.Unmarshal(raw, &spec); err != nil {
		return nil, err
	}

	return yaml.Marshal(spec)
}
