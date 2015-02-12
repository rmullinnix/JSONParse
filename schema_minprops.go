package JSONParse

import (
	"fmt"
	"strconv"
)

func validMinProperties(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	value := mem.GetValue().(*JSONNode)

	propCount := value.GetMemberCount()
	minCount, _ := strconv.Atoi(schema.GetValue().(string))

	fmt.Println("    min properties: ", minCount, " mem count: ", propCount)

	return propCount < minCount
}
