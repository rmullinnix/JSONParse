package JSONParse

import (
)

// 5.5.6.  not
// 
// 5.5.6.1.  Valid values
// 
// This keyword's value MUST be an object. This object MUST be a valid JSON Schema.
// 
// 5.5.6.2.  Conditions for successful validation
// 
// An instance is valid against this keyword if it fails to validate successfully against the schema defined by this keyword.
// 
func validNot(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "vaildNot")

	schema.ResetIterate()
	items := schema.GetNext() // items is of type object

	suppress = true
	valid := validMember(stack_id, "not", mem, items)
	suppress = false
	
	if valid {
		errs.Add(mem, "Encountered valid components in NOT section", JP_ERROR)
	}

	suppressErrors = NewSchemaErrors()

	Trace.Println(stack_id, "vaildNot", !valid)
	return !valid
}
