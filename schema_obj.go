package JSONParse

import (
//	"fmt"
)

func (js *JSONSchema) validObject(doc *JSONNode, schema *JSONNode) bool {
	var item	*JSONNode
	var found	bool

	//fmt.Println("    validate members")

	if item, found = schema.Find("properties"); found {
		item = item.GetValue().(*JSONNode)
	} else {
		panic("invalid schema")
	}

	validMember("root", doc, schema)

	return true
}
