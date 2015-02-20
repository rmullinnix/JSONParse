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

		valid = validMember("allOf", mem, item)

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

		inSuppression := false
		if suppress  {
			inSuppression = true	
		} else {
			suppress = true
		}

		valid = validMember("anyOf", mem, of)

		if !inSuppression {
			suppress = false
		}

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

	Trace.Println("   oneOf", depth)
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

		inSuppression := false
		if suppress  {
			inSuppression = true	
		} else {
			suppress = true
		}

		valid := validMember("oneOf", mem, of)

		if !inSuppression {
			suppress = false
		}

		Warning.Println("   oneOf", depth, name, "valid", valid, "match", match)

		if valid {
			match++
		}
	}

	if match == 1 {
		Warning.Println("  oneOf", depth, "matched one in oneOf")
		mem.SetState(VALID)
		return true
	} else {
		mem.SetState(INVALID)
		if match > 1 {
			Warning.Println("  oneOf", depth, "matched more than one in oneOf")
			errs.Add(mem, "Matched more than one in a oneOf section", JP_ERROR)
		} else {
			Warning.Println("  oneOf", depth, "failed to match one in oneOf")
			errs.Add(mem, "Failed to match one in a oneOf section", JP_ERROR)
		}
		return false
	}
}
