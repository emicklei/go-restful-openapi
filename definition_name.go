package restfulspec

import "strings"

// DefinitionNameHandler generate name by this handler for definition without json tag.
// example: (for more, see file definition_name_test.go)
//   field	      			 definition_name
//   Name `json:"name"`  ->  name
// 	 Name                ->  Name
//
// there are some example provided for use
//   DefaultNameHandler         GoRestfulDefinition -> GoRestfulDefinition (not changed)
//   LowerSnakeCasedNameHandler  GoRestfulDefinition -> go_restful_definition
//   LowerCamelCasedNameHandler  GoRestfulDefinition -> goRestfulDefinition
//   GoLowerCamelCasedNameHandler HTTPRestfulDefinition -> httpRestfulDefinition
//
type DefinitionNameHandler interface {
	GetDefinitionName(string) string
}


// DefaultNameHandler GoRestfulDefinition -> GoRestfulDefinition (not changed)
type DefaultNameHandler struct{}

func (_ DefaultNameHandler) GetDefinitionName(name string) string {
	return name
}


// LowerSnakeCasedNameHandler GoRestfulDefinition -> go_restful_definition
type LowerSnakeCasedNameHandler struct{}

func (_ LowerSnakeCasedNameHandler) GetDefinitionName(name string) string {
	definitionName := make([]byte, 0, len(name)+1)
	for i := 0; i < len(name); i++ {
		c := name[i]
		if isUpper := 'A' <= c && c <= 'Z'; isUpper {
			if i > 0 {
				definitionName = append(definitionName, '_')
			}
			c += 'a' - 'A'
		}
		definitionName = append(definitionName, c)
	}

	return string(definitionName)
}


// LowerCamelCasedNameHandler GoRestfulDefinition -> goRestfulDefinition
type LowerCamelCasedNameHandler struct{}

func (_ LowerCamelCasedNameHandler) GetDefinitionName(name string) string {
	definitionName := make([]byte, 0, len(name)+1)
	for i := 0; i < len(name); i++ {
		c := name[i]
		if isUpper(c) && i == 0 {
			c += 'a' - 'A'
		}
		definitionName = append(definitionName, c)
	}

	return string(definitionName)
}


// GoLowerCamelCasedNameHandler HTTPRestfulDefinition -> httpRestfulDefinition
type GoLowerCamelCasedNameHandler struct{}

func (_ GoLowerCamelCasedNameHandler) GetDefinitionName(name string) string {
	var i = 0
	// for continuous Upper letters, check whether is it a common Initialisms
	for ; i < len(name) && isUpper(name[i]); i++ {}
	if len(name) != i && i != 1 {
		i-- // for continuous Upper letters, the last Upper is should not be check, eg: S for HTTPStatus
	}
	for ;i > 1;i-- {
		if _, ok := commonInitialisms[name[:i]];ok {
			break
		}
	}

	return strings.ToLower(name[:i])+ name[i:]
}


// commonInitialisms is a set of common initialisms. (from https://github.com/golang/lint/blob/master/lint.go)
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

func isUpper(r uint8) bool {
	return 'A' <= r && r <= 'Z'
}
