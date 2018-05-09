package restfulspec

import "testing"
import (
	"github.com/go-openapi/spec"
)

type Apple struct {
	Species string
	Volume  int `json:"vol"`
	Things  *[]string
}

func TestAppleDef(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(Apple{})

	schema := db.Definitions["restfulspec.Apple"]
	if got, want := len(schema.Required), 3; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Required[0], "Species"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Required[1], "vol"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.ID, ""; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Properties["Things"].Items.Schema.Type.Contains("string"), true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Properties["Things"].Items.Schema.Ref.String(), ""; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
