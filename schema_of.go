package JSONParse

import (
)

// 5.5.3.  allOf
// 
// 5.5.3.1.  Valid values
// 
// This keyword's value MUST be an array. This array MUST have at least one element.
// 
// Elements of the array MUST be objects. Each object MUST be a valid JSON Schema.
// 
// 5.5.3.2.  Conditions for successful validation
// 
// An instance validates successfully against this keyword if it validates successfully against all schemas defined by this keyword's value.
//
func validAllOf(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	return true
}

// 5.5.4.  anyOf
// 
// 5.5.4.1.  Valid values
// 
// This keyword's value MUST be an array. This array MUST have at least one element.
// 
// Elements of the array MUST be objects. Each object MUST be a valid JSON Schema.
// 
// 5.5.4.2.  Conditions for successful validation
// 
// An instance validates successfully against this keyword if it validates successfully against at least one schema defined by this keyword's value.
// 
func validAnyOf(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	return true
}

// 5.5.5.  oneOf
// 
// 5.5.5.1.  Valid values
// 
// This keyword's value MUST be an array. This array MUST have at least one element.
// 
// Elements of the array MUST be objects. Each object MUST be a valid JSON Schema.
// 
//  5.5.5.2.  Conditions for successful validation
//
// An instance validates successfully against this keyword if it validates successfully against exactly one schema defined by this keyword's value.
//
func validOneOf(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	return true
}

