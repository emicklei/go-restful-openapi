package restfulspec

import (
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

// nolint:paralleltest
func TestBuildDoc(t *testing.T) {
	path := "/testPath"

	ws1 := new(restful.WebService)
	ws1.Path(path)
	ws1.Route(ws1.GET("").To(dummy))

	ws2 := new(restful.WebService)
	ws2.Path(path)
	ws2.Route(ws2.DELETE("").To(dummy))

	c := Config{}
	c.WebServices = []*restful.WebService{ws1, ws2}
	doc := BuildDoc(c)

	if !(doc.Paths.Find(path).Get != nil && doc.Paths.Find(path).Delete != nil) {
		t.Errorf("Swagger spec should have methods for GET and DELETE")
	}

}

// nolint:paralleltest
func TestEnablingCORS(t *testing.T) {
	ws := NewOpenAPIService(Config{})
	wc := restful.NewContainer().Add(ws)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	origin := "somewhere.over.the.rainbow"
	request.Header[restful.HEADER_Origin] = []string{origin}

	wc.Dispatch(recorder, request)

	responseHeader := recorder.Result().Header.Get(restful.HEADER_AccessControlAllowOrigin)
	if responseHeader != origin {
		t.Errorf("The CORS header was set to the wrong value, expected: %s but got: %s", origin, responseHeader)
	}
}

// nolint:paralleltest
func TestDisablingCORS(t *testing.T) {
	ws := NewOpenAPIService(Config{DisableCORS: true})
	wc := restful.NewContainer().Add(ws)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	request.Header[restful.HEADER_Origin] = []string{"somewhere.over.the.rainbow"}

	wc.Dispatch(recorder, request)

	responseHeader := recorder.Result().Header.Get(restful.HEADER_AccessControlAllowOrigin)
	if responseHeader != "" {
		t.Errorf("The CORS header was set to %s but it was disabled so should not be set", responseHeader)
	}
}
