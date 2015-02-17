package JSONParse

import (
	"strings"
)

// 5.5.1.  enum
//
// 5.5.1.1.  Valid values
// 
// The value of this keyword MUST be an array. This array MUST have at least one element. Elements in the array MUST be unique.
// 
// Elements in the array MAY be of any type, including null.
//
// 5.5.1.2.  Conditions for successful validation
// 
// An instance validates successfully against this keyword if its value is equal to one of the elements in this keyword's array value.
//
// === this function implments 5.5.1.2 of the json schema validation spec
func validEnum(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	value := mem.GetValue().(string)

	match := false
	valStr := "["
	if schema.GetValueType() == V_ARRAY  {
		schema.ResetIterate()
		for {
			val := schema.GetNextValue()
			if val == nil {
				break
			}

			if value == val.(string)  {
				match = true
				break
			}
			valStr += val.(string) + ", "
		}
	}

	Trace.Println("  validEnum() match - ", match)
	if !match {
		valStr = strings.TrimSuffix(valStr, ", ")
		valStr += "]"
		OutputError(mem, "invalid enum <" + value + "> specified in document. Must be one of " + valStr)
	}

	return match
}
