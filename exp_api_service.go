package restfulspec

import (
	restful "github.com/emicklei/go-restful/v3"
	"github.com/getkin/kin-openapi/openapi3"
)

// NewOpenAPIService returns a new WebService that provides the API documentation of all services
// conform the OpenAPI documentation specifcation.
func NewOpenAPIService(config Config) *restful.WebService {

	ws := new(restful.WebService)
	ws.Path(config.APIPath)
	ws.Produces(restful.MIME_JSON)
	if !config.DisableCORS {
		ws.Filter(enableCORS)
	}

	doc := BuildDoc(config)
	resource := specResource{doc: doc}
	ws.Route(ws.GET("/").To(resource.getSchema))
	return ws
}

// BuildDoc returns a OpenAPI object for all services' API endpoints.
func BuildDoc(config Config) *openapi3.T { return nil }

func enableCORS(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if origin := req.HeaderParameter(restful.HEADER_Origin); origin != "" {
		// prevent duplicate header
		if len(resp.Header().Get(restful.HEADER_AccessControlAllowOrigin)) == 0 {
			resp.AddHeader(restful.HEADER_AccessControlAllowOrigin, origin)
		}
	}
	chain.ProcessFilter(req, resp)
}

// specResource is a REST resource to serve the Open-API spec.
type specResource struct {
	doc *openapi3.T
}

func (s specResource) getSchema(req *restful.Request, resp *restful.Response) {
	resp.WriteAsJson(s.doc)
}

func asParamType(kind int) string {
	switch {
	case kind == restful.PathParameterKind:
		return "path"
	case kind == restful.QueryParameterKind:
		return "query"
	case kind == restful.BodyParameterKind:
		return "body"
	case kind == restful.HeaderParameterKind:
		return "header"
	case kind == restful.FormParameterKind:
		return "formData"
	}
	return ""
}
