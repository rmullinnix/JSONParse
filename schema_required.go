package JSONParse

import (
)

// 5.4.3.  required
// 
// 5.4.3.1.  Valid values
// 
// The value of this keyword MUST be an array. This array MUST have at least one element. 
// Elements of this array MUST be strings, and MUST be unique.
// 
// 5.4.3.2.  Conditions for successful validation
// 
// An object instance is valid against this keyword if its property set contains all elements in this keyword's array value.
// 
func validRequired(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validRequired")

	missing := ""
	arr := schema

	if mem.GetType() != N_OBJECT {
		mem.ResetIterate()
		mem = mem.GetNext()
	}

	match := 0
	count := 0
//	if arr.GetType() == N_ARRAY  {
		arr.ResetIterate()
		for {
			val := arr.GetNextValue()
			if val == nil {
				break
			}
			count++
			Trace.Println(" req element:", val)

			if _, found := mem.Find(val.(string)); found  {
				match++
			} else {
				missing = val.(string)
				break
			}
		}
//	}

	if match != count {
		errs.Add(mem, "document missing required property " + missing, JP_ERROR)
	}

	Trace.Println(stack_id, "validRequired", match == count)
	return match == count
}

