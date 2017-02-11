package restfulspec

import (
	"reflect"

	restful "github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

// MapSchemaFormatFunc can be used to modify typeName at definition time.
type MapSchemaFormatFunc func(typeName string) string

// MapModelTypeNameFunc can be used to return the desired typeName for a given
// type. It will return false if the default name should be used.
type MapModelTypeNameFunc func(t reflect.Type) (string, bool)

// Config holds service api metadata.
type Config struct {
	// WebServicesURL is the url where the services are available, e.g. http://localhost:8080
	// if left empty then the basePath of Swagger is taken from the actual request
	WebServicesURL string
	// APIPath is the path where the JSON api is avaiable , e.g. /apidocs.json
	APIPath string
	// api listing is constructed from this list of restful WebServices.
	WebServices []*restful.WebService
	// [optional] on default CORS (Cross-Origin-Resource-Sharing) is enabled.
	DisableCORS bool
	// Top-level API version. Is reflected in the resource listing.
	APIVersion string
	// OpenAPI global info struct
	Info spec.Info
	// [optional] If set, model builder should call this handler to get addition typename-to-swagger-format-field conversion.
	SchemaFormatHandler MapSchemaFormatFunc
	// [optional] If set, model builder should call this handler to retrieve the name for a given type.
	ModelTypeNameHandler MapModelTypeNameFunc
}
