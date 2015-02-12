package JSONParse

import (
	"fmt"
)

func (js *JSONSchema) validObject(doc *JSONNode, schema *JSONNode) bool {

	fmt.Println("  validate object")
	js.validMembers(doc, schema)

	if required, found := schema.Find("required"); found {
		fmt.Println("  validate required members")
		js.requiredMembers(doc, required)
	}

	return true
}
