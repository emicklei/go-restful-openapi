package restfulspec

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

// nolint:paralleltest
func TestBuildSwagger(t *testing.T) {
	path := "/testPath"

	ws1 := new(restful.WebService)
	ws1.Path(path)
	ws1.Route(ws1.GET("").To(dummy))

	ws2 := new(restful.WebService)
	ws2.Path(path)
	ws2.Route(ws2.DELETE("").To(dummy))

	c := Config{}
	c.WebServices = []*restful.WebService{ws1, ws2}
	s := BuildSwagger(c)

	if s.Paths.Paths[path].Get == nil || s.Paths.Paths[path].Delete == nil {
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

// nolint:paralleltest
func TestNewOpenAPIServiceWithOAS2(t *testing.T) {
	ws1 := new(restful.WebService)
	ws1.Path("/api")
	ws1.Route(ws1.GET("/users").To(dummy).Doc("Get users"))

	config := Config{
		WebServices: []*restful.WebService{ws1},
		APIPath:     "/apidocs.json",
		OASVersion:  OASVersion20,
	}
	ws := NewOpenAPIService(config)
	wc := restful.NewContainer().Add(ws)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/apidocs.json/", nil)

	wc.Dispatch(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("Expected status 200 but got %d", recorder.Code)
	}

	// Verify it's a valid OAS 2.0 document
	var doc map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &doc); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if doc["swagger"] != "2.0" {
		t.Errorf("Expected swagger version 2.0 but got %v", doc["swagger"])
	}
}

// nolint:paralleltest
func TestNewOpenAPIServiceWithOAS3(t *testing.T) {
	ws1 := new(restful.WebService)
	ws1.Path("/api")
	ws1.Route(ws1.GET("/users").To(dummy).Doc("Get users"))

	config := Config{
		WebServices: []*restful.WebService{ws1},
		APIPath:     "/apidocs.json",
		OASVersion:  OASVersion320,
	}
	ws := NewOpenAPIService(config)
	wc := restful.NewContainer().Add(ws)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/apidocs.json/", nil)

	wc.Dispatch(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("Expected status 200 but got %d", recorder.Code)
	}

	// Verify it's a valid OAS 3.x document
	var doc map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &doc); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if doc["openapi"] == nil {
		t.Error("Expected openapi field to be present for OAS 3.x document")
	}
	openapi, ok := doc["openapi"].(string)
	if !ok || len(openapi) < 3 || openapi[:2] != "3." {
		t.Errorf("Expected openapi version 3.x but got %v", doc["openapi"])
	}
}
