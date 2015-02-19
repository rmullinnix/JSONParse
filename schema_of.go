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
func validAllOf(mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println("  allOf()")
	
	mem.SetState(NODE_SEMAPHORE)

	valid := false
	schema.ResetIterate()
	for {
		item := schema.GetNext()
		if item == nil {
			break
		}

		if item.GetValueType () == V_OBJECT {
			item.ResetIterate()
			item = item.GetNext()
		}

		valid = validMember("allOf", mem, item, false)

		Trace.Println("   allOf valid", valid)
		if !valid {
			break
		}
	}

	if !valid {
		mem.SetState(INVALID)
		errs.Add(mem, "Did not match all the allOf constraints", JP_ERROR)
	} else {
		mem.SetState(VALID)
	}

	return valid
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
func validAnyOf(mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println("  anyOf()")
	
	mem.SetState(NODE_SEMAPHORE)

	valid := false
	schema.ResetIterate()
	for {
		item := schema.GetNext()
		if item == nil {
			break
		}

		of := item
		if of.GetValueType () == V_OBJECT {
			of.ResetIterate()
			of = of.GetNext()
		}

		valid = validMember("anyOf", mem, of, true)

		Trace.Println("   anyOf valid", valid)
		if valid {
			break
		}
	}

	return valid
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
func validOneOf(mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println("  oneOf()")

	mem.SetState(NODE_SEMAPHORE)

	match := 0
	schema.ResetIterate()
	for {
		item := schema.GetNext()
		if item == nil {
			break
		}

		of := item
		if of.GetValueType () == V_OBJECT {
			of.ResetIterate()
			of = of.GetNext()
		}

		var name	string
		ref, found := of.namedKids["$ref"]
		if found {
			name = ref.GetValue().(string)
		}

		valid := validMember("oneOf", mem, of, true)

		Warning.Println("   oneOf ", name, "valid", valid)

		if valid {
			match++
		}
	}

	if match == 1 {
		mem.SetState(VALID)
		return true
	} else {
		mem.SetState(INVALID)
		if match > 1 {
			errs.Add(mem, "Matched more than one in a oneOf section", JP_ERROR)
		} else {
			errs.Add(mem, "Failed to match one in a oneOf section", JP_ERROR)
		}
		return false
	}
}
