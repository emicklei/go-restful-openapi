package restfulapi

import restful "github.com/emicklei/go-restful"

// specResource is a REST resource to serve the Open-API spec.
type specResource struct {
	config Config
}

// RegisterOpenAPIService add the WebService that provides the API documentation of all services
// conform the OpenAPI documentation specifcation.
func RegisterOpenAPIService(config Config, wsContainer *restful.Container) {}
