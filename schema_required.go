package JSONParse

import (
	"fmt"
)

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
		arr.dump()
		fmt.Println("        missing required property ", missing)
	}

	return match
}

