package restfulspec

import (
	"net/http"
	"testing"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/erraggy/oastools/parser"
)

// Test types for schema generation
type TestUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name" description:"User's full name"`
	Email    string `json:"email" format:"email"`
	Age      int    `json:"age,omitempty" minimum:"0" maximum:"150"`
	Role     string `json:"role" enum:"admin|user|guest"`
	IsActive bool   `json:"is_active"`
}

type TestUserWithOASTags struct {
	ID    int    `json:"id" oas:"description=User ID,minimum=1"`
	Name  string `json:"name" oas:"description=User name,minLength=1,maxLength=100"`
	Email string `json:"email" oas:"format=email"`
}

func TestBuildOAS2_Basic(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users")
	ws.Route(ws.GET("").To(dummyHandler).
		Doc("List all users").
		Operation("listUsers").
		Writes([]TestUser{}).
		Returns(http.StatusOK, "OK", []TestUser{}))

	ws.Route(ws.POST("").To(dummyHandler).
		Doc("Create a user").
		Operation("createUser").
		Reads(TestUser{}).
		Returns(http.StatusCreated, "Created", TestUser{}))

	ws.Route(ws.GET("/{id}").To(dummyHandler).
		Doc("Get a user by ID").
		Operation("getUser").
		Param(ws.PathParameter("id", "User ID").DataType("integer")).
		Writes(TestUser{}).
		Returns(http.StatusOK, "OK", TestUser{}).
		Returns(http.StatusNotFound, "Not found", nil))

	config := Config{
		WebServices: []*restful.WebService{ws},
		APIVersion:  "1.0.0",
		Host:        "api.example.com",
		Schemes:     []string{"https"},
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	// Verify basic structure
	if doc.Swagger != "2.0" {
		t.Errorf("Expected swagger version 2.0, got %s", doc.Swagger)
	}
	if doc.Host != "api.example.com" {
		t.Errorf("Expected host api.example.com, got %s", doc.Host)
	}
	if len(doc.Schemes) != 1 || doc.Schemes[0] != "https" {
		t.Errorf("Expected schemes [https], got %v", doc.Schemes)
	}

	// Verify paths exist
	if doc.Paths == nil {
		t.Fatal("Expected paths to be non-nil")
	}
	if _, ok := doc.Paths["/users"]; !ok {
		t.Error("Expected /users path")
	}
	if _, ok := doc.Paths["/users/{id}"]; !ok {
		t.Error("Expected /users/{id} path")
	}

	// Verify GET /users operation
	usersPath := doc.Paths["/users"]
	if usersPath.Get == nil {
		t.Error("Expected GET operation on /users")
	} else {
		if usersPath.Get.OperationID != "listUsers" {
			t.Errorf("Expected operation ID listUsers, got %s", usersPath.Get.OperationID)
		}
	}

	// Verify POST /users operation
	if usersPath.Post == nil {
		t.Error("Expected POST operation on /users")
	} else {
		if usersPath.Post.OperationID != "createUser" {
			t.Errorf("Expected operation ID createUser, got %s", usersPath.Post.OperationID)
		}
	}
}

func TestBuildOAS3_Basic(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users")
	ws.Route(ws.GET("").To(dummyHandler).
		Doc("List all users").
		Operation("listUsers").
		Writes([]TestUser{}).
		Returns(http.StatusOK, "OK", []TestUser{}))

	config := Config{
		WebServices: []*restful.WebService{ws},
		APIVersion:  "1.0.0",
		OASVersion:  OASVersion320,
	}

	doc, err := BuildOAS3(config)
	if err != nil {
		t.Fatalf("BuildOAS3 failed: %v", err)
	}

	// Verify version
	if doc.OpenAPI != "3.2.0" {
		t.Errorf("Expected OpenAPI version 3.2.0, got %s", doc.OpenAPI)
	}

	// Verify paths
	if doc.Paths == nil {
		t.Fatal("Expected paths to be non-nil")
	}
	if _, ok := doc.Paths["/users"]; !ok {
		t.Error("Expected /users path")
	}
}

func TestBuildOAS3_WithServers(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")
	ws.Route(ws.GET("/health").To(dummyHandler).
		Operation("healthCheckServers").
		Doc("Health check"))

	config := Config{
		WebServices: []*restful.WebService{ws},
		Servers: []*parser.Server{
			{URL: "https://api.example.com"},
			{URL: "https://staging.example.com"},
		},
	}

	doc, err := BuildOAS3(config)
	if err != nil {
		t.Fatalf("BuildOAS3 failed: %v", err)
	}

	if len(doc.Servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(doc.Servers))
	}
}

func TestBuildOAS3_HostToServer(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")
	ws.Route(ws.GET("/health").To(dummyHandler).
		Operation("healthCheckHost"))

	config := Config{
		WebServices: []*restful.WebService{ws},
		Host:        "api.example.com",
		Schemes:     []string{"https"},
	}

	doc, err := BuildOAS3(config)
	if err != nil {
		t.Fatalf("BuildOAS3 failed: %v", err)
	}

	if len(doc.Servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(doc.Servers))
	}
	if doc.Servers[0].URL != "https://api.example.com" {
		t.Errorf("Expected server URL https://api.example.com, got %s", doc.Servers[0].URL)
	}
}

func TestBuildOAS2_WithSchemaNaming(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users")
	ws.Route(ws.GET("").To(dummyHandler).
		Writes(TestUser{}).
		Returns(http.StatusOK, "OK", TestUser{}))

	config := Config{
		WebServices:  []*restful.WebService{ws},
		SchemaNaming: SchemaNamingTypeOnly, // Just type name without package
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	// Verify definitions exist (schema names depend on naming strategy)
	if doc.Definitions == nil {
		t.Fatal("Expected definitions to be non-nil")
	}
	// With SchemaNamingTypeOnly, we should have "TestUser" not "restfulspec.TestUser"
	if _, ok := doc.Definitions["TestUser"]; !ok {
		t.Logf("Available definitions: %v", keysOf(doc.Definitions))
		t.Error("Expected TestUser definition with SchemaNamingTypeOnly")
	}
}

func TestBuildOAS2_PostBuildHandler(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")
	ws.Route(ws.GET("/health").To(dummyHandler).
		Operation("healthCheckPostBuild2"))

	handlerCalled := false
	config := Config{
		WebServices: []*restful.WebService{ws},
		PostBuildOAS2Handler: func(doc *parser.OAS2Document) {
			handlerCalled = true
			doc.Info = &parser.Info{
				Title:   "Modified API",
				Version: "2.0.0",
			}
		},
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	if !handlerCalled {
		t.Error("Expected PostBuildOAS2Handler to be called")
	}
	if doc.Info == nil || doc.Info.Title != "Modified API" {
		t.Error("Expected Info to be modified by handler")
	}
}

func TestBuildOAS3_PostBuildHandler(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")
	ws.Route(ws.GET("/health").To(dummyHandler).
		Operation("healthCheckPostBuild3"))

	handlerCalled := false
	config := Config{
		WebServices: []*restful.WebService{ws},
		PostBuildOAS3Handler: func(doc *parser.OAS3Document) {
			handlerCalled = true
			doc.Info = &parser.Info{
				Title:   "Modified API",
				Version: "2.0.0",
			}
		},
	}

	doc, err := BuildOAS3(config)
	if err != nil {
		t.Fatalf("BuildOAS3 failed: %v", err)
	}

	if !handlerCalled {
		t.Error("Expected PostBuildOAS3Handler to be called")
	}
	if doc.Info == nil || doc.Info.Title != "Modified API" {
		t.Error("Expected Info to be modified by handler")
	}
}

func TestBuildOAS2_Parameters(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/items")
	ws.Route(ws.GET("").To(dummyHandler).
		Doc("List items").
		Operation("listItems").
		Param(ws.QueryParameter("limit", "Max items to return").DataType("integer").DefaultValue("10")).
		Param(ws.QueryParameter("offset", "Offset for pagination").DataType("integer")).
		Returns(http.StatusOK, "OK", nil))

	ws.Route(ws.GET("/{id}").To(dummyHandler).
		Doc("Get item").
		Operation("getItem").
		Param(ws.PathParameter("id", "Item ID").DataType("integer")).
		Returns(http.StatusOK, "OK", nil))

	config := Config{
		WebServices: []*restful.WebService{ws},
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	// Check GET /items parameters
	listOp := doc.Paths["/items"].Get
	if listOp == nil {
		t.Fatal("Expected GET operation on /items")
	}
	if len(listOp.Parameters) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(listOp.Parameters))
	}

	// Check GET /items/{id} path parameter
	getOp := doc.Paths["/items/{id}"].Get
	if getOp == nil {
		t.Fatal("Expected GET operation on /items/{id}")
	}
	if len(getOp.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(getOp.Parameters))
	}
	if getOp.Parameters[0].In != "path" {
		t.Errorf("Expected path parameter, got %s", getOp.Parameters[0].In)
	}
}

func TestBuildOAS2_Tags(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users")
	ws.Route(ws.GET("").To(dummyHandler).
		Doc("List users").
		Operation("listUsersWithTags").
		Metadata(KeyOpenAPITags, []string{"users", "public"}).
		Returns(http.StatusOK, "OK", nil))

	config := Config{
		WebServices: []*restful.WebService{ws},
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	listOp := doc.Paths["/users"].Get
	if listOp == nil {
		t.Fatal("Expected GET operation")
	}
	if len(listOp.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(listOp.Tags))
	}
}

func TestBuildOAS2_Deprecated(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/old")
	ws.Route(ws.GET("").To(dummyHandler).
		Doc("Old endpoint").
		Operation("deprecatedEndpoint").
		Deprecate().
		Returns(http.StatusOK, "OK", nil))

	config := Config{
		WebServices: []*restful.WebService{ws},
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	op := doc.Paths["/old"].Get
	if op == nil {
		t.Fatal("Expected GET operation")
	}
	if !op.Deprecated {
		t.Error("Expected operation to be deprecated")
	}
}

func TestBuildOAS3_DefaultVersion(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")
	ws.Route(ws.GET("/health").To(dummyHandler).
		Operation("healthCheckDefaultVersion"))

	// Don't set OASVersion, should default to 3.2.0
	config := Config{
		WebServices: []*restful.WebService{ws},
	}

	doc, err := BuildOAS3(config)
	if err != nil {
		t.Fatalf("BuildOAS3 failed: %v", err)
	}

	if doc.OpenAPI != "3.2.0" {
		t.Errorf("Expected default OpenAPI version 3.2.0, got %s", doc.OpenAPI)
	}
}

func TestBuildOAS3_SpecificVersion(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")
	ws.Route(ws.GET("/health").To(dummyHandler).
		Operation("healthCheckSpecificVersion"))

	config := Config{
		WebServices: []*restful.WebService{ws},
		OASVersion:  OASVersion310,
	}

	doc, err := BuildOAS3(config)
	if err != nil {
		t.Fatalf("BuildOAS3 failed: %v", err)
	}

	if doc.OpenAPI != "3.1.0" {
		t.Errorf("Expected OpenAPI version 3.1.0, got %s", doc.OpenAPI)
	}
}

// Helper function
func dummyHandler(req *restful.Request, resp *restful.Response) {}

func keysOf(m map[string]*parser.Schema) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestBuildOAS3_WithOASVersion20_ReturnsError(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")
	ws.Route(ws.GET("/health").To(dummyHandler))

	config := Config{
		WebServices: []*restful.WebService{ws},
		OASVersion:  OASVersion20, // Explicitly set OAS 2.0
	}

	_, err := BuildOAS3(config)
	if err == nil {
		t.Error("Expected error when calling BuildOAS3 with OASVersion20")
	}
}

func TestBuildOAS2_DefaultResponse(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")

	type ErrorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	ws.Route(ws.GET("/test").To(dummyHandler).
		Operation("testDefaultResponse").
		DefaultReturns("Unexpected error", ErrorResponse{}))

	config := Config{
		WebServices: []*restful.WebService{ws},
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	op := doc.Paths["/api/test"].Get
	if op == nil {
		t.Fatal("Expected GET operation")
	}

	// Check that default response exists
	if op.Responses == nil || op.Responses.Default == nil {
		t.Fatal("Expected default response")
	}

	if op.Responses.Default.Description != "Unexpected error" {
		t.Errorf("Expected description 'Unexpected error', got '%s'", op.Responses.Default.Description)
	}
}

func TestBuildOAS2_HeaderParameter(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")
	ws.Route(ws.GET("/test").To(dummyHandler).
		Operation("testHeaderParam").
		Param(ws.HeaderParameter("X-Request-ID", "Request tracking ID").DataType("string").Required(true)).
		Returns(http.StatusOK, "OK", nil))

	config := Config{
		WebServices: []*restful.WebService{ws},
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	op := doc.Paths["/api/test"].Get
	if op == nil {
		t.Fatal("Expected GET operation")
	}

	if len(op.Parameters) != 1 {
		t.Fatalf("Expected 1 parameter, got %d", len(op.Parameters))
	}

	param := op.Parameters[0]
	if param.In != "header" {
		t.Errorf("Expected header parameter, got %s", param.In)
	}
	if param.Name != "X-Request-ID" {
		t.Errorf("Expected param name 'X-Request-ID', got '%s'", param.Name)
	}
	if !param.Required {
		t.Error("Expected parameter to be required")
	}
}

func TestBuildOAS2_FormParameter(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")
	ws.Route(ws.POST("/upload").To(dummyHandler).
		Operation("testFormParam").
		Consumes("application/x-www-form-urlencoded").
		Param(ws.FormParameter("filename", "File name").DataType("string")).
		Returns(http.StatusOK, "OK", nil))

	config := Config{
		WebServices: []*restful.WebService{ws},
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	op := doc.Paths["/api/upload"].Post
	if op == nil {
		t.Fatal("Expected POST operation")
	}

	if len(op.Parameters) != 1 {
		t.Fatalf("Expected 1 parameter, got %d", len(op.Parameters))
	}

	param := op.Parameters[0]
	if param.In != "formData" {
		t.Errorf("Expected formData parameter, got %s", param.In)
	}
}

func TestGetTypeForDataType_AllTypes(t *testing.T) {
	testCases := []struct {
		dataType string
		expected any
	}{
		{"string", ""},
		{"integer", int(0)},
		{"int", int(0)},
		{"int32", int32(0)},
		{"int64", int64(0)},
		{"number", float64(0)},
		{"float64", float64(0)},
		{"float32", float32(0)},
		{"boolean", false},
		{"bool", false},
		{"file", nil},
		{"", ""}, // Empty defaults to string silently
	}

	for _, tc := range testCases {
		t.Run(tc.dataType, func(t *testing.T) {
			result := getTypeForDataType(tc.dataType)
			if result != tc.expected {
				t.Errorf("getTypeForDataType(%q) = %T(%v), expected %T(%v)",
					tc.dataType, result, result, tc.expected, tc.expected)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	t.Run("empty WebServices returns error", func(t *testing.T) {
		config := Config{}
		err := config.Validate()
		if err == nil {
			t.Error("Expected error for empty WebServices")
		}
	})

	t.Run("nil WebService in slice returns error", func(t *testing.T) {
		config := Config{
			WebServices: []*restful.WebService{nil},
		}
		err := config.Validate()
		if err == nil {
			t.Error("Expected error for nil WebService")
		}
	})

	t.Run("valid config returns nil", func(t *testing.T) {
		ws := new(restful.WebService)
		ws.Path("/api")
		config := Config{
			WebServices: []*restful.WebService{ws},
		}
		err := config.Validate()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestBuildOAS2_ParameterWithConstraints(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api")

	minVal := float64(1)
	maxVal := float64(100)
	ws.Route(ws.GET("/test").To(dummyHandler).
		Operation("testParamConstraints").
		Param(ws.QueryParameter("limit", "Limit results").
			DataType("integer").
			Minimum(minVal).
			Maximum(maxVal).
			DefaultValue("10")).
		Returns(http.StatusOK, "OK", nil))

	config := Config{
		WebServices: []*restful.WebService{ws},
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	op := doc.Paths["/api/test"].Get
	if op == nil {
		t.Fatal("Expected GET operation")
	}

	if len(op.Parameters) != 1 {
		t.Fatalf("Expected 1 parameter, got %d", len(op.Parameters))
	}

	param := op.Parameters[0]
	if param.Minimum == nil || *param.Minimum != 1 {
		t.Errorf("Expected minimum 1, got %v", param.Minimum)
	}
	if param.Maximum == nil || *param.Maximum != 100 {
		t.Errorf("Expected maximum 100, got %v", param.Maximum)
	}
}

func TestBuildOAS2_ParameterWithFormat(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users")
	ws.Route(ws.GET("/{user_id}").To(dummyHandler).
		Operation("getUserWithFormat").
		Param(ws.PathParameter("user_id", "User ID in UUID format").
			DataType("string").
			DataFormat("uuid")).
		Param(ws.QueryParameter("created_after", "Filter by creation date").
			DataType("string").
			DataFormat("date")).
		Param(ws.QueryParameter("version", "API version").
			DataType("integer").
			DataFormat("int64")).
		Returns(http.StatusOK, "OK", nil))

	config := Config{
		WebServices: []*restful.WebService{ws},
	}

	doc, err := BuildOAS2(config)
	if err != nil {
		t.Fatalf("BuildOAS2 failed: %v", err)
	}

	op := doc.Paths["/users/{user_id}"].Get
	if op == nil {
		t.Fatal("Expected GET operation")
	}

	if len(op.Parameters) != 3 {
		t.Fatalf("Expected 3 parameters, got %d", len(op.Parameters))
	}

	// Find parameters by name
	params := make(map[string]*parser.Parameter)
	for _, p := range op.Parameters {
		params[p.Name] = p
	}

	// Verify path parameter with uuid format
	userIDParam := params["user_id"]
	if userIDParam == nil {
		t.Fatal("Expected user_id parameter")
	}
	if userIDParam.Format != "uuid" {
		t.Errorf("Expected user_id format 'uuid', got '%s'", userIDParam.Format)
	}
	if userIDParam.In != "path" {
		t.Errorf("Expected user_id in 'path', got '%s'", userIDParam.In)
	}

	// Verify query parameter with date format
	createdAfterParam := params["created_after"]
	if createdAfterParam == nil {
		t.Fatal("Expected created_after parameter")
	}
	if createdAfterParam.Format != "date" {
		t.Errorf("Expected created_after format 'date', got '%s'", createdAfterParam.Format)
	}

	// Verify query parameter with int64 format
	versionParam := params["version"]
	if versionParam == nil {
		t.Fatal("Expected version parameter")
	}
	if versionParam.Format != "int64" {
		t.Errorf("Expected version format 'int64', got '%s'", versionParam.Format)
	}
}

func TestBuildOAS3_ParameterWithFormat(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/items")
	ws.Route(ws.GET("/{item_id}").To(dummyHandler).
		Operation("getItemWithFormat").
		Param(ws.PathParameter("item_id", "Item ID in UUID format").
			DataType("string").
			DataFormat("uuid")).
		Param(ws.QueryParameter("updated_after", "Filter by update date").
			DataType("string").
			DataFormat("date-time")).
		Returns(http.StatusOK, "OK", nil))

	config := Config{
		WebServices: []*restful.WebService{ws},
		OASVersion:  OASVersion310,
	}

	doc, err := BuildOAS3(config)
	if err != nil {
		t.Fatalf("BuildOAS3 failed: %v", err)
	}

	pathItem := doc.Paths["/items/{item_id}"]
	if pathItem == nil {
		t.Fatal("Expected /items/{item_id} path")
	}

	op := pathItem.Get
	if op == nil {
		t.Fatal("Expected GET operation")
	}

	if len(op.Parameters) != 2 {
		t.Fatalf("Expected 2 parameters, got %d", len(op.Parameters))
	}

	// Find parameters by name (OAS3 uses Schema for format)
	params := make(map[string]*parser.Parameter)
	for _, p := range op.Parameters {
		params[p.Name] = p
	}

	// Verify path parameter with uuid format in OAS3 (format is in Schema)
	itemIDParam := params["item_id"]
	if itemIDParam == nil {
		t.Fatal("Expected item_id parameter")
	}
	if itemIDParam.Schema == nil {
		t.Fatal("Expected item_id schema")
	}
	if itemIDParam.Schema.Format != "uuid" {
		t.Errorf("Expected item_id format 'uuid', got '%s'", itemIDParam.Schema.Format)
	}

	// Verify query parameter with date-time format in OAS3
	updatedAfterParam := params["updated_after"]
	if updatedAfterParam == nil {
		t.Fatal("Expected updated_after parameter")
	}
	if updatedAfterParam.Schema == nil {
		t.Fatal("Expected updated_after schema")
	}
	if updatedAfterParam.Schema.Format != "date-time" {
		t.Errorf("Expected updated_after format 'date-time', got '%s'", updatedAfterParam.Schema.Format)
	}
}
