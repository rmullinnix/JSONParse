package JSONParse

import (
)

// 5.5.6.  not
// 
// 5.5.6.1.  Valid values
// 
// This keyword's value MUST be an object. This object MUST be a valid JSON Schema.
// 
// 5.5.6.2.  Conditions for successful validation
// 
// An instance is valid against this keyword if it fails to validate successfully against the schema defined by this keyword.
// 
func validNot(mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	schema.ResetIterate()
	items := schema.GetNext() // items is of type object
	
	Trace.Println("  vaildNot()")

	suppress = true
	valid := validMember("not", mem, items)
	suppress = false
	
	return !valid
}
