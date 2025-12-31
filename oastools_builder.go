package restfulspec

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/erraggy/oastools/builder"
	"github.com/erraggy/oastools/parser"
)

// BuildOAS2 builds an OAS 2.0 (Swagger) document using the oastools builder.
// This is an alternative to BuildSwagger that uses the oastools library for document generation.
func BuildOAS2(config Config) (*OAS2Document, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	b := newOASBuilder(config, parser.OASVersion20)

	// Configure info
	configureInfo(b, config)

	// Process all WebServices
	for _, ws := range config.WebServices {
		addWebServiceToBuilder(b, ws, config)
	}

	// Apply semantic deduplication if enabled
	if config.SemanticDeduplication {
		b.DeduplicateSchemas()
	}

	// Build the document
	doc, err := b.BuildOAS2()
	if err != nil {
		return nil, err
	}

	// Set host and schemes from config
	if config.Host != "" {
		doc.Host = config.Host
	}
	if len(config.Schemes) > 0 {
		doc.Schemes = config.Schemes
	}

	// Apply legacy optional tag post-processing
	if config.LegacyStructTags && doc.Definitions != nil {
		applyLegacyOptionalToSchemas(doc.Definitions)
	}

	// Call post-build handler if set
	if config.PostBuildOAS2Handler != nil {
		config.PostBuildOAS2Handler(doc)
	}

	return doc, nil
}

// BuildOAS3 builds an OAS 3.x document using the oastools builder.
// The OAS version can be configured via Config.OASVersion (defaults to 3.2.0).
// Returns an error if OASVersion20 is explicitly set - use BuildOAS2 instead.
func BuildOAS3(config Config) (*OAS3Document, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	version := config.OASVersion
	if version == parser.OASVersion20 {
		return nil, fmt.Errorf("BuildOAS3 called with OASVersion20; use BuildOAS2 for OAS 2.0 output")
	}
	if version == 0 {
		version = parser.OASVersion320 // Default to 3.2.0
	}

	b := newOASBuilder(config, version)

	// Configure info
	configureInfo(b, config)

	// Add servers from config
	for _, server := range config.Servers {
		b.AddServer(server.URL)
	}

	// If no servers but host is specified, create a server from it
	if len(config.Servers) == 0 && config.Host != "" {
		scheme := "https"
		if len(config.Schemes) > 0 {
			scheme = config.Schemes[0]
		}
		b.AddServer(scheme + "://" + config.Host)
	}

	// Process all WebServices
	for _, ws := range config.WebServices {
		addWebServiceToBuilder(b, ws, config)
	}

	// Apply semantic deduplication if enabled
	if config.SemanticDeduplication {
		b.DeduplicateSchemas()
	}

	// Build the document
	doc, err := b.BuildOAS3()
	if err != nil {
		return nil, err
	}

	// Apply legacy optional tag post-processing
	if config.LegacyStructTags && doc.Components != nil && doc.Components.Schemas != nil {
		applyLegacyOptionalToSchemas(doc.Components.Schemas)
	}

	// Call post-build handler if set
	if config.PostBuildOAS3Handler != nil {
		config.PostBuildOAS3Handler(doc)
	}

	return doc, nil
}

// newOASBuilder creates a new oastools builder with configuration from Config.
func newOASBuilder(config Config, version parser.OASVersion) *builder.Builder {
	var opts []builder.BuilderOption

	// Schema naming strategy
	if config.SchemaNameFunc != nil {
		opts = append(opts, builder.WithSchemaNameFunc(config.SchemaNameFunc))
	} else if config.SchemaNaming != 0 {
		opts = append(opts, builder.WithSchemaNaming(config.SchemaNaming))
	}

	// Generic naming strategy
	if config.GenericNaming != 0 {
		opts = append(opts, builder.WithGenericNaming(config.GenericNaming))
	}

	// Semantic deduplication
	if config.SemanticDeduplication {
		opts = append(opts, builder.WithSemanticDeduplication(true))
	}

	// Legacy struct tag support
	if config.LegacyStructTags {
		opts = append(opts, builder.WithSchemaFieldProcessor(legacyTagFieldProcessor))
	}

	return builder.New(version, opts...)
}

// configureInfo sets up the Info object on the builder.
func configureInfo(b *builder.Builder, config Config) {
	if config.Info != nil {
		b.SetInfo(config.Info)
	} else {
		// Use APIVersion as the version if no Info is provided
		if config.APIVersion != "" {
			b.SetVersion(config.APIVersion)
		}
	}
}

// addWebServiceToBuilder processes a WebService and adds its routes to the builder.
func addWebServiceToBuilder(b *builder.Builder, ws *restful.WebService, config Config) {
	for _, route := range ws.Routes() {
		addRouteToBuilder(b, ws, route, config)
	}
}

// addRouteToBuilder adds a single route to the builder.
func addRouteToBuilder(b *builder.Builder, ws *restful.WebService, route restful.Route, config Config) {
	path, patterns := sanitizePath(route.Path)

	var opts []builder.OperationOption

	// Operation metadata
	if route.Operation != "" {
		opts = append(opts, builder.WithOperationID(route.Operation))
	}
	if route.Doc != "" {
		opts = append(opts, builder.WithSummary(stripTags(route.Doc)))
	}
	if route.Notes != "" {
		opts = append(opts, builder.WithDescription(route.Notes))
	}
	if route.Deprecated {
		opts = append(opts, builder.WithDeprecated(true))
	}

	// Tags from metadata
	if route.Metadata != nil {
		if tags, ok := route.Metadata[KeyOpenAPITags]; ok {
			if tagList, ok := tags.([]string); ok {
				opts = append(opts, builder.WithTags(tagList...))
			}
		}
	}

	// OAS 2.0 specific: consumes and produces
	if len(route.Consumes) > 0 {
		opts = append(opts, builder.WithConsumes(route.Consumes...))
	}
	if len(route.Produces) > 0 {
		opts = append(opts, builder.WithProduces(route.Produces...))
	}

	// Vendor extensions
	if len(route.Extensions) > 0 {
		for key, value := range route.Extensions {
			if strings.HasPrefix(key, ExtensionPrefix) {
				opts = append(opts, builder.WithOperationExtension(key, value))
			}
		}
	}

	// Path parameters from WebService level
	for _, param := range ws.PathParameters() {
		if opt := mapParameter(param.Data(), patterns[param.Data().Name], config); opt != nil {
			opts = append(opts, opt)
		}
	}

	// Route-specific parameters
	for _, param := range route.ParameterDocs {
		if opt := mapParameter(param.Data(), patterns[param.Data().Name], config); opt != nil {
			opts = append(opts, opt)
		}
	}

	// Request body from ReadSample
	if route.ReadSample != nil {
		contentType := "application/json"
		if len(route.Consumes) > 0 {
			contentType = route.Consumes[0]
		}
		opts = append(opts, builder.WithRequestBody(contentType, route.ReadSample))
	}

	// Responses
	for code, respErr := range route.ResponseErrors {
		if opt := mapResponse(code, respErr, config); opt != nil {
			opts = append(opts, opt)
		}
	}

	// Default response
	if route.DefaultResponse != nil {
		if opt := mapDefaultResponse(*route.DefaultResponse, config); opt != nil {
			opts = append(opts, opt)
		}
	}

	// If no responses defined, add a default 200 OK
	if len(route.ResponseErrors) == 0 && route.DefaultResponse == nil {
		// Use WriteSample if available
		if route.WriteSample != nil {
			opts = append(opts, builder.WithResponse(http.StatusOK, route.WriteSample,
				builder.WithResponseDescription(http.StatusText(http.StatusOK))))
		} else {
			// Add empty 200 response
			opts = append(opts, builder.WithResponse(http.StatusOK, nil,
				builder.WithResponseDescription(http.StatusText(http.StatusOK))))
		}
	}

	b.AddOperation(route.Method, path, opts...)
}

// mapParameter converts a go-restful parameter to an oastools operation option.
func mapParameter(param restful.ParameterData, pattern string, _ Config) builder.OperationOption {
	var paramOpts []builder.ParamOption

	if param.Description != "" {
		paramOpts = append(paramOpts, builder.WithParamDescription(param.Description))
	}
	if param.Required {
		paramOpts = append(paramOpts, builder.WithParamRequired(true))
	}
	if param.DefaultValue != "" {
		paramOpts = append(paramOpts, builder.WithParamDefault(stringAutoType(param.DataType, param.DefaultValue)))
	}

	// Enum values - prefer PossibleValues over deprecated AllowableValues
	if len(param.PossibleValues) > 0 {
		enumVals := make([]any, len(param.PossibleValues))
		for i, v := range param.PossibleValues {
			enumVals[i] = stringAutoType(param.DataType, v)
		}
		paramOpts = append(paramOpts, builder.WithParamEnum(enumVals...))
	} else if len(param.AllowableValues) > 0 {
		enumVals := make([]any, 0, len(param.AllowableValues))
		for k := range param.AllowableValues {
			enumVals = append(enumVals, stringAutoType(param.DataType, k))
		}
		paramOpts = append(paramOpts, builder.WithParamEnum(enumVals...))
	}

	// Constraints
	if param.Minimum != nil {
		paramOpts = append(paramOpts, builder.WithParamMinimum(*param.Minimum))
	}
	if param.Maximum != nil {
		paramOpts = append(paramOpts, builder.WithParamMaximum(*param.Maximum))
	}
	if param.MinLength != nil {
		paramOpts = append(paramOpts, builder.WithParamMinLength(int(*param.MinLength)))
	}
	if param.MaxLength != nil {
		paramOpts = append(paramOpts, builder.WithParamMaxLength(int(*param.MaxLength)))
	}

	// Pattern - use extracted regex pattern for path params, or param.Pattern otherwise
	if param.Kind == restful.PathParameterKind && pattern != "" {
		paramOpts = append(paramOpts, builder.WithParamPattern(pattern))
	} else if param.Pattern != "" {
		paramOpts = append(paramOpts, builder.WithParamPattern(param.Pattern))
	}

	// DataFormat - explicitly set the parameter format if provided
	if param.DataFormat != "" {
		paramOpts = append(paramOpts, builder.WithParamFormat(param.DataFormat))
	}

	// AllowEmptyValue - only valid for query and form parameters
	if param.AllowEmptyValue && (param.Kind == restful.QueryParameterKind || param.Kind == restful.FormParameterKind) {
		paramOpts = append(paramOpts, builder.WithParamAllowEmptyValue(true))
	}

	// AllowMultiple - handle array parameters
	if param.AllowMultiple {
		// For array parameters, add collection format and item constraints
		if param.CollectionFormat != "" {
			paramOpts = append(paramOpts, builder.WithParamCollectionFormat(param.CollectionFormat))
		}
		if param.MinItems != nil {
			paramOpts = append(paramOpts, builder.WithParamMinItems(int(*param.MinItems)))
		}
		if param.MaxItems != nil {
			paramOpts = append(paramOpts, builder.WithParamMaxItems(int(*param.MaxItems)))
		}
		if param.UniqueItems {
			paramOpts = append(paramOpts, builder.WithParamUniqueItems(true))
		}
	}

	// Get the Go type for the parameter
	paramType := getTypeForDataType(param.DataType)

	// For AllowMultiple, wrap the type in a slice
	if param.AllowMultiple {
		paramType = []string{} // Use slice type for array parameters
	}

	// Map parameter kind to oastools parameter type
	switch param.Kind {
	case restful.PathParameterKind:
		return builder.WithPathParam(param.Name, paramType, paramOpts...)
	case restful.QueryParameterKind:
		return builder.WithQueryParam(param.Name, paramType, paramOpts...)
	case restful.HeaderParameterKind:
		return builder.WithHeaderParam(param.Name, paramType, paramOpts...)
	case restful.FormParameterKind:
		return builder.WithFormParam(param.Name, paramType, paramOpts...)
	case restful.BodyParameterKind:
		// Body parameters are handled separately via ReadSample
		return nil
	default:
		log.Printf("restfulspec: unknown parameter kind %d for parameter %q, skipping", param.Kind, param.Name)
		return nil
	}
}

// buildResponseOptions creates common response options from a ResponseError.
// This shared helper is used by both mapResponse and mapDefaultResponse.
func buildResponseOptions(resp restful.ResponseError) []builder.ResponseOption {
	respOpts := make([]builder.ResponseOption, 0, 1+len(resp.Headers)+len(resp.Extensions))

	if resp.Message != "" {
		respOpts = append(respOpts, builder.WithResponseDescription(resp.Message))
	}

	// Handle response headers
	for name, header := range resp.Headers {
		respOpts = append(respOpts, builder.WithResponseHeader(name, &parser.Header{
			Description: header.Description,
			Schema:      &parser.Schema{Type: header.Type, Format: header.Format},
		}))
	}

	// Handle extensions
	for key, value := range resp.Extensions {
		if strings.HasPrefix(key, ExtensionPrefix) {
			respOpts = append(respOpts, builder.WithResponseExtension(key, value))
		}
	}

	return respOpts
}

// mapResponse converts a go-restful ResponseError to an oastools response option.
func mapResponse(code int, resp restful.ResponseError, _ Config) builder.OperationOption {
	return builder.WithResponse(code, resp.Model, buildResponseOptions(resp)...)
}

// mapDefaultResponse converts a go-restful ResponseError to a default response option.
func mapDefaultResponse(resp restful.ResponseError, _ Config) builder.OperationOption {
	return builder.WithDefaultResponse(resp.Model, buildResponseOptions(resp)...)
}

// getTypeForDataType returns a Go type that matches the given data type string.
// This is used to provide type hints to the oastools builder for schema generation.
// Unknown data types will log a warning and default to string.
func getTypeForDataType(dataType string) any {
	switch dataType {
	case "string":
		return ""
	case "integer", "int":
		return int(0)
	case "int32":
		return int32(0)
	case "int64":
		return int64(0)
	case "number", "float64":
		return float64(0)
	case "float32":
		return float32(0)
	case "boolean", "bool":
		return false
	case "file":
		return nil // File uploads handled specially
	case "":
		// Empty data type defaults to string silently
		return ""
	default:
		// For unknown types, log a warning and default to string
		log.Printf("restfulspec: unknown data type %q, defaulting to string", dataType)
		return ""
	}
}
