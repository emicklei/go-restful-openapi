# go-restful-openapi

[![Build Status](https://travis-ci.org/emicklei/go-restful-openapi.png)](https://travis-ci.org/emicklei/go-restful-openapi)
[![GoDoc](https://godoc.org/github.com/emicklei/go-restful-openapi?status.svg)](https://godoc.org/github.com/emicklei/go-restful-openapi)

Generate [OpenAPI](https://www.openapis.org) specifications from [go-restful](https://github.com/emicklei/go-restful) WebServices.

**Supported OpenAPI versions:**
- OpenAPI 3.2.0, 3.1.x, 3.0.x
- OpenAPI 2.0 (Swagger)

## Installation

```bash
go get github.com/emicklei/go-restful-openapi/v2
```

## Quick Start

### OpenAPI 3.x (Recommended)

```go
import (
    restfulspec "github.com/emicklei/go-restful-openapi/v2"
    restful "github.com/emicklei/go-restful/v3"
)

func main() {
    ws := new(restful.WebService)
    ws.Path("/users")
    ws.Route(ws.GET("/").To(listUsers).
        Doc("List all users").
        Writes([]User{}))

    config := restfulspec.Config{
        WebServices: []*restful.WebService{ws},
        OASVersion:  restfulspec.OASVersion320,
        Servers: []*restfulspec.Server{
            {URL: "https://api.example.com"},
        },
        Info: &restfulspec.Info{
            Title:       "User API",
            Description: "API for managing users",
            Version:     "1.0.0",
        },
    }

    doc, err := restfulspec.BuildOAS3(config)
    if err != nil {
        log.Fatal(err)
    }
    // Serialize to JSON/YAML, serve via endpoint, write to file, etc.
}
```

### OpenAPI 2.0

```go
config := restfulspec.Config{
    WebServices: []*restful.WebService{ws},
    Host:        "api.example.com",
    Schemes:     []string{"https"},
}

doc, err := restfulspec.BuildOAS2(config)
if err != nil {
    log.Fatal(err)
}
```

## API

| Function | Description |
|----------|-------------|
| `BuildOAS3(config) (*OAS3Document, error)` | Build an OpenAPI 3.x document |
| `BuildOAS2(config) (*OAS2Document, error)` | Build an OpenAPI 2.0 document |

### OAS Version Constants

```go
restfulspec.OASVersion320  // OpenAPI 3.2.0 (default)
restfulspec.OASVersion311  // OpenAPI 3.1.1
restfulspec.OASVersion310  // OpenAPI 3.1.0
restfulspec.OASVersion304  // OpenAPI 3.0.4
restfulspec.OASVersion303  // OpenAPI 3.0.3
restfulspec.OASVersion302  // OpenAPI 3.0.2
restfulspec.OASVersion301  // OpenAPI 3.0.1
restfulspec.OASVersion300  // OpenAPI 3.0.0
restfulspec.OASVersion20   // OpenAPI 2.0 (Swagger)
```

## Configuration

```go
type Config struct {
    // Required: services to document
    WebServices []*restful.WebService

    // API metadata
    Info       *Info      // Title, description, version, contact, license
    APIVersion string     // Alternative to Info.Version

    // Server configuration
    Servers []*Server     // OAS 3.x server objects
    Host    string        // OAS 2.0 host (also used to generate Server for OAS 3.x)
    Schemes []string      // OAS 2.0 schemes (used with Host for OAS 3.x Server)

    // OAS version (for BuildOAS3, defaults to 3.2.0)
    OASVersion OASVersion

    // Schema naming
    SchemaNaming          SchemaNamingStrategy  // How to name schemas
    GenericNaming         GenericNamingStrategy // How to format generic types
    SemanticDeduplication bool                  // Consolidate identical schemas
    SchemaNameFunc        SchemaNameFunc        // Custom naming function

    // Post-build customization
    PostBuildOAS2Handler func(*OAS2Document)
    PostBuildOAS3Handler func(*OAS3Document)
}
```

## Schema Naming Strategies

Control how Go types are converted to schema names:

| Strategy | Example Output |
|----------|---------------|
| `SchemaNamingDefault` | `models.User` |
| `SchemaNamingTypeOnly` | `User` |
| `SchemaNamingPascalCase` | `ModelsUser` |
| `SchemaNamingCamelCase` | `modelsUser` |
| `SchemaNamingSnakeCase` | `models_user` |
| `SchemaNamingKebabCase` | `models-user` |
| `SchemaNamingFullPath` | `github.com_org_models_User` |

```go
config := restfulspec.Config{
    WebServices:  []*restful.WebService{ws},
    SchemaNaming: restfulspec.SchemaNamingTypeOnly,
}
```

## Generic Naming Strategies

Control how generic type parameters appear in schema names:

| Strategy | Example Output |
|----------|---------------|
| `GenericNamingUnderscore` | `Response_User_` |
| `GenericNamingOf` | `ResponseOfUser` |
| `GenericNamingFor` | `ResponseForUser` |
| `GenericNamingAngleBrackets` | `Response<User>` |
| `GenericNamingFlattened` | `ResponseUser` |

## Struct Tags

Customize schema generation using struct tags on your model types.

| Tag | Description | Example |
|-----|-------------|---------|
| `json` | Field name in JSON | `json:"user_name"` |
| `description` | Field description | `description:"User's email"` |
| `type` | Override type | `type:"string"` or `type:"[]string"` |
| `format` | OpenAPI format | `format:"email"` |
| `enum` | Allowed values (pipe-separated) | `enum:"active\|inactive"` |
| `default` | Default value | `default:"active"` |
| `example` | Example value | `example:"user@example.com"` |
| `minimum` | Minimum numeric value | `minimum:"0"` |
| `maximum` | Maximum numeric value | `maximum:"100"` |
| `optional` | Mark as not required | `optional:"true"` |
| `unique` | Array items must be unique | `unique:"true"` |
| `readOnly` | Field is read-only | `readOnly:"true"` |
| `x-nullable` | Field can be null | `x-nullable:"true"` |
| `modelDescription` | Model description | `modelDescription:"A user"` |

### Example

```go
type User struct {
    ID        int       `json:"id" readOnly:"true"`
    Email     string    `json:"email" format:"email" description:"User's email address"`
    Name      string    `json:"name" description:"Full name" example:"John Doe"`
    Age       int       `json:"age,omitempty" minimum:"0" maximum:"150"`
    Status    string    `json:"status" enum:"active|inactive|pending" default:"active"`
    CreatedAt time.Time `json:"created_at" readOnly:"true"`
}
```

## Post-Build Customization

Modify the generated document before it's returned:

```go
config := restfulspec.Config{
    WebServices: []*restful.WebService{ws},
    PostBuildOAS3Handler: func(doc *restfulspec.OAS3Document) {
        // Add security schemes, tags, or any other customization
        doc.Components.SecuritySchemes = map[string]*SecurityScheme{
            "bearerAuth": {
                Type:   "http",
                Scheme: "bearer",
            },
        }
    },
}
```

## Semantic Deduplication

When enabled, structurally identical schemas are consolidated into a single definition:

```go
config := restfulspec.Config{
    WebServices:           []*restful.WebService{ws},
    SemanticDeduplication: true,
}
```

---

## Legacy API

The following API uses [go-openapi/spec](https://github.com/go-openapi/spec) types and only supports OpenAPI 2.0. Consider migrating to `BuildOAS2` or `BuildOAS3` for new projects.

```go
// BuildSwagger returns a *spec.Swagger (go-openapi/spec types)
swagger := restfulspec.BuildSwagger(config)

// NewOpenAPIService creates a WebService that serves the spec at config.APIPath
restful.Add(restfulspec.NewOpenAPIService(config))
```

### Legacy Configuration Options

These options only apply to `BuildSwagger` and `NewOpenAPIService`:

```go
type Config struct {
    // Path to serve the spec (e.g., "/apidocs.json")
    APIPath string

    // Disable CORS headers on the spec endpoint
    DisableCORS bool

    // Customize the Swagger object after generation
    PostBuildSwaggerObjectHandler func(*spec.Swagger)

    // Custom type format mapping
    SchemaFormatHandler MapSchemaFormatFunc

    // Custom type name mapping
    ModelTypeNameHandler MapModelTypeNameFunc

    // Custom definition name handling
    DefinitionNameHandler DefinitionNameHandlerFunc
}
```

---

## Go Modules

Version `v2` of this package requires:
- Go 1.18+
- `github.com/emicklei/go-restful/v3`

```go
import (
    restfulspec "github.com/emicklei/go-restful-openapi/v2"
    restful "github.com/emicklei/go-restful/v3"
)
```

For go-restful v2, use v1 of this package (OpenAPI 2.0 only):
```go
import (
    restfulspec "github.com/emicklei/go-restful-openapi"
    restful "github.com/emicklei/go-restful/v2"
)
```

## Dependencies

- [go-restful v3](https://github.com/emicklei/go-restful) - REST framework
- [oastools](https://github.com/erraggy/oastools) - OpenAPI document builder

## License

MIT License. Contributions welcome.

2017-2024, ernestmicklei.com
