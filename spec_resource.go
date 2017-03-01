package restfulspec

import (
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

// RegisterOpenAPIService adds the WebService that provides the API documentation of all services
// conform the OpenAPI documentation specifcation.
func RegisterOpenAPIService(config Config, wsContainer *restful.Container) {

	ws := new(restful.WebService)
	ws.Path(config.APIPath)
	ws.Produces(restful.MIME_JSON)
	if config.DisableCORS {
		ws.Filter(enableCORS)
	}

	// TEMP
	paths := spec.Paths{Paths: map[string]spec.PathItem{}}
	definitions := spec.Definitions{}
	for _, each := range config.WebServices {
		po, defs := buildPathsAndDefs(each)
		for path, item := range po.Paths {
			paths.Paths[path] = item
		}
		for name, d := range defs {
			definitions[name] = d
		}
	}

	sw := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger:     "2.0",
			Paths:       &paths,
			Definitions: definitions,
			Info:        &config.Info,
		},
	}
	if config.PostBuildSwaggerObjectHandler != nil {
		config.PostBuildSwaggerObjectHandler(sw)
	}

	res := specResource{swaggerSpec: sw}
	ws.Route(ws.GET("/").To(res.getSwagger))

	wsContainer.Add(ws)

	// Check paths for UI serving
	if config.SwaggerFilePath != "" && config.SwaggerPath != "" {
		swaggerPathSlash := config.SwaggerPath
		// path must end with slash /
		if "/" != config.SwaggerPath[len(config.SwaggerPath)-1:] {
			swaggerPathSlash += "/"
		}

		wsContainer.Handle(swaggerPathSlash, http.StripPrefix(swaggerPathSlash, http.FileServer(http.Dir(config.SwaggerFilePath))))
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
	swaggerSpec *spec.Swagger
}

func (s specResource) getSwagger(req *restful.Request, resp *restful.Response) {
	resp.WriteAsJson(s.swaggerSpec)
}
