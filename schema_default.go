package JSONParse

import (
)

// 6.2.  "default"
//
// 6.2.1.  Valid values
//
// There are no restrictions placed on the value of this keyword.
//
// 6.2.2.  Purpose
//
// This keyword can be used to supply a default JSON value associated with a
// particular schema. It is RECOMMENDED that a default value be valid against
// the associated schema.
//
// This keyword MAY be used in root schemas, and in any subschemas.
// 
func validDefault(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	return true
}
