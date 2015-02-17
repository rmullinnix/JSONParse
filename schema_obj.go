package JSONParse

import (
)

func (js *JSONSchema) validObject(doc *JSONNode, schema *JSONNode) bool {
	var item	*JSONNode
	var found	bool

	Trace.Println(" Entering validObject")

	schema.dump()
	if item, found = schema.Find("properties"); found {
		item.ResetIterate()
		item = item.GetNext()
	} else {
		Error.Panicln("invalid schema: no properties: definition")
	}

	validMember("root", doc, schema)

	return true
}
