package JSONParse

import (
	"fmt"
)

// primitive types - valid in value ("type": <value>)
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

// pointers to the document node, scheme node, and schme node parent
type validator func(*JSONNode, *JSONNode, *JSONNode) bool

type keyword struct {
	keywordValidator	validator
}

var keywords		map[string]validator


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

func (js *JSONSchema) ValidateDocument(source string) (bool, []ParseError) {
	var errors	[]ParseError

	jp := NewJSONParser(source, 10)
	jp.Parse()

	js.doc = jp

	fmt.Println("Validate Document: ", source)

	result := js.validObject(jp.jsonDoc, js.schema.jsonDoc)
	return result, errors
}
