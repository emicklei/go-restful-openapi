package restfulspec

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/go-openapi/spec"
)

type Child struct {
	SingleByteArray []byte   `json:"sba,omitempty" `
	DoubleByteArray [][]byte `json:"dba,omitempty" `
}

func TestParentChildArray(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(Child{})
	s := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Definitions: db.Definitions,
		},
	}
	data, _ := json.MarshalIndent(s, "", "  ")
	log.Fatalln(string(data))
}
