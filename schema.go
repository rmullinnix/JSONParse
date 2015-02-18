package JSONParse

import (
)

type JSONSchema struct {
	schema		*JSONParser
	doc		*JSONParser
}

// validate function signature used to add keyword validators
// as keywords are encountered, the validator function is called
// to validate the document section based on the keyword
type validator func(*JSONNode, *JSONNode, *JSONNode, *SchemaErrors) bool

// a list of keywords and associated validators
//   todo:  add func AddKeywordValidator
var keywords		map[string]validator
var schemaErrors	*SchemaErrors


//  == from the json schema core spec ==
// A JSON Schema is a JSON document, and that document MUST be an object. 
// Object members (or properties) defined by JSON Schema (this specification,
// or related specifications) are called keywords, or schema keywords.
// A JSON Schema MAY contain properties which are not schema keywords.

// Initializes a new json schema object used to parse the json schema and
// use it validate a json document agains that schema
func NewJSONSchema(source string, level string) *JSONSchema {
	js := new(JSONSchema)

	schemaErrors = NewSchemaErrors()

	keywords = make(map[string]validator)

	keywords["type"] = validType
	keywords["enum"] = validEnum
	keywords["required"] = validRequired
	keywords["properties"] = validProperties
	keywords["additionalProperties"] = validProperties
	keywords["patternProperties"] = validProperties
	keywords["minProperties"] = validProperties
	keywords["maxProperties"] = validProperties
	keywords["items"] = validItems
	keywords["uniqueItems"] = validUnique
	keywords["maxItems"] = validMaxItems
	keywords["minItems"] = validMinItems
	keywords["additionalItems"] = validAdditionalItems
	keywords["maxLength"] = validMaxLength
	keywords["minLength"] = validMinLength
	keywords["maximum"] = validMaximum
	keywords["exclusiveMaximum"] = validMaximum
	keywords["minimum"] = validMinimum
	keywords["exclusiveMinimum"] = validMinimum
	keywords["pattern"] = validPattern
	keywords["anyOf"] = validAnyOf
	keywords["allOf"] = validAllOf
	keywords["oneOf"] = validOneOf
	keywords["multipleOf"] = validMultipleOf
	keywords["default"] = validDefault

	js.schema = NewJSONParser(source, 1, level)

	js.schema.Parse()

	return js
}

// validates a document against the schema
// JSONParser is used to parse the document
func (js *JSONSchema) ValidateDocument(source string) (bool, *SchemaErrors) {

	jp := NewJSONParser(source, 10, "default")
	jp.Parse()

	js.doc = jp

	Trace.Println("Validate Document: ", source)

	result := js.validObject(jp.jsonDoc, js.schema.jsonDoc)

	schemaErrors.Output()
	return result, schemaErrors
}
