package JSONParse

import (
)

func (js *JSONSchema) validObject(doc *JSONNode, schema *JSONNode) bool {
	var item	*JSONNode
	var found	bool

	Trace.Println(" Entering validObject")

	if item, found = schema.Find("properties"); found {
		item.ResetIterate()
		item = item.GetNext()
	} else {
		Error.Panicln("invalid schema: no properties: definition")
	}

	depth = 0
	validMember("root", doc, schema, false)

	return true
}
