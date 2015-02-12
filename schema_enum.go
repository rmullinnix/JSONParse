package JSONParse

import (
	"fmt"
)

func validEnum(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	value := mem.GetValue().(string)

	match := false
	arr := schema.GetValue().(*JSONNode)
	if arr.GetType() == "array"  {
		arr.ResetIterate()
		for {
			val := arr.GetNextValue()
			if val == nil {
				break
			}

			if value == val.(string)  {
				match = true
				break
			}
		}
	}

	if !match {
		fmt.Println("        invalid enum ", value)
	}

	return match
}
