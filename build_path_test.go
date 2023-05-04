package restfulspec

import (
	"testing"

	"github.com/go-openapi/spec"

	"github.com/emicklei/go-restful/v3"
)

func TestRouteToPath(t *testing.T) {
	description := "get the <strong>a</strong> <em>b</em> test\nthis is the test description"
	notes := "notes\nblah blah"

	ws := new(restful.WebService)
	ws.Path("/tests/{v}")
	ws.Param(ws.PathParameter("v", "value of v").DefaultValue("default-v"))
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/a/{b}").To(dummy).
		Doc(description).
		Metadata(KeyOpenAPITags, []string{"tests"}).
		AddExtension("x-restful-test", "test value").
		Notes(notes).
		Param(ws.PathParameter("b", "value of b").DefaultValue("default-b")).
		Param(ws.QueryParameter("q", "value of q").DefaultValue("default-q")).
		Returns(200, "list of a b tests", []Sample{}).
		DefaultReturns("default", Sample{}).
		Writes([]Sample{}))
	ws.Route(ws.GET("/a/{b}/{c:[a-z]+}/{d:[1-9]+}/e/{f:*}").To(dummy).
		Param(ws.PathParameter("b", "value of b").DefaultValue("default-b")).
		Param(ws.PathParameter("c", "with regex").DefaultValue("abc")).
		Param(ws.PathParameter("d", "with regex").DefaultValue("abcef")).
		Param(ws.PathParameter("f", "with regex")).
		Param(ws.QueryParameter("q", "value of q").DataType("string").DataFormat("date").
			DefaultValue("default-q").AllowMultiple(true)).
		Returns(200, "list of a b tests", []Sample{}).
		Writes([]Sample{}))

	p := buildPaths(ws, Config{})
	t.Log(asJSON(p))

	if op := p.Paths["/tests/{v}/a/{b}"].Get; op.Parameters[0].Type != "string" {
		t.Error("Parameter type is not set.")
	} else if op.Extensions["x-restful-test"] != "test value" {
		t.Error("Extensions not set.")
	} else if len(op.Tags) != 1 || op.Tags[0] != "tests" {
		t.Error("Metadata/openapi.tags not set.")
	}
	if _, exists := p.Paths["/tests/{v}/a/{b}/{c}/{d}/e/{f}"]; !exists {
		t.Error("Expected path to exist after it was sanitized.")
	}

	q, exists := getParameter(p.Paths["/tests/{v}/a/{b}/{c}/{d}/e/{f}"], "q")
	if !exists {
		t.Errorf("get parameter q failed")
	}
	if q.Type != "array" || q.Items.Type != "string" || q.Format != "date" {
		t.Errorf("parameter q expected to be a date array")
	}

	if p.Paths["/tests/{v}/a/{b}"].Get.Description != notes {
		t.Errorf("GET description incorrect")
	}
	if p.Paths["/tests/{v}/a/{b}"].Get.Summary != "get the a b test\nthis is the test description" {
		t.Errorf("GET summary incorrect")
	}
	response := p.Paths["/tests/{v}/a/{b}"].Get.Responses.StatusCodeResponses[200]
	if response.Schema.Type[0] != "array" {
		t.Errorf("response type incorrect")
	}
	defaultResponse := p.Paths["/tests/{v}/a/{b}"].Get.Responses.Default
	if defaultResponse.Description != "default" {
		t.Errorf("default response description incorrect")
	}
	if response.Schema.Items.Schema.Ref.String() != "#/definitions/restfulspec.Sample" {
		t.Errorf("response element type incorrect")
	}

	// Test for patterns
	path := p.Paths["/tests/{v}/a/{b}/{c}/{d}/e/{f}"]
	checkPattern(t, path, "c", "[a-z]+")
	checkPattern(t, path, "d", "[1-9]+")
	checkPattern(t, path, "v", "")
	checkPattern(t, path, "f", ".*")
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
		t.Errorf("Expected Parameter %s to exist", paramName)
	}
	if param.Pattern != pattern {
		t.Errorf("Expected pattern %s to equal %s", param.Pattern, pattern)
	}
}

func TestGoogleCustomMethods(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests").Consumes(restful.MIME_JSON).Produces(restful.MIME_XML)
	ws.Route(ws.GET("/resource:validate").To(dummy).
		Doc("validate resource").
		Returns(200, "validate resource successfully", []Sample{}).
		Writes([]Sample{}))
	ws.Route(ws.POST("/resource/{resourceId}:init").To(dummy).
		Doc("init the specified resource").
		Returns(200, "init the specified resource successfully", []Sample{}).
		Returns(500, "internal server error", []Sample{}).
		Reads(Sample{}).
		Writes([]Sample{}))

	p := buildPaths(ws, Config{})
	t.Log(asJSON(p))

	path, exists := p.Paths["/tests/resource:validate"]
	if !exists {
		t.Error("Expected path to exist after it was sanitized.")
	}
	if path.Get.Summary != "validate resource" {
		t.Errorf("GET description incorrect")
	}
	response := path.Get.Responses.StatusCodeResponses[200]
	if response.Schema.Type[0] != "array" {
		t.Errorf("response type incorrect")
	}
	if response.Schema.Items.Schema.Ref.String() != "#/definitions/restfulspec.Sample" {
		t.Errorf("response element type incorrect")
	}
	path, exists = p.Paths["/tests/resource/{resourceId}:init"]
	if !exists {
		t.Error("Expected path to exist after it was sanitized.")
	}
	if path.Post.Summary != "init the specified resource" {
		t.Errorf("POST description incorrect")
	}
	if _, exists := p.Paths["/tests/resource/{resourceId}:init"].Post.Responses.StatusCodeResponses[500]; !exists {
		t.Errorf("Response code 500 not added to spec.")
	}
	response = path.Post.Responses.StatusCodeResponses[200]
	if response.Schema.Type[0] != "array" {
		t.Errorf("response type incorrect")
	}
	if response.Schema.Items.Schema.Ref.String() != "#/definitions/restfulspec.Sample" {
		t.Errorf("response element type incorrect")
	}
}

func TestRouteToPathForAllowableValues(t *testing.T) {
	description := "get the <strong>a</strong> <em>b</em> test\nthis is the test description"
	notes := "notes\nblah blah"

	allowedCheeses := map[string]string{"cheddar": "cheddar", "feta": "feta", "colby-jack": "colby-jack", "mozzerella": "mozzerella"}
	ws := new(restful.WebService)
	ws.Path("/tests/{v}")
	ws.Param(ws.PathParameter("v", "value of v").DefaultValue("default-v"))
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/a/{cheese}").To(dummy).
		Doc(description).
		Notes(notes).
		Param(ws.PathParameter("cheese", "value of cheese").DefaultValue("cheddar").AllowableValues(allowedCheeses)).
		Param(ws.QueryParameter("q", "value of q").DefaultValue("default-q")).
		Returns(200, "list of a b tests", []Sample{}).
		Writes([]Sample{}))

	p := buildPaths(ws, Config{})
	t.Log(asJSON(p))

	if p.Paths["/tests/{v}/a/{cheese}"].Get.Parameters[0].Type != "string" {
		t.Error("Parameter type is not set.")
	}

	path, exists := p.Paths["/tests/{v}/a/{cheese}"]
	if !exists {
		t.Error("Expected path to exist after it was sanitized.")
	}

	if path.Get.Description != notes {
		t.Errorf("GET description incorrect")
	}
	if path.Get.Summary != "get the a b test\nthis is the test description" {
		t.Errorf("GET summary incorrect")
	}
	if len(path.Get.Parameters) != 3 {
		t.Errorf("Path expected to have 3 parameters but had %d", len(path.Get.Parameters))
	}
	if path.Get.Parameters[1].Name != "cheese" || len(path.Get.Parameters[1].Enum) != len(allowedCheeses) {
		t.Errorf("Path parameter 'cheese' expected to have enum of allowable values of lenght %d but was of length %d", len(allowedCheeses), len(path.Get.Parameters[1].Enum))
	}
	for _, cheeseIface := range path.Get.Parameters[1].Enum {
		cheese := cheeseIface.(string)
		if _, found := allowedCheeses[cheese]; !found {
			t.Errorf("Path parameter 'cheese' had enum value of %s which was not found in the allowableValues map", cheese)
		}
	}

	response := path.Get.Responses.StatusCodeResponses[200]
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
		Reads(Sample{}).
		Writes([]Sample{}))

	p := buildPaths(ws, Config{})
	t.Log(asJSON(p))

	if p.Paths["/tests/a/a/b"].Get.Summary != "get a b test" {
		t.Errorf("GET summary incorrect")
	}
	if p.Paths["/tests/a/a/b"].Post.Summary != "post a b test" {
		t.Errorf("POST summary incorrect")
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

func TestReadArrayObjectInBody(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests/a")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)

	ws.Route(ws.POST("/a/b").To(dummy).
		Doc("post a b test with array in body").
		Returns(200, "list of a b tests", []Sample{}).
		Returns(500, "internal server error", []Sample{}).
		Reads([]Sample{}).
		Writes([]Sample{}))

	p := buildPaths(ws, Config{})
	t.Log(asJSON(p))

	postInfo := p.Paths["/tests/a/a/b"].Post

	if postInfo.Summary != "post a b test with array in body" {
		t.Errorf("POST description incorrect")
	}
	if _, exists := postInfo.Responses.StatusCodeResponses[500]; !exists {
		t.Errorf("Response code 500 not added to spec.")
	}
	// identify element model type in body array
	expectedItemRef := spec.MustCreateRef("#/definitions/restfulspec.Sample")
	postBody := postInfo.Parameters[0]
	if postBody.Schema.Ref.String() != "" {
		t.Errorf("you shouldn't have body Ref setting when using array in body!")
	}
	// check body array dy item ref
	postBodyitems := postBody.Schema.Items.Schema.Ref
	if postBodyitems.String() != expectedItemRef.String() {
		t.Errorf("Expected: %s, Got: %s", expectedItemRef.String(), expectedItemRef.String())
	}

	if postBody.Format != "" || postBody.Type != "" || postBody.Default != nil {
		t.Errorf("Invalid parameter property is set on body property")
	}
}

// TestWritesPrimitive ensures that if an operation returns a primitive, then it
// is used as such (and not a ref to a definition).
func TestWritesPrimitive(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests/returns")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/primitive").To(dummy).
		Doc("post that returns a string").
		Returns(200, "primitive string", "(this is a string)").
		Writes("(this is a string)"))

	ws.Route(ws.POST("/custom").To(dummy).
		Doc("post that returns a custom structure").
		Returns(200, "sample object", Sample{}).
		Writes(Sample{}))

	p := buildPaths(ws, Config{})
	t.Log(asJSON(p))

	// Make sure that the operation that returns a primitive type is correct.
	if pathInfo, okay := p.Paths["/tests/returns/primitive"]; !okay {
		t.Errorf("Could not find path")
	} else {
		postInfo := pathInfo.Post

		if postInfo.Summary != "post that returns a string" {
			t.Errorf("POST description incorrect")
		}
		response := postInfo.Responses.StatusCodeResponses[200]
		if response.Schema.Ref.String() != "" {
			t.Errorf("Expected no ref; got: %s", response.Schema.Ref.String())
		}
		if len(response.Schema.Type) != 1 {
			t.Errorf("Expected exactly one type; got: %d", len(response.Schema.Type))
		}
		if response.Schema.Type[0] != "string" {
			t.Errorf("Expected a type of string; got: %s", response.Schema.Type[0])
		}
	}

	// Make sure that the operation that returns a custom type is correct.
	if pathInfo, okay := p.Paths["/tests/returns/custom"]; !okay {
		t.Errorf("Could not find path")
	} else {
		postInfo := pathInfo.Post

		if postInfo.Summary != "post that returns a custom structure" {
			t.Errorf("POST description incorrect")
		}
		response := postInfo.Responses.StatusCodeResponses[200]
		if response.Schema.Ref.String() != "#/definitions/restfulspec.Sample" {
			t.Errorf("Expected ref '#/definitions/restfulspec.Sample'; got: %s", response.Schema.Ref.String())
		}
		if len(response.Schema.Type) != 0 {
			t.Errorf("Expected exactly zero types; got: %d", len(response.Schema.Type))
		}
	}
}
