package restfulspec

import (
	restful "github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

func buildPaths(ws *restful.WebService) spec.Paths {
	p := spec.Paths{Paths: map[string]spec.PathItem{}}
	for _, each := range ws.Routes() {
		p.Paths[each.Path] = buildPathItem(ws, each)
	}
	return p
}

func buildPathItem(ws *restful.WebService, r restful.Route) spec.PathItem {
	op := buildOperation(ws, r)
	props := spec.PathItemProps{}
	switch r.Method {
	case "GET":
		props.Get = op
	case "POST":
		props.Post = op
	case "PUT":
		props.Put = op
	case "DELETE":
		props.Delete = op
	case "PATCH":
		props.Patch = op
	case "OPTIONS":
		props.Options = op
	case "HEAD":
		props.Head = op
	}
	p := spec.PathItem{
		PathItemProps: props,
	}
	return p
}

func buildOperation(ws *restful.WebService, r restful.Route) *spec.Operation {
	o := spec.NewOperation(r.Operation)
	o.Description = r.Doc
	o.Consumes = r.Consumes
	o.Produces = r.Produces
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
	return r
}
