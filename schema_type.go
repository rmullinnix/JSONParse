package JSONParse

import (
	"fmt"
)

// 5.5.2.  type
// 
// 5.5.2.1.  Valid values
// 
// The value of this keyword MUST be either a string or an array. If it is an array, elements of the array MUST be strings and MUST be unique.
// 
// String values MUST be one of the seven primitive types defined by the core specification.
// 
// 5.5.2.2.  Conditions for successful validation
// 
// An instance matches successfully if its primitive type is one of the types defined by keyword. Recall: "number" includes "integer".
//
// === validates that the type specified in the document is the same as specified in the schema
//  todo:  support validation against arry of types
//
func validType(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	schemaValue := schema.GetValue().(string)
	value := mem.GetType()
	if value == "member" {
		value = mem.GetMemberType()
	}

	if value == schemaValue {
		return true
	} else {
		fmt.Println("        invalid type: expecting - ", schemaValue, " - found - ", value, "instead")
		return false
	}
}
