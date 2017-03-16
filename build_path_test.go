package restfulspec

import (
	"testing"

	restful "github.com/emicklei/go-restful"
)

func TestRouteToPath(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests/{v}")
	ws.Param(ws.PathParameter("v", "value of v").DefaultValue("default-v"))
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/a/{b}").To(dummy).
		Doc("get the a b test").
		Param(ws.PathParameter("b", "value of b").DefaultValue("default-b")).
		Param(ws.QueryParameter("q", "value of q").DefaultValue("default-q")).
		Returns(200, "list of a b tests", []Sample{}).
		Writes([]Sample{}))

	p := BuildPaths(ws)
	t.Log(asJSON(p))
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
		Writes([]Sample{}))

	p := BuildPaths(ws)
	t.Log(asJSON(p))

	if p.Paths["/tests/a/a/b"].Get.Description != "get a b test" {
		t.Errorf("GET description incorrect")
	}
	if p.Paths["/tests/a/a/b"].Post.Description != "post a b test" {
		t.Errorf("POST description incorrect")
	}
}
