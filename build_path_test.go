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
	r := ws.GET("/a").To(dummy).Writes(Sample{}).Build()

	p := buildPath(r)

	if got, want := p.Get.OperationProps.ID, "dummy"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
