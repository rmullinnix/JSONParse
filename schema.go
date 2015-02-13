package JSONParse

import (
	"fmt"
)

// primitive types - valid in value ("type": <value>)
// not currently being used
const (
	S_ARRAY=iota
	S_BOOLEAN
	S_INTEGER
	S_NUMBER
	S_OBJECT
	S_NULL
	S_STRING
)

type JSONSchema struct {
	schema		*JSONParser
	doc		*JSONParser
}

// validate function signature used to add keyword validators
// as keywords are encountered, the validator function is called
// to validate the document section based on the keyword
type validator func(*JSONNode, *JSONNode, *JSONNode) bool

// a list of keywords and associated validators
//   todo:  add func AddKeywordValidator
var keywords		map[string]validator


//  == from the json schema core spec ==
// A JSON Schema is a JSON document, and that document MUST be an object. 
// Object members (or properties) defined by JSON Schema (this specification,
// or related specifications) are called keywords, or schema keywords.
// A JSON Schema MAY contain properties which are not schema keywords.

// Initializes a new json schema object used to parse the json schema and
// use it validate a json document agains that schema
func NewJSONSchema(source string) *JSONSchema {
	js := new(JSONSchema)

	keywords = make(map[string]validator)

	keywords["type"] = validType
	keywords["enum"] = validEnum
	keywords["required"] = validRequired
	keywords["properties"] = validProperties
	keywords["additionalProperties"] = validProperties
	keywords["patternProperties"] = validProperties
	keywords["uniqueItems"] = validUnique
	keywords["minProperties"] = validMinProperties
	keywords["maxProperties"] = validMaxProperties
//	keywords["maxLength"] = validMaxLength
//	keywords["minLength"] = validMinLength
//	keywords["maximum"] = validMaximum
//	keywords["exclusiveMaximum"] = validExclusiveMaximum
//	keywords["minimum"] = validMinimum
//	keywords["exclusiveMinimum"] = validExclusiveMinimum
//	keywords["pattern"] = validPattern
//	keywords["anyOf"] = validAnyOf
//	keywords["allOf"] = validAllOf
//	keywords["oneOf"] = validOneOf
//	keywords["multipleOf"] = validMultipleOf
//	keywords["default"] = validDefault
//	keywords["maxItems"] = validMaxItems
//	keywords["minItems"] = validMinItems
//	keywords["maxProperties"] = validMaxProperties
//	keywords["minProperties"] = validMinProperties
//	keywords["additionalItems"] = validAdditionalItems

	js.schema = NewJSONParser(source, 1)

	js.schema.Parse()

	return js
}

// validates a document against the schema
// JSONParser is used to parse the document
func (js *JSONSchema) ValidateDocument(source string) (bool, []ParseError) {
	var errors	[]ParseError

	jp := NewJSONParser(source, 10)
	jp.Parse()

	js.doc = jp

	fmt.Println("Validate Document: ", source)

	result := js.validObject(jp.jsonDoc, js.schema.jsonDoc)
	return result, errors
}
