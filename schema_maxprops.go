package JSONParse

import (
	"fmt"
)

func validMaxProperties(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	value := mem.GetValue().(*JSONNode)

	propCount := value.GetMemberCount()
	maxCount := schema.GetValue().(int)

	fmt.Println("    max properties: ", maxCount, " mem count: ", propCount)

	return propCount > maxCount
}
