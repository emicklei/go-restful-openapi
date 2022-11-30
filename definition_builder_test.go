package restfulspec

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/spec"
)

type StringAlias string
type IntAlias int

type Apple struct {
	Species     string
	Volume      int `json:"vol"`
	Things      *[]string
	Weight      json.Number
	StringAlias StringAlias
	IntAlias    IntAlias
}

func TestAppleDef(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(Apple{})

	if got, want := len(db.Definitions), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}

	schema := db.Definitions["restfulspec.Apple"]
	if got, want := len(schema.Required), 6; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Required[0], "Species"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Required[1], "vol"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.ID, ""; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Properties["Things"].Items.Schema.Type.Contains("string"), true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Properties["Things"].Items.Schema.Ref.String(), ""; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Properties["Weight"].Type.Contains("number"), true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Properties["StringAlias"].Type.Contains("string"), true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Properties["IntAlias"].Type.Contains("integer"), true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

type MyDictionaryResponse struct {
	Dictionary1 map[string]DictionaryValue   `json:"dictionary1"`
	Dictionary2 map[string]interface{}       `json:"dictionary2"`
	Dictionary3 map[string][]byte            `json:"dictionary3"`
	Dictionary4 map[string]string            `json:"dictionary4"`
	Dictionary5 map[string][]DictionaryValue `json:"dictionary5"`
	Dictionary6 map[string][]string          `json:"dictionary6"`
}
type DictionaryValue struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}

func TestDictionarySupport(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(MyDictionaryResponse{})

	// Make sure that only the types that we want were created.
	if got, want := len(db.Definitions), 2; got != want {
		t.Errorf("got %v want %v", got, want)
	}

	schema, schemaFound := db.Definitions["restfulspec.MyDictionaryResponse"]
	if !schemaFound {
		t.Errorf("could not find schema")
	} else {
		if got, want := len(schema.Required), 6; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if got, want := schema.Required[0], "dictionary1"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := schema.Required[1], "dictionary2"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := schema.Required[2], "dictionary3"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := schema.Required[3], "dictionary4"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := schema.Required[4], "dictionary5"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := schema.Required[5], "dictionary6"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		}
		if got, want := len(schema.Properties), 6; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if property, found := schema.Properties["dictionary1"]; !found {
				t.Errorf("could not find property")
			} else {
				if got, want := property.AdditionalProperties.Schema.SchemaProps.Ref.String(), "#/definitions/restfulspec.DictionaryValue"; got != want {
					t.Errorf("got %v want %v", got, want)
				}
			}
			if property, found := schema.Properties["dictionary2"]; !found {
				t.Errorf("could not find property")
			} else {
				if property.AdditionalProperties != nil {
					t.Errorf("unexpected additional properties")
				}
			}
			if property, found := schema.Properties["dictionary3"]; !found {
				t.Errorf("could not find property")
			} else {
				if got, want := property.AdditionalProperties.Schema.Type[0], "string"; got != want {
					t.Errorf("got %v want %v", got, want)
				}
			}
			if property, found := schema.Properties["dictionary4"]; !found {
				t.Errorf("could not find property")
			} else {
				if got, want := property.AdditionalProperties.Schema.Type[0], "string"; got != want {
					t.Errorf("got %v want %v", got, want)
				}
			}
			if property, found := schema.Properties["dictionary5"]; !found {
				t.Errorf("could not find property")
			} else {
				if got, want := len(property.AdditionalProperties.Schema.Type), 1; got != want {
					t.Errorf("got %v want %v", got, want)
				} else {
					if got, want := property.AdditionalProperties.Schema.Type[0], "array"; got != want {
						t.Errorf("got %v want %v", got, want)
					}
					if property.AdditionalProperties.Schema.Items == nil {
						t.Errorf("Items not set")
					} else {
						if got, want := property.AdditionalProperties.Schema.Items.Schema.Ref.String(), "#/definitions/restfulspec.DictionaryValue"; got != want {
							t.Errorf("got %v want %v", got, want)
						}
					}
				}
			}
			if property, found := schema.Properties["dictionary6"]; !found {
				t.Errorf("could not find property")
			} else {
				if got, want := len(property.AdditionalProperties.Schema.Type), 1; got != want {
					t.Errorf("got %v want %v", got, want)
				} else {
					if got, want := property.AdditionalProperties.Schema.Type[0], "array"; got != want {
						t.Errorf("got %v want %v", got, want)
					}
					if property.AdditionalProperties.Schema.Items == nil {
						t.Errorf("Items not set")
					} else {
						if got, want := property.AdditionalProperties.Schema.Items.Schema.Type[0], "string"; got != want {
							t.Errorf("got %v want %v", got, want)
						}
					}
				}
			}
		}
	}

	schema, schemaFound = db.Definitions["restfulspec.DictionaryValue"]
	if !schemaFound {
		t.Errorf("could not find schema")
	} else {
		if got, want := len(schema.Required), 2; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if got, want := schema.Required[0], "key1"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := schema.Required[1], "key2"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		}
	}
}

type MyRecursiveDictionaryResponse struct {
	Dictionary1 map[string]RecursiveDictionaryValue `json:"dictionary1"`
}
type RecursiveDictionaryValue struct {
	Key1 string                              `json:"key1"`
	Key2 map[string]RecursiveDictionaryValue `json:"key2"`
}

func TestRecursiveDictionarySupport(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(MyRecursiveDictionaryResponse{})

	// Make sure that only the types that we want were created.
	if got, want := len(db.Definitions), 2; got != want {
		t.Errorf("got %v want %v", got, want)
	}

	schema, schemaFound := db.Definitions["restfulspec.MyRecursiveDictionaryResponse"]
	if !schemaFound {
		t.Errorf("could not find schema")
	} else {
		if got, want := len(schema.Required), 1; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if got, want := schema.Required[0], "dictionary1"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		}
		if got, want := len(schema.Properties), 1; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if property, found := schema.Properties["dictionary1"]; !found {
				t.Errorf("could not find property")
			} else {
				if got, want := property.AdditionalProperties.Schema.SchemaProps.Ref.String(), "#/definitions/restfulspec.RecursiveDictionaryValue"; got != want {
					t.Errorf("got %v want %v", got, want)
				}
			}
		}
	}

	schema, schemaFound = db.Definitions["restfulspec.RecursiveDictionaryValue"]
	if !schemaFound {
		t.Errorf("could not find schema")
	} else {
		if got, want := len(schema.Required), 2; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if got, want := schema.Required[0], "key1"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := schema.Required[1], "key2"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		}
		if got, want := len(schema.Properties), 2; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if property, found := schema.Properties["key1"]; !found {
				t.Errorf("could not find property")
			} else {
				if property.AdditionalProperties != nil {
					t.Errorf("unexpected additional properties")
				}
			}
			if property, found := schema.Properties["key2"]; !found {
				t.Errorf("could not find property")
			} else {
				if got, want := property.AdditionalProperties.Schema.SchemaProps.Ref.String(), "#/definitions/restfulspec.RecursiveDictionaryValue"; got != want {
					t.Errorf("got %v want %v", got, want)
				}
			}
		}
	}
}

func TestReturningStringToStringDictionary(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(map[string]string{})

	if got, want := len(db.Definitions), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}

	schema, schemaFound := db.Definitions["map[string]string"]
	if !schemaFound {
		t.Errorf("could not find schema")
	} else {
		if schema.AdditionalProperties == nil {
			t.Errorf("AdditionalProperties not set")
		} else {
			if got, want := schema.AdditionalProperties.Schema.Type[0], "string"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		}
	}
}

func TestReturningStringToSliceObjectDictionary(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(map[string][]DictionaryValue{})

	if got, want := len(db.Definitions), 2; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	schema, schemaFound := db.Definitions["map[string]||restfulspec.DictionaryValue"]
	if !schemaFound {
		t.Errorf("could not find schema")
	} else {
		if schema.AdditionalProperties == nil {
			t.Errorf("AdditionalProperties not set")
		} else {
			if got, want := len(schema.AdditionalProperties.Schema.Type), 1; got != want {
				t.Errorf("got %v want %v", got, want)
			} else {
				if got, want := schema.AdditionalProperties.Schema.Type[0], "array"; got != want {
					t.Errorf("got %v want %v", got, want)
				}
				if schema.AdditionalProperties.Schema.Items == nil {
					t.Errorf("Items not set")
				} else {
					if got, want := schema.AdditionalProperties.Schema.Items.Schema.Ref.String(), "#/definitions/restfulspec.DictionaryValue"; got != want {
						t.Errorf("got %v want %v", got, want)
					}
				}
			}
		}
	}

	schema, schemaFound = db.Definitions["restfulspec.DictionaryValue"]
	if !schemaFound {
		t.Errorf("could not find schema")
	} else {
		if got, want := len(schema.Required), 2; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if got, want := schema.Required[0], "key1"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := schema.Required[1], "key2"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		}
	}
}

func TestAddSliceOfPrimitiveCreatesNoType(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom([]string{})

	if got, want := len(db.Definitions), 0; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

type StructForSlice struct {
	Value string
}

func TestAddSliceOfStructCreatesTypeForStruct(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom([]StructForSlice{})

	if got, want := len(db.Definitions), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	schema, schemaFound := db.Definitions["restfulspec.StructForSlice"]
	if !schemaFound {
		t.Errorf("could not find schema")
	} else {
		if got, want := len(schema.Required), 1; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if got, want := schema.Required[0], "Value"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		}
		if got, want := len(schema.Properties), 1; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if property, found := schema.Properties["Value"]; !found {
				t.Errorf("could not find property")
			} else {
				if property.AdditionalProperties != nil {
					t.Errorf("unexpected additional properties")
				}
			}
		}
	}
}

type (
	X struct {
		yy []Y
	}
	Y struct {
		X
	}
)

func TestPotentialStackOverflow(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(X{})

	if got, want := len(db.Definitions), 2; got != want {
		t.Errorf("got %v want %v", got, want)
	}

	schema := db.Definitions["restfulspec.X"]
	if got, want := len(schema.Required), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.Required[0], "yy"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := schema.ID, ""; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

type Foo struct {
	b *Bar
}

// nolint:unused
type Bar struct {
	Foo `json:"foo"`
	B   struct {
		Foo
		f Foo
		b []*Bar
	}
}

func TestRecursiveFieldStructure(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(Foo{})
	t.Log(db)
}

type Baz struct {
	Foo    `json:"foo"`
	person struct {
		name string `json:"name"`
		age  uint32 `json:"age"`
	} `json:"person"`
}

func TestAnonymousStructure(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(Baz{})
	schema, schemaFound := db.Definitions[".restfulspec.Baz.person"]
	if !schemaFound {
		t.Errorf("could not find schema")
	} else {
		if got, want := len(schema.Required), 2; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if got, want := schema.Required[0], "name"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := schema.Required[1], "age"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		}
		if got, want := len(schema.Properties), 2; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if property, found := schema.Properties["name"]; !found {
				t.Errorf("could not find property")
			} else {
				if property.AdditionalProperties != nil {
					t.Errorf("unexpected additional properties")
				}
			}
			if property, found := schema.Properties["age"]; !found {
				t.Errorf("could not find property")
			} else {
				if property.AdditionalProperties != nil {
					t.Errorf("unexpected additional properties")
				}
			}
		}
	}
	schema, schemaFound = db.Definitions["restfulspec.Baz"]
	if !schemaFound {
		t.Errorf("could not find schema")
	} else {
		if got, want := len(schema.Required), 2; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if got, want := schema.Required[0], "foo"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := schema.Required[1], "restfulspec.Baz.person"; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		}
		if got, want := len(schema.Properties), 2; got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			if property, found := schema.Properties["foo"]; !found {
				t.Errorf("could not find property")
			} else {
				if property.AdditionalProperties != nil {
					t.Errorf("unexpected additional properties")
				}
				if got, want := property.Ref.String(), "#/definitions/restfulspec.Foo"; got != want {
					t.Errorf("got %v want %v", got, want)
				}
			}
			if property, found := schema.Properties["restfulspec.Baz.person"]; !found {
				t.Errorf("could not find property")
			} else {
				if property.AdditionalProperties != nil {
					t.Errorf("unexpected additional properties")
				}
				if got, want := property.Ref.String(), "#/definitions/.restfulspec.Baz.person"; got != want {
					t.Errorf("got %v want %v", got, want)
				}
			}
		}
	}
}

type email struct {
	Attachments [][]byte `json:"attachments,omitempty" optional:"true"`
}

// Definition Builder fails with [][]byte #77
func TestDoubleByteArray(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(email{})
	sc, ok := db.Definitions["restfulspec.email.attachments"]
	if !ok {
		t.Fail()
	}
	t.Log(sc)
	if got, want := sc.Type[0], "array"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

type matrix struct {
	Cells [][]string
}

func TestDoubleStringArray(t *testing.T) {
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(matrix{})
	sc, ok := db.Definitions["restfulspec.matrix.Cells"]
	if !ok {
		t.Log(db.Definitions)
		t.Fail()
	}
	t.Log(sc)
	if got, want := sc.Type[0], "array"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	tp := sc.Items.Schema.Type
	t.Log(tp)
}

type childPostBuildSwaggerSchema struct {
	Value string
}

func (m childPostBuildSwaggerSchema) PostBuildSwaggerSchemaHandler(sm *spec.Schema) {
	sm.Description = "child's description"
}

type parentPostBuildSwaggerSchema struct {
	Node childPostBuildSwaggerSchema
}

func (m parentPostBuildSwaggerSchema) PostBuildSwaggerSchemaHandler(sm *spec.Schema) {
	sm.Description = "parent's description"
}

func TestPostBuildSwaggerSchema(t *testing.T) {
	var obj interface{} = parentPostBuildSwaggerSchema{}
	if _, ok := obj.(PostBuildSwaggerSchema); !ok {
		t.Fatalf("object does not implement PostBuildSwaggerSchema interface")
	}
	db := definitionBuilder{Definitions: spec.Definitions{}, Config: Config{}}
	db.addModelFrom(obj)
	sc, ok := db.Definitions["restfulspec.parentPostBuildSwaggerSchema"]
	if !ok {
		t.Logf("definitions: %#v", db.Definitions)
		t.Fail()
	}
	t.Logf("sc: %#v", sc)
	if got, want := sc.Description, "parent's description"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	t.Log(sc.Description)

	sc, ok = db.Definitions["restfulspec.childPostBuildSwaggerSchema"]
	if !ok {
		t.Logf("definitions: %#v", db.Definitions)
		t.Fail()
	}
	t.Logf("sc: %#v", sc)
	if got, want := sc.Description, "child's description"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	t.Log(sc.Description)
}
