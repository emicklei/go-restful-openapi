package restfulspec

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

// nolint:paralleltest
func TestThatExtraTagsAreReadIntoModel(t *testing.T) {
	type fakeint int
	type fakearray string
	type Anything struct {
		Name             string    `description:"name" modelDescription:"a test" readOnly:"false"`
		Size             int       `minimum:"0" maximum:"10"`
		Stati            string    `enum:"off|on" default:"on" modelDescription:"more description"`
		ID               string    `unique:"true"`
		FakeInt          fakeint   `type:"integer"`
		FakeArray        fakearray `type:"[]string"`
		IP               net.IP    `type:"string"`
		Password         string
		Optional         bool   `optional:"true"`
		Created          string `readOnly:"true"`
		NullableField    string `x-nullable:"true"`
		NotNullableField string `x-nullable:"false"`
		UUID             string `type:"string" format:"UUID"`
	}
	d := definitionsFromStruct(Anything{})
	props, _ := d["restfulspec.Anything"]
	p1, _ := props.Properties["Name"]
	if got, want := p1.Description, "name"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := p1.ReadOnly, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p2, _ := props.Properties["Size"]
	if got, want := *p2.Minimum, 0.0; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := p2.ReadOnly, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := *p2.Maximum, 10.0; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p3, _ := props.Properties["Stati"]
	if got, want := p3.Enum[0], "off"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := p3.Enum[1], "on"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p4, _ := props.Properties["ID"]
	if got, want := p4.UniqueItems, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p5, _ := props.Properties["Password"]
	if got, want := p5.Type[0], "string"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p6, _ := props.Properties["FakeInt"]
	if got, want := p6.Type[0], "integer"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p7, _ := props.Properties["FakeArray"]
	if got, want := p7.Type[0], "array"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p7p, _ := props.Properties["FakeArray"]
	if got, want := p7p.Items.Schema.Type[0], "string"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p8, _ := props.Properties["IP"]
	if got, want := p8.Type[0], "string"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p9, _ := props.Properties["Created"]
	if got, want := p9.ReadOnly, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := strings.Contains(fmt.Sprintf("%v", props.Required), "Optional"), false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := props.Description, "a test\nmore description"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p10, _ := props.Properties["NullableField"]
	if got, want := p10.Extensions["x-nullable"], true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p11, _ := props.Properties["NotNullableField"]
	if got, want := p11.Extensions["x-nullable"], false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p12, _ := props.Properties["UUID"]
	if got, want := p12.Type[0], "string"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := p12.Format, "UUID"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
