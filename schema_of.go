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
	Trace.Println("    allOf")
	
	mem.SetState(NODE_SEMAPHORE)

	node := schema.GetValue().(*JSONNode)
	valid := false
	node.ResetIterate()
	for {
		item := node.GetNext()
		if item == nil {
			break
		}

		item = item.GetValue().(*JSONNode)
		valid = validMember("allOf", mem, item)

		Trace.Println("   allOf valid", valid)
		if !valid {
			break
		}
	}

	if !valid {
		mem.SetState(INVALID)
		OutputError(mem, "Did not match all the allOf constraints")
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
func validAnyOf(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	Trace.Println("    anyOf")
	
	mem.SetState(NODE_SEMAPHORE)

	node := schema.GetValue().(*JSONNode)
	valid := false
	node.ResetIterate()
	for {
		item := node.GetNext()
		if item == nil {
			break
		}

		item = item.GetValue().(*JSONNode)
		valid = validMember("anyOf", mem, item)

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
func validOneOf(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	Trace.Println("    oneOf")

	mem.SetState(NODE_SEMAPHORE)

	node := schema.GetValue().(*JSONNode)
	match := 0
	node.ResetIterate()
	for {
		item := node.GetNext()
		if item == nil {
			break
		}

		item = item.GetValue().(*JSONNode)

		valid := validMember("oneOf", mem, item)

		Trace.Println("   oneOf valid", valid)
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
			OutputError(mem, "Matched more than one in a oneOf section")
		} else {
			OutputError(mem, "Failed to match one in a oneOf section")
		}
		return false
	}
}
