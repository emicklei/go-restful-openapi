package restfulspec

import (
	"testing"

	restful "github.com/emicklei/go-restful"
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

	p := buildPaths(ws, Config{})
	t.Log(asJSON(p))

	if p.Paths["/tests/{v}/a/{b}"].Get.Parameters[0].Type != "string" {
		t.Error("Parameter type is not set.")

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
}
