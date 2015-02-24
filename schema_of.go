package JSONParse

import (
	"strconv"
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
func validAllOf(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "allOf")
	
	valid := false
	schema.ResetIterate()
	index := 1
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

		new_stack_id := stack_id + "." + strconv.Itoa(index)
		index++
		valid = validMember(new_stack_id, "allOf", mem, of)

		Trace.Println("   allOf valid", valid)
		if !valid {
			break
		}
	}

	if !valid {
		errs.Add(mem, "Did not match all the allOf constraints", JP_ERROR)
	}

	Trace.Println(stack_id, "allOf", valid)
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
func validAnyOf(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "anyOf")
	
	valid := false
	schema.ResetIterate()
	index := 1
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

		new_stack_id := stack_id + "." + strconv.Itoa(index)
		index++
		valid = validMember(new_stack_id, "anyOf", mem, of)

		if !inSuppression {
			suppress = false
		}

		Trace.Println("   anyOf valid", valid)
		if valid {
			break
		}
	}

	Trace.Println(stack_id, "anyOf", valid)
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
// An instance validates successfully against this keyword if it validates
// successfully against exactly one schema defined by this keyword's value.
//
func validOneOf(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "oneOf")

	match := 0
	schema.ResetIterate()
	index := 1
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

		new_stack_id := stack_id + "." + strconv.Itoa(index)
		index++
		valid := validMember(new_stack_id, "oneOf", mem, of)

		if !inSuppression {
			suppress = false
		}

		Trace.Println("   oneOf", depth, name, "valid", valid, "match", match)

		if valid {
			match++
		}
	}

	if match == 1 {
		Trace.Println(stack_id, "oneOf", true)
		return true
	} else {
		if match > 1 {
			Trace.Println("  oneOf", depth, "matched more than one in oneOf")
			errs.Add(mem, "Matched more than one in a oneOf section", JP_ERROR)
		} else {
			Trace.Println("  oneOf", depth, "failed to match one in oneOf")
			errs.Add(mem, "Failed to match one in a oneOf section", JP_ERROR)
		}
		Trace.Println(stack_id, "oneOf", false)
		return false
	}
}
