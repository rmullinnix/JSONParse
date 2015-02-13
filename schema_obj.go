package JSONParse

import (
)

func (js *JSONSchema) validObject(doc *JSONNode, schema *JSONNode) bool {
	var item	*JSONNode
	var found	bool

	Trace.Println(" Entering validObject")

	if item, found = schema.Find("properties"); found {
		item = item.GetValue().(*JSONNode)
	} else {
		Error.Panicln("invalid schema: no properties: definition")
	}

	validMember("root", doc, schema)

	return true
}
