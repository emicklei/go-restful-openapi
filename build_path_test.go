package restfulspec

import (
	"testing"

	restful "github.com/emicklei/go-restful"
)

func TestRouteToPath(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	r := ws.GET("/a/{b}").To(dummy).
		Param(ws.PathParameter("b", "value of b")).
		Param(ws.QueryParameter("q", "value of q")).
		Writes(Sample{}).
		Build()

	p := buildPath(r)
	t.Log(asJSON(p))
}
