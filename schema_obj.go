package JSONParse

import (
)

func (js *JSONSchema) validObject(doc *JSONNode, schema *JSONNode) bool {

	Trace.Println(" Entering validObject")

	depth = 0
	doc.ResetIterate()
	schema.ResetIterate()
	validMember("root", doc, schema)

	return true
}
