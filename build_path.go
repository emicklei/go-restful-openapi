package restfulspec

import (
	"reflect"
	"strings"

	restful "github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

// KeyOpenAPITags is a Metadata key for a restful Route
const KeyOpenAPITags = "openapi.tags"

func buildPathsAndDefs(ws *restful.WebService) (spec.Paths, spec.Definitions) {
	p := spec.Paths{Paths: map[string]spec.PathItem{}}
	d := spec.Definitions{}
	for _, each := range ws.Routes() {
		op := buildOperation(ws, each)
		existingPathItem, ok := p.Paths[each.Path]
		if !ok {
			existingPathItem = spec.PathItem{}
		}

		switch each.Method {
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

		p.Paths[each.Path] = existingPathItem
		definitions := buildDefinitions(ws, each)
		for name, defn := range definitions {
			d[name] = defn
		}
	}
	return p, d
}

func buildOperation(ws *restful.WebService, r restful.Route) *spec.Operation {
	o := spec.NewOperation(r.Operation)
	o.Description = r.Doc
	// take the first line to be the summary
	if lines := strings.Split(r.Doc, "\n"); len(lines) > 0 {
		o.Summary = lines[0]
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
		o.Parameters = append(o.Parameters, buildParameter(param))
	}
	// route specific params
	for _, each := range r.ParameterDocs {
		o.Parameters = append(o.Parameters, buildParameter(each))
	}
	o.Responses = new(spec.Responses)
	props := o.Responses.ResponsesProps
	props.StatusCodeResponses = map[int]spec.Response{}
	for k, v := range r.ResponseErrors {
		r := buildResponse(v)
		props.StatusCodeResponses[k] = r
		if 200 == k { // any 2xx code?
			o.Responses.Default = &r
		}
	}
	return o
}

func buildParameter(r *restful.Parameter) spec.Parameter {
	p := spec.Parameter{}
	param := r.Data()
	p.In = asParamType(param.Kind)
	p.Description = param.Description
	p.Name = param.Name
	p.Required = param.Required
	p.Default = param.DefaultValue
	p.Format = param.DataFormat
	return p
}

func buildResponse(e restful.ResponseError) (r spec.Response) {
	r.Description = e.Message
	if e.Model != nil {
		st := reflect.TypeOf(e.Model)
		modelName := modelBuilder{}.keyFrom(st)
		r.Schema = &spec.Schema{}
		r.Schema.Ref = spec.MustCreateRef("#/definitions/" + modelName)
	}
	return r
}

func buildDefinitions(ws *restful.WebService, r restful.Route) spec.Definitions {
	definitions := spec.Definitions{}
	for _, v := range r.ResponseErrors {
		if v.Model == nil {
			continue
		}
		st := reflect.TypeOf(v.Model)
		mb := modelBuilder{Definitions: definitions, Config: nil}
		mb.addModel(st, "")
	}
	return definitions
}
