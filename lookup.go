package restfulspec

import restful "github.com/emicklei/go-restful/v3"

func asParamType(kind int) string {
	switch kind {
	case restful.PathParameterKind:
		return "path"
	case restful.QueryParameterKind:
		return "query"
	case restful.BodyParameterKind:
		return "body"
	case restful.HeaderParameterKind:
		return "header"
	case restful.FormParameterKind:
		return "formData"
	default:
		return ""
	}
}
