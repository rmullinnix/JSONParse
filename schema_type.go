package JSONParse

import (
	"fmt"
)

func validType(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	schemaValue := schema.GetValue().(string)
	value := mem.GetType()
	if value == "member" {
		value = mem.GetMemberType()
	}

	if value == schemaValue {
		return true
	} else {
		fmt.Println("        invalid type: expecting - ", schemaValue, " - found - ", value, "instead")
		return false
	}
}
