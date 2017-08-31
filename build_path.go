package restfulspec

import (
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	restful "github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

// KeyOpenAPITags is a Metadata key for a restful Route
const KeyOpenAPITags = "openapi.tags"

func buildPaths(ws *restful.WebService, cfg Config) spec.Paths {
	p := spec.Paths{Paths: map[string]spec.PathItem{}}
	for _, each := range ws.Routes() {
		existingPathItem, ok := p.Paths[each.Path]
		if !ok {
			existingPathItem = spec.PathItem{}
		}
		p.Paths[each.Path] = buildPathItem(ws, each, existingPathItem, cfg)
	}
	return p
}

func buildPathItem(ws *restful.WebService, r restful.Route, existingPathItem spec.PathItem, cfg Config) spec.PathItem {
	op := buildOperation(ws, r, cfg)
	switch r.Method {
	case "GET":
		existingPathItem.Get = op
	case "POST":
		existingPathItem.Post = op
	case "PUT":
		existingPathItem.Put = op
	case "DELETE":
		existingPathItem.Delete = op
	case "PATCH":
		existingPathItem.Patch = op
	case "OPTIONS":
		existingPathItem.Options = op
	case "HEAD":
		existingPathItem.Head = op
	}
	return existingPathItem
}

func buildOperation(ws *restful.WebService, r restful.Route, cfg Config) *spec.Operation {
	o := spec.NewOperation(r.Operation)
	o.Description = r.Doc
	// take the first line (stripping HTML tags) to be the summary
	if lines := strings.Split(r.Doc, "\n"); len(lines) > 0 {
		o.Summary = stripTags(lines[0])
	}
	o.Consumes = r.Consumes
	o.Produces = r.Produces
	if r.Metadata != nil {
		if tags, ok := r.Metadata[KeyOpenAPITags]; ok {
			if tagList, ok := tags.([]string); ok {
				o.Tags = tagList
			}
		}
	}
	// collect any path parameters
	for _, param := range ws.PathParameters() {
		o.Parameters = append(o.Parameters, buildParameter(r, param, cfg))
	}
	// route specific params
	for _, each := range r.ParameterDocs {
		o.Parameters = append(o.Parameters, buildParameter(r, each, cfg))
	}
	o.Responses = new(spec.Responses)
	props := &o.Responses.ResponsesProps
	props.StatusCodeResponses = map[int]spec.Response{}
	for k, v := range r.ResponseErrors {
		r := buildResponse(v, cfg)
		props.StatusCodeResponses[k] = r
		if 200 == k { // any 2xx code?
			o.Responses.Default = &r
		}
	}
	if len(o.Responses.StatusCodeResponses) == 0 {
		o.Responses.StatusCodeResponses[200] = spec.Response{ResponseProps: spec.ResponseProps{Description: http.StatusText(http.StatusOK)}}
	}
	return o
}

// stringAutoType automatically picks the correct type from an ambiguously typed
// string. Ex. numbers become int, true/false become bool, etc.
func stringAutoType(ambiguous string) interface{} {
	if ambiguous == "" {
		return nil
	}
	if parsedInt, err := strconv.ParseInt(ambiguous, 10, 64); err == nil {
		return parsedInt
	}
	if parsedBool, err := strconv.ParseBool(ambiguous); err == nil {
		return parsedBool
	}
	return ambiguous
}

func buildParameter(r restful.Route, restfulParam *restful.Parameter, cfg Config) spec.Parameter {
	p := spec.Parameter{}
	param := restfulParam.Data()
	p.In = asParamType(param.Kind)
	p.Type = param.DataType
	p.Description = param.Description
	p.Name = param.Name
	p.Required = param.Required
	p.Default = stringAutoType(param.DefaultValue)
	p.Format = param.DataFormat

	if p.In == "body" && r.ReadSample != nil && p.Type == reflect.TypeOf(r.ReadSample).String() {
		p.Schema = new(spec.Schema)
		p.Schema.Ref = spec.MustCreateRef("#/definitions/" + p.Type)
		p.SimpleSchema = spec.SimpleSchema{}
	}
	return p
}

func buildResponse(e restful.ResponseError, cfg Config) (r spec.Response) {
	r.Description = e.Message
	if e.Model != nil {
		st := reflect.TypeOf(e.Model)
		if st.Kind() == reflect.Ptr {
			// For pointer type, use element type as the key; otherwise we'll
			// endup with '#/definitions/*Type' which violates openapi spec.
			st = st.Elem()
		}
		r.Schema = new(spec.Schema)
		if st.Kind() == reflect.Array || st.Kind() == reflect.Slice {
			modelName := definitionBuilder{}.keyFrom(st.Elem())
			r.Schema.Type = []string{"array"}
			r.Schema.Items = &spec.SchemaOrArray{
				Schema: &spec.Schema{},
			}
			isPrimitive := isPrimitiveType(modelName)
			if isPrimitive {
				mapped := jsonSchemaType(modelName)
				r.Schema.Items.Schema.Type = []string{mapped}
			} else {
				r.Schema.Items.Schema.Ref = spec.MustCreateRef("#/definitions/" + modelName)
			}
		} else {
			modelName := definitionBuilder{}.keyFrom(st)
			r.Schema.Ref = spec.MustCreateRef("#/definitions/" + modelName)
		}
	}
	return r
}

// stripTags takes a snippet of HTML and returns only the text content.
// For example, `<b>&lt;Hi!&gt;</b> <br>` -> `&lt;Hi!&gt; `.
func stripTags(html string) string {
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(html, "")
}

func isPrimitiveType(modelName string) bool {
	if len(modelName) == 0 {
		return false
	}
	return strings.Contains("uint uint8 uint16 uint32 uint64 int int8 int16 int32 int64 float32 float64 bool string byte rune time.Time", modelName)
}

func jsonSchemaType(modelName string) string {
	schemaMap := map[string]string{
		"uint":   "integer",
		"uint8":  "integer",
		"uint16": "integer",
		"uint32": "integer",
		"uint64": "integer",

		"int":   "integer",
		"int8":  "integer",
		"int16": "integer",
		"int32": "integer",
		"int64": "integer",

		"byte":      "integer",
		"float64":   "number",
		"float32":   "number",
		"bool":      "boolean",
		"time.Time": "string",
	}
	mapped, ok := schemaMap[modelName]
	if !ok {
		return modelName // use as is (custom or struct)
	}
	return mapped
}
