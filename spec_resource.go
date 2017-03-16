package restfulspec

import restful "github.com/emicklei/go-restful"
import "github.com/go-openapi/spec"

// NewOpenAPIService returns a new WebService that provides the API documentation of all services
// conform the OpenAPI documentation specifcation.
func NewOpenAPIService(config Config) *restful.WebService {

	ws := new(restful.WebService)
	ws.Path(config.APIPath)
	ws.Produces(restful.MIME_JSON)
	if config.DisableCORS {
		ws.Filter(enableCORS)
	}

	resource := specResource{swagger: NewOpenAPISpecFromServices(config)}
	ws.Route(ws.GET("/").To(resource.getSwagger))

	return ws
}

func NewOpenAPISpecFromServices(config Config) *spec.Swagger {

	// collect paths and model definitions to build Swagger object.
	paths := &spec.Paths{Paths: map[string]spec.PathItem{}}
	definitions := spec.Definitions{}

	for _, each := range config.WebServices {
		for path, item := range BuildPaths(each).Paths {
			paths.Paths[path] = item
		}
		for name, def := range BuildDefinitions(each, config) {
			definitions[name] = def
		}
	}
	return &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger:     "2.0",
			Paths:       paths,
			Definitions: definitions,
		},
	}
}

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
	swagger *spec.Swagger
}

func (s specResource) getSwagger(req *restful.Request, resp *restful.Response) {
	resp.WriteAsJson(s.swagger)
}
