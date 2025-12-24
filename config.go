package restfulspec

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/emicklei/go-restful/v3"
	"github.com/erraggy/oastools/builder"
	"github.com/erraggy/oastools/parser"
	"github.com/go-openapi/spec"
)

// MapSchemaFormatFunc can be used to modify typeName at definition time.
// To use it set the SchemaFormatHandler in the config.
type MapSchemaFormatFunc func(typeName string) string

// MapModelTypeNameFunc can be used to return the desired typeName for a given
// type. It will return false if the default name should be used.
// To use it set the ModelTypeNameHandler in the config.
type MapModelTypeNameFunc func(t reflect.Type) (string, bool)

// PostBuildSwaggerObjectFunc can be used to change the creates Swagger Object
// before serving it. To use it set the PostBuildSwaggerObjectHandler in the config.
type PostBuildSwaggerObjectFunc func(s *spec.Swagger)

// DefinitionNameHandlerFunc generate name by this handler for definition without json tag.
// example: (for more, see file definition_name_test.go)
//
//	  field	      			 definition_name
//	  Name `json:"name"`  ->  name
//		 Name                ->  Name
//
// there are some example provided for use
//
//	DefaultNameHandler         GoRestfulDefinition -> GoRestfulDefinition (not changed)
//	LowerSnakeCasedNameHandler  GoRestfulDefinition -> go_restful_definition
//	LowerCamelCasedNameHandler  GoRestfulDefinition -> goRestfulDefinition
//	GoLowerCamelCasedNameHandler HTTPRestfulDefinition -> httpRestfulDefinition
type DefinitionNameHandlerFunc func(string) string

// Config holds service api metadata.
type Config struct {
	// [optional] If set then set this field with the generated Swagger Object
	Host string
	// [optional] If set then set this field with the generated Swagger Object
	Schemes []string
	// WebServicesURL is a DEPRECATED field; it never had any effect in this package.
	WebServicesURL string
	// APIPath is the path where the JSON api is available, e.g. /apidocs.json
	APIPath string
	// api listing is constructed from this list of restful WebServices.
	WebServices []*restful.WebService
	// [optional] on default CORS (Cross-Origin-Resource-Sharing) is enabled.
	DisableCORS bool
	// Top-level API version. Is reflected in the resource listing.
	APIVersion string
	// [optional] If set, model builder should call this handler to get addition typename-to-swagger-format-field conversion.
	SchemaFormatHandler MapSchemaFormatFunc
	// [optional] If set, model builder should call this handler to retrieve the name for a given type.
	ModelTypeNameHandler MapModelTypeNameFunc
	// [optional] If set then call this function with the generated Swagger Object
	PostBuildSwaggerObjectHandler PostBuildSwaggerObjectFunc
	// [optional] If set then call handler's function for to generate name by this handler for definition without json tag,
	//   you can use you DefinitionNameHandler, also, there are four DefinitionNameHandler provided, see definition_name.go
	DefinitionNameHandler DefinitionNameHandlerFunc

	// --- oastools builder configuration (for BuildOAS2/BuildOAS3) ---

	// [optional] OASVersion specifies target OAS version for BuildOAS3.
	// Default is OASVersion320 (3.2.0). Only applies to BuildOAS3; BuildOAS2 always uses 2.0.
	OASVersion parser.OASVersion
	// [optional] SchemaNaming controls how Go type names are converted to schema names.
	// Default is SchemaNamingDefault (package.TypeName format).
	SchemaNaming builder.SchemaNamingStrategy
	// [optional] GenericNaming controls how generic type parameters are formatted in schema names.
	// Default is GenericNamingUnderscore (Response_User_ format).
	GenericNaming builder.GenericNamingStrategy
	// [optional] SemanticDeduplication enables schema deduplication.
	// When true, structurally identical schemas are consolidated. Default is false.
	SemanticDeduplication bool
	// [optional] SchemaNameFunc provides custom schema naming logic.
	// Takes precedence over SchemaNaming if set.
	SchemaNameFunc builder.SchemaNameFunc
	// [optional] Servers specifies server objects for OAS 3.x output (BuildOAS3 only).
	Servers []*parser.Server
	// [optional] Info provides additional metadata for the Info object.
	// Title and Version can come from existing APIVersion field or this Info.
	Info *parser.Info
	// [optional] PostBuildOAS2Handler is called after building an OAS 2.0 document.
	PostBuildOAS2Handler func(doc *parser.OAS2Document)
	// [optional] PostBuildOAS3Handler is called after building an OAS 3.x document.
	PostBuildOAS3Handler func(doc *parser.OAS3Document)
}

// Validate checks the Config for common misconfigurations and returns an error if any are found.
// This can be called before BuildSwagger, BuildOAS2, or BuildOAS3 to catch configuration issues early.
func (c Config) Validate() error {
	if len(c.WebServices) == 0 {
		return errors.New("restfulspec: WebServices is required but empty")
	}

	// Check for nil WebServices in the slice
	for i, ws := range c.WebServices {
		if ws == nil {
			return fmt.Errorf("restfulspec: WebServices contains nil entry at index %d", i)
		}
	}

	return nil
}
