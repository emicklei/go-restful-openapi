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
		XGoName          string `x-go-name:"specgoname"`
		ByteArray        []byte `format:"binary"`
	}
	d := definitionsFromStruct(Anything{})
	props := d["restfulspec.Anything"]
	p1 := props.Properties["Name"]
	if got, want := p1.Description, "name"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := p1.ReadOnly, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p2 := props.Properties["Size"]
	if got, want := *p2.Minimum, 0.0; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := p2.ReadOnly, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := *p2.Maximum, 10.0; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p3 := props.Properties["Stati"]
	if got, want := p3.Enum[0], "off"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := p3.Enum[1], "on"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p4 := props.Properties["ID"]
	if got, want := p4.UniqueItems, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p5 := props.Properties["Password"]
	if got, want := p5.Type[0], "string"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p6 := props.Properties["FakeInt"]
	if got, want := p6.Type[0], "integer"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p7 := props.Properties["FakeArray"]
	if got, want := p7.Type[0], "array"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p7p := props.Properties["FakeArray"]
	if got, want := p7p.Items.Schema.Type[0], "string"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p8 := props.Properties["IP"]
	if got, want := p8.Type[0], "string"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p9 := props.Properties["Created"]
	if got, want := p9.ReadOnly, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := strings.Contains(fmt.Sprintf("%v", props.Required), "Optional"), false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := props.Description, "a test\nmore description"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p10 := props.Properties["NullableField"]
	if got, want := p10.Extensions["x-nullable"], true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p11 := props.Properties["NotNullableField"]
	if got, want := p11.Extensions["x-nullable"], false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p12 := props.Properties["UUID"]
	if got, want := p12.Type[0], "string"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := p12.Format, "UUID"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	p13 := props.Properties["XGoName"]
	if got, want := p13.Extensions["x-go-name"], "specgoname"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
