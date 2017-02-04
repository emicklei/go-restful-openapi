package restfulspec

import restful "github.com/emicklei/go-restful"
import "github.com/go-openapi/spec"

// RegisterOpenAPIService adds the WebService that provides the API documentation of all services
// conform the OpenAPI documentation specifcation.
func RegisterOpenAPIService(config Config, wsContainer *restful.Container) {

	ws := new(restful.WebService)
	ws.Path(config.APIPath)
	ws.Produces(restful.MIME_JSON)
	if config.DisableCORS {
		ws.Filter(enableCORS)
	}

	res := specResource{config: config, paths: spec.Paths{Paths: map[string]spec.PathItem{}}}
	ws.Route(ws.GET("/").To(res.getSwagger))

	// TEMP
	for _, each := range config.WebServices {
		po := buildPaths(each)
		for path, item := range po.Paths {
			res.paths.Paths[path] = item
		}
	}
	wsContainer.Add(ws)
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
	config Config
	paths  spec.Paths
}

func (s specResource) getSwagger(req *restful.Request, resp *restful.Response) {
	sw := spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Info:    &s.config.Info,
			Swagger: "2.0",
			Paths:   &(s.paths),
		},
	}
	resp.WriteAsJson(sw)
}
