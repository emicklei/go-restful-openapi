package restfulspec

import "github.com/go-openapi/spec"

// Documented is
type Documented interface {
	SwaggerDoc() map[string]string
}

type PostBuildSwaggerSchema interface {
	PostBuildSwaggerSchemaHandler(sm *spec.Schema)
}

const (
	// KeyOpenAPITags is a Metadata key for a restful Route
	KeyOpenAPITags = "openapi.tags"

	// ExtensionPrefix is the only prefix accepted for VendorExtensible extension keys
	ExtensionPrefix = "x-"
)

// SchemaType is used to wrap any raw types
// For example, to return a "schema": "file" one can use
// Returns(http.StatusOK, http.StatusText(http.StatusOK), SchemaType{RawType: "file"})
type SchemaType struct {
	RawType string
	Format  string
}
