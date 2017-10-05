package main

import (
	"testing"
	"io/ioutil"
	"github.com/emicklei/go-restful-openapi"
	"github.com/emicklei/go-restful"
	spec2 "github.com/go-openapi/spec"
	"encoding/json"
	"reflect"
)

func TestAppleDef(t *testing.T) {

	raw, err := ioutil.ReadFile("./openapi.json")
	if err != nil {
		t.Error("Loading the openapi specification failed.")
	}
	ws := UserResource{}.WebService()
	expected := &spec2.Swagger{}
	if err := json.Unmarshal(raw, expected); err != nil {
		t.Error("Unmarshaling the openapi specification failed.")
	}

	config := restfulspec.Config{
		WebServices:    []*restful.WebService{ws}, // you control what services are visible
		WebServicesURL: "http://localhost:8080",
		APIPath:        "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}

	actual := restfulspec.BuildSwagger(config)

	if reflect.DeepEqual(expected, actual) != true {
		t.Errorf("Got:\n%v\nWanted:\n%v", asJSON(actual), asJSON(expected))
	}
}

func asJSON(v interface{}) string {
	data, _ := json.MarshalIndent(v, " ", " ")
	return string(data)
}

