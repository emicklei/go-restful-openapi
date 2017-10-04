package restfulspec

import (
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

func TestRouteToPath(t *testing.T) {
	description := "get the <strong>a</strong> <em>b</em> test\nthis is the test description"

	ws := new(restful.WebService)
	ws.Path("/tests/{v}")
	ws.Param(ws.PathParameter("v", "value of v").DefaultValue("default-v"))
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/a/{b}").To(dummy).
		Doc(description).
		Param(ws.PathParameter("b", "value of b").DefaultValue("default-b")).
		Param(ws.QueryParameter("q", "value of q").DefaultValue("default-q")).
		Returns(200, "list of a b tests", []Sample{}).
		Writes([]Sample{}))
	ws.Route(ws.GET("/a/{b}/{c:[a-z]+}/{d:[1-9]+}/e").To(dummy).
		Doc("get the a b test").
		Param(ws.PathParameter("b", "value of b").DefaultValue("default-b")).
		Param(ws.PathParameter("c", "with regex").DefaultValue("abc")).
		Param(ws.PathParameter("d", "with regex").DefaultValue("abcef")).
		Param(ws.QueryParameter("q", "value of q").DefaultValue("default-q")).
		Returns(200, "list of a b tests", []Sample{}).
		Writes([]Sample{}))

	p := buildPaths(ws, Config{})
	t.Log(asJSON(p))

	if p.Paths["/tests/{v}/a/{b}"].Get.Parameters[0].Type != "string" {
		t.Error("Parameter type is not set.")
	}
	if _, exists := p.Paths["/tests/{v}/a/{b}/{c}/{d}/e"]; !exists {
		t.Error("Expected path to exist after it was sanitized.")
	}

	if p.Paths["/tests/{v}/a/{b}"].Get.Description != description {
		t.Errorf("GET description incorrect")
	}
	if p.Paths["/tests/{v}/a/{b}"].Get.Summary != "get the a b test" {
		t.Errorf("GET summary incorrect")
	}
	response := p.Paths["/tests/{v}/a/{b}"].Get.Responses.StatusCodeResponses[200]
	if response.Schema.Type[0] != "array" {
		t.Errorf("response type incorrect")
	}
	if response.Schema.Items.Schema.Ref.String() != "#/definitions/restfulspec.Sample" {
		t.Errorf("response element type incorrect")
	}

	// Test for patterns
	path := p.Paths["/tests/{v}/a/{b}/{c}/{d}/e"]
	checkPattern(t, path, "c", "[a-z]+")
	checkPattern(t, path, "d", "[1-9]+")
	checkPattern(t, path, "v", "")
}

func getParameter(path spec.PathItem, name string) (*spec.Parameter, bool) {
	for _, param := range path.Get.Parameters {
		if param.Name == name {
			return &param, true
		}
	}
	return nil, false
}

func checkPattern(t *testing.T, path spec.PathItem, paramName string, pattern string) {
	param, exists := getParameter(path, paramName)
	if !exists {
		t.Error("Expected Parameter %s to exist", paramName)
	}
	if param.Pattern != pattern {
		t.Error("Expected pattern %s to equal %s", param.Pattern, pattern)
	}
}

func TestMultipleMethodsRouteToPath(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests/a")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/a/b").To(dummy).
		Doc("get a b test").
		Returns(200, "list of a b tests", []Sample{}).
		Writes([]Sample{}))
	ws.Route(ws.POST("/a/b").To(dummy).
		Doc("post a b test").
		Returns(200, "list of a b tests", []Sample{}).
		Returns(500, "internal server error", []Sample{}).
		Reads(Sample{}).
		Writes([]Sample{}))

	p := buildPaths(ws, Config{})
	t.Log(asJSON(p))

	if p.Paths["/tests/a/a/b"].Get.Description != "get a b test" {
		t.Errorf("GET description incorrect")
	}
	if p.Paths["/tests/a/a/b"].Post.Description != "post a b test" {
		t.Errorf("POST description incorrect")
	}
	if _, exists := p.Paths["/tests/a/a/b"].Post.Responses.StatusCodeResponses[500]; !exists {
		t.Errorf("Response code 500 not added to spec.")
	}

	expectedRef := spec.MustCreateRef("#/definitions/restfulspec.Sample")
	postBodyparam := p.Paths["/tests/a/a/b"].Post.Parameters[0]
	postBodyRef := postBodyparam.Schema.Ref
	if postBodyRef.String() != expectedRef.String() {
		t.Errorf("Expected: %s, Got: %s", expectedRef.String(), postBodyRef.String())
	}

	if postBodyparam.Format != "" || postBodyparam.Type != "" || postBodyparam.Default != nil {
		t.Errorf("Invalid parameter property is set on body property")
	}
}
