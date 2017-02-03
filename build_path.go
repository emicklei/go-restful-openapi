package restfulspec

import (
	restful "github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

func buildPath(r restful.Route) spec.PathItem {
	op := buildOperation(r)
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

func buildOperation(r restful.Route) *spec.Operation {
	o := spec.NewOperation(r.Operation)
	o.Description = r.Doc
	o.Consumes = r.Consumes
	o.Produces = r.Produces
	for _, each := range r.ParameterDocs {
		o.Parameters = append(o.Parameters, buildParameter(each))
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
