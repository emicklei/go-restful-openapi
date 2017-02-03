package restfulspec

import (
	restful "github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

func buildPath(r restful.Route) spec.PathItem {
	p := spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Get: buildOperation(r),
		},
	}
	return p
}

func buildOperation(r restful.Route) *spec.Operation {
	o := spec.NewOperation(r.Operation)
	return o
}
