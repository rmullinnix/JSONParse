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
func validRequired(mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	match := false
	missing := ""
	arr := schema

	Trace.Println("  validRequired()")
	if mem.GetValueType() == V_OBJECT {
		mem.ResetIterate()
		mem = mem.GetNext()
	}

//	if arr.GetType() == N_ARRAY  {
		arr.ResetIterate()
		for {
			val := arr.GetNextValue()
			if val == nil {
				break
			}
			Trace.Println(" req element:", val)

			if _, found := mem.Find(val.(string)); found  {
				match = true
				break
			} else {
				missing = val.(string)
			}
		}
//	}

	if !match {
		errs.Add(mem, "document missing required property " + missing, JP_ERROR)
	}

	return match
}

