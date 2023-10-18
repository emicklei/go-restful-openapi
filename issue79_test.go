package restfulspec

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/spec"
)

type Parent struct {
	FieldA string
}

func (v Parent) MarshalJSON() ([]byte, error) { return nil, nil }

type Child struct {
	FieldA Parent
	FieldB int
}

func TestParentChildArray(t *testing.T) {
	t.Skip()
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(Child{})
	s := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Definitions: db.Definitions,
		},
	}
	data, _ := json.MarshalIndent(s, "", "  ")
	t.Log(string(data))
}
