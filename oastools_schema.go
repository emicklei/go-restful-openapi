package restfulspec

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/erraggy/oastools/parser"
)

// hasOASTag returns true if the field has an oas:"..." struct tag.
// This function is reserved for future dual tag support integration.
//
//nolint:unused
func hasOASTag(field reflect.StructField) bool {
	return field.Tag.Get("oas") != ""
}

// applyLegacyTags applies go-restful-openapi's legacy struct tags to a schema.
// This is only called when the field does NOT have an oas:"..." tag.
// Supported legacy tags: description, minimum, maximum, enum, format, type,
// unique, readOnly, optional, example, default, x-nullable, x-go-name
// This function is reserved for future dual tag support integration.
//
//nolint:unused
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
		}
	}

	// maximum tag
	if tag := field.Tag.Get("maximum"); tag != "" {
		if f, err := strconv.ParseFloat(tag, 64); err == nil {
			result.Maximum = &f
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

	return result
}

// isLegacyFieldRequired determines if a field should be required based on legacy tags.
// Returns true if the field should be required, false otherwise.
// This function is reserved for future dual tag support integration.
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
// This function is reserved for future dual tag support integration.
//
//nolint:unused
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
// This function is reserved for future dual tag support integration.
//
//nolint:unused
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
	case "number":
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	case "boolean":
		return value == "true"
	}

	return value
}
