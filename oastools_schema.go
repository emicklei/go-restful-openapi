package restfulspec

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/erraggy/oastools/parser"
)

// legacyOptionalExtension is the extension key used to mark fields as optional
// during schema generation. This is cleaned up in post-processing.
const legacyOptionalExtension = "x-restful-optional"

// legacyTagFieldProcessor is a schema field processor that applies go-restful-openapi's
// legacy struct tags to schemas. This is used with oastools' WithSchemaFieldProcessor option.
// Legacy tags are only applied to fields that do NOT have an oas:"..." tag.
func legacyTagFieldProcessor(schema *parser.Schema, field reflect.StructField) *parser.Schema {
	// Only apply legacy tags if the field doesn't have an oas:"..." tag
	if hasOASTag(field) {
		return schema
	}
	return applyLegacyTags(schema, field)
}

// hasOASTag returns true if the field has an oas:"..." struct tag.
func hasOASTag(field reflect.StructField) bool {
	return field.Tag.Get("oas") != ""
}

// applyLegacyTags applies go-restful-openapi's legacy struct tags to a schema.
// This is only called when the field does NOT have an oas:"..." tag.
// Supported legacy tags: description, minimum, maximum, enum, format, type,
// unique, readOnly, optional, example, default, x-nullable, x-go-name
func applyLegacyTags(schema *parser.Schema, field reflect.StructField) *parser.Schema {
	if schema == nil {
		return nil
	}

	// Make a copy to avoid mutating the original
	result := shallowCopySchema(schema)

	// description tag
	if tag := field.Tag.Get("description"); tag != "" {
		result.Description = tag
	}

	// minimum tag
	if tag := field.Tag.Get("minimum"); tag != "" {
		if f, err := strconv.ParseFloat(tag, 64); err == nil {
			result.Minimum = &f
		} else {
			log.Printf("restfulspec: field %q has invalid minimum tag %q: %v", field.Name, tag, err)
		}
	}

	// maximum tag
	if tag := field.Tag.Get("maximum"); tag != "" {
		if f, err := strconv.ParseFloat(tag, 64); err == nil {
			result.Maximum = &f
		} else {
			log.Printf("restfulspec: field %q has invalid maximum tag %q: %v", field.Name, tag, err)
		}
	}

	// enum tag (pipe-separated values)
	if tag := field.Tag.Get("enum"); tag != "" {
		enumValues := strings.Split(tag, "|")
		result.Enum = make([]any, len(enumValues))
		for i, v := range enumValues {
			result.Enum[i] = strings.TrimSpace(v)
		}
	}

	// format tag
	if tag := field.Tag.Get("format"); tag != "" {
		result.Format = tag
	}

	// type tag - can override the Go type
	if tag := field.Tag.Get("type"); tag != "" {
		// Check if the type is intended to be an array (e.g., "[]string")
		if len(tag) > 2 && tag[0:2] == "[]" {
			result.Type = "array"
			result.Items = &parser.Schema{Type: tag[2:]}
		} else {
			result.Type = tag
		}
	}

	// unique tag
	if tag := field.Tag.Get("unique"); tag == "true" {
		result.UniqueItems = true
	}

	// readOnly tag
	if tag := field.Tag.Get("readOnly"); tag == "true" {
		result.ReadOnly = true
	}

	// example tag
	if tag := field.Tag.Get("example"); tag != "" {
		result.Example = tag
	}

	// default tag
	if tag := field.Tag.Get("default"); tag != "" {
		result.Default = parseLegacyDefaultValue(tag, result.Type)
	}

	// x-nullable tag - stored as extension
	if tag := field.Tag.Get("x-nullable"); tag != "" {
		if result.Extra == nil {
			result.Extra = make(map[string]any)
		}
		result.Extra["x-nullable"] = tag == "true"
	}

	// x-go-name tag - stored as extension
	if tag := field.Tag.Get("x-go-name"); tag != "" {
		if result.Extra == nil {
			result.Extra = make(map[string]any)
		}
		result.Extra["x-go-name"] = tag
	}

	// optional tag - marks field as not required
	// This is processed via extension and cleaned up in post-processing
	// because the required array is on the parent schema, not the field schema.
	if tag := field.Tag.Get("optional"); tag == "true" {
		if result.Extra == nil {
			result.Extra = make(map[string]any)
		}
		// Get the JSON field name for removal from required array
		jsonName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				jsonName = parts[0]
			}
		}
		result.Extra[legacyOptionalExtension] = jsonName
	}

	return result
}

// isLegacyFieldRequired determines if a field should be required based on legacy tags.
// Returns true if the field should be required, false otherwise.
// Note: This function is reserved for future integration when oastools adds support
// for customizing required field behavior. Currently, required fields are determined
// by oastools based on pointer types and omitempty tags.
//
//nolint:unused
func isLegacyFieldRequired(field reflect.StructField) bool {
	// Check for explicit optional tag
	if optionalTag := field.Tag.Get("optional"); optionalTag == "true" {
		return false
	}

	// Check json tag for omitempty
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		for _, part := range parts[1:] {
			if part == "omitempty" {
				return false
			}
		}
	}

	// Pointer fields are optional by default
	if field.Type.Kind() == reflect.Pointer {
		return false
	}

	return true
}

// shallowCopySchema creates a shallow copy of a schema for modification.
func shallowCopySchema(s *parser.Schema) *parser.Schema {
	if s == nil {
		return nil
	}

	result := &parser.Schema{
		Ref:         s.Ref,
		Type:        s.Type,
		Format:      s.Format,
		Title:       s.Title,
		Description: s.Description,
		Default:     s.Default,
		Nullable:    s.Nullable,
		ReadOnly:    s.ReadOnly,
		WriteOnly:   s.WriteOnly,
		Deprecated:  s.Deprecated,
		Pattern:     s.Pattern,
		UniqueItems: s.UniqueItems,
		Example:     s.Example,
	}

	// Copy pointer fields
	if s.Minimum != nil {
		minCopy := *s.Minimum
		result.Minimum = &minCopy
	}
	if s.Maximum != nil {
		maxCopy := *s.Maximum
		result.Maximum = &maxCopy
	}
	if s.MinLength != nil {
		minLenCopy := *s.MinLength
		result.MinLength = &minLenCopy
	}
	if s.MaxLength != nil {
		maxLenCopy := *s.MaxLength
		result.MaxLength = &maxLenCopy
	}
	if s.MinItems != nil {
		minItemsCopy := *s.MinItems
		result.MinItems = &minItemsCopy
	}
	if s.MaxItems != nil {
		maxItemsCopy := *s.MaxItems
		result.MaxItems = &maxItemsCopy
	}

	// Copy slices
	if s.Enum != nil {
		result.Enum = make([]any, len(s.Enum))
		copy(result.Enum, s.Enum)
	}
	if s.Required != nil {
		result.Required = make([]string, len(s.Required))
		copy(result.Required, s.Required)
	}

	// Reference shared structures (immutable from our perspective)
	result.Properties = s.Properties
	result.Items = s.Items
	result.AdditionalProperties = s.AdditionalProperties
	result.AllOf = s.AllOf
	result.AnyOf = s.AnyOf
	result.OneOf = s.OneOf

	// Deep-copy Extra map to avoid shared mutation
	if s.Extra != nil {
		result.Extra = make(map[string]any, len(s.Extra))
		for k, v := range s.Extra {
			result.Extra[k] = v
		}
	}

	return result
}

// parseLegacyDefaultValue parses a default value string based on the schema type.
func parseLegacyDefaultValue(value string, schemaType any) any {
	typeStr, ok := schemaType.(string)
	if !ok {
		return value
	}

	switch typeStr {
	case "integer":
		if n, err := strconv.ParseInt(value, 10, 64); err == nil {
			return n
		}
		log.Printf("restfulspec: invalid default value %q for integer type: falling back to string", value)
	case "number":
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
		log.Printf("restfulspec: invalid default value %q for number type: falling back to string", value)
	case "boolean":
		return value == "true"
	}

	return value
}

// applyLegacyOptionalToSchemas processes all schemas in a map and removes fields
// marked with x-restful-optional from the required arrays, then cleans up the extension.
func applyLegacyOptionalToSchemas(schemas map[string]*parser.Schema) {
	for _, schema := range schemas {
		applyLegacyOptionalToSchema(schema)
	}
}

// applyLegacyOptionalToSchema processes a single schema and its nested properties.
func applyLegacyOptionalToSchema(schema *parser.Schema) {
	if schema == nil {
		return
	}

	// Collect field names that should be removed from required
	var optionalFields []string

	// Process properties
	for propName, propSchema := range schema.Properties {
		if propSchema == nil {
			continue
		}

		// Check for the optional extension
		if propSchema.Extra != nil {
			if jsonName, ok := propSchema.Extra[legacyOptionalExtension].(string); ok {
				optionalFields = append(optionalFields, jsonName)
				// Clean up the extension
				delete(propSchema.Extra, legacyOptionalExtension)
				if len(propSchema.Extra) == 0 {
					propSchema.Extra = nil
				}
			}
		}

		// Recursively process nested schemas
		applyLegacyOptionalToSchema(propSchema)

		// Also check if the property name matches what we collected
		_ = propName // propName is used implicitly via propSchema
	}

	// Remove optional fields from the required array
	if len(optionalFields) > 0 && len(schema.Required) > 0 {
		optionalSet := make(map[string]bool, len(optionalFields))
		for _, f := range optionalFields {
			optionalSet[f] = true
		}

		newRequired := make([]string, 0, len(schema.Required))
		for _, req := range schema.Required {
			if !optionalSet[req] {
				newRequired = append(newRequired, req)
			}
		}
		schema.Required = newRequired
	}

	// Process nested schemas in Items, AllOf, AnyOf, OneOf
	if schema.Items != nil {
		if itemSchema, ok := schema.Items.(*parser.Schema); ok {
			applyLegacyOptionalToSchema(itemSchema)
		}
	}
	for _, s := range schema.AllOf {
		applyLegacyOptionalToSchema(s)
	}
	for _, s := range schema.AnyOf {
		applyLegacyOptionalToSchema(s)
	}
	for _, s := range schema.OneOf {
		applyLegacyOptionalToSchema(s)
	}
	if schema.AdditionalProperties != nil {
		if addPropSchema, ok := schema.AdditionalProperties.(*parser.Schema); ok {
			applyLegacyOptionalToSchema(addPropSchema)
		}
	}
}
