package JSONParse

import (
	"fmt"
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
	var value		*JSONNode
	if mem.GetMemberType() == "object" {
		value = mem.GetValue().(*JSONNode)
	} else {	
		value = mem
	}

	match := true
	missing := ""
	arr := schema.GetValue().(*JSONNode)
	if arr.GetType() == "array"  {
		arr.ResetIterate()
		for {
			val := arr.GetNextValue()
			if val == nil {
				break
			}

			if _, found := value.Find(val.(string)); !found  {
				match = false
				missing = val.(string)
				break
			}
		}
	}

	if !match {
		// need to implement error stack
		fmt.Println("        missing required property ", missing)
	}

	return match
}

