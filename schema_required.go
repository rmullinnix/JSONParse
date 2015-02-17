package JSONParse

import (
)

// 5.4.3.  required
// 
// 5.4.3.1.  Valid values
// 
// The value of this keyword MUST be an array. This array MUST have at least one element. Elements of this array MUST be strings, and MUST be unique.
// 
// 5.4.3.2.  Conditions for successful validation
// 
// An object instance is valid against this keyword if its property set contains all elements in this keyword's array value.
// 
func validRequired(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	match := true
	missing := ""
	arr := schema

	Trace.Println("  validRequired()")
	if arr.GetType() == N_ARRAY  {
		arr.ResetIterate()
		for {
			val := arr.GetNextValue()
			if val == nil {
				break
			}

			if _, found := mem.Find(val.(string)); !found  {
				match = false
				missing = val.(string)
				break
			}
		}
	}

	if !match {
		OutputError(mem, "document missing required property " + missing)
	}

	return match
}

