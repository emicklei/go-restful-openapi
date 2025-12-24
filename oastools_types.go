package restfulspec

import (
	"github.com/erraggy/oastools/builder"
	"github.com/erraggy/oastools/parser"
)

// OAS2Document is an alias for parser.OAS2Document representing an OpenAPI 2.0 (Swagger) document.
type OAS2Document = parser.OAS2Document

// OAS3Document is an alias for parser.OAS3Document representing an OpenAPI 3.x document.
type OAS3Document = parser.OAS3Document

// Type aliases for oastools types used in Config.
type (
	// SchemaNamingStrategy defines built-in schema naming conventions.
	SchemaNamingStrategy = builder.SchemaNamingStrategy

	// GenericNamingStrategy defines how generic type parameters are formatted in schema names.
	GenericNamingStrategy = builder.GenericNamingStrategy

	// SchemaNameFunc provides custom schema naming logic.
	SchemaNameFunc = builder.SchemaNameFunc

	// OASVersion represents each canonical version of the OpenAPI Specification.
	OASVersion = parser.OASVersion

	// Server represents an OAS 3.x server object.
	Server = parser.Server

	// Info represents the API metadata.
	Info = parser.Info
)

// OAS version constants.
const (
	// OASVersionUnknown represents an unknown or unsupported OAS version.
	OASVersionUnknown = parser.Unknown
	// OASVersion20 represents OpenAPI Specification Version 2.0 (Swagger).
	OASVersion20 = parser.OASVersion20
	// OASVersion300 represents OpenAPI Specification Version 3.0.0.
	OASVersion300 = parser.OASVersion300
	// OASVersion301 represents OpenAPI Specification Version 3.0.1.
	OASVersion301 = parser.OASVersion301
	// OASVersion302 represents OpenAPI Specification Version 3.0.2.
	OASVersion302 = parser.OASVersion302
	// OASVersion303 represents OpenAPI Specification Version 3.0.3.
	OASVersion303 = parser.OASVersion303
	// OASVersion304 represents OpenAPI Specification Version 3.0.4.
	OASVersion304 = parser.OASVersion304
	// OASVersion310 represents OpenAPI Specification Version 3.1.0.
	OASVersion310 = parser.OASVersion310
	// OASVersion311 represents OpenAPI Specification Version 3.1.1.
	OASVersion311 = parser.OASVersion311
	// OASVersion312 represents OpenAPI Specification Version 3.1.2.
	OASVersion312 = parser.OASVersion312
	// OASVersion320 represents OpenAPI Specification Version 3.2.0.
	OASVersion320 = parser.OASVersion320
)

// Schema naming strategy constants.
const (
	// SchemaNamingDefault uses "package.TypeName" format (e.g., models.User).
	SchemaNamingDefault = builder.SchemaNamingDefault
	// SchemaNamingPascalCase uses "PackageTypeName" format (e.g., ModelsUser).
	SchemaNamingPascalCase = builder.SchemaNamingPascalCase
	// SchemaNamingCamelCase uses "packageTypeName" format (e.g., modelsUser).
	SchemaNamingCamelCase = builder.SchemaNamingCamelCase
	// SchemaNamingSnakeCase uses "package_type_name" format (e.g., models_user).
	SchemaNamingSnakeCase = builder.SchemaNamingSnakeCase
	// SchemaNamingKebabCase uses "package-type-name" format (e.g., models-user).
	SchemaNamingKebabCase = builder.SchemaNamingKebabCase
	// SchemaNamingTypeOnly uses just "TypeName" without package (may cause conflicts).
	SchemaNamingTypeOnly = builder.SchemaNamingTypeOnly
	// SchemaNamingFullPath uses full package path (e.g., github.com_org_models_User).
	SchemaNamingFullPath = builder.SchemaNamingFullPath
)

// Generic naming strategy constants.
const (
	// GenericNamingUnderscore replaces brackets with underscores (e.g., Response_User_).
	GenericNamingUnderscore = builder.GenericNamingUnderscore
	// GenericNamingOf uses "Of" separator (e.g., ResponseOfUser).
	GenericNamingOf = builder.GenericNamingOf
	// GenericNamingFor uses "For" separator (e.g., ResponseForUser).
	GenericNamingFor = builder.GenericNamingFor
	// GenericNamingAngleBrackets uses angle brackets, URI-encoded in $ref (e.g., Response<User>).
	GenericNamingAngleBrackets = builder.GenericNamingAngleBrackets
	// GenericNamingFlattened removes brackets entirely (e.g., ResponseUser).
	GenericNamingFlattened = builder.GenericNamingFlattened
)
