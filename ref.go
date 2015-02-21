package JSONParse

import (
	"strings"
)

// steps the reference map and stores the pointer to the json node
// in the json tree.
// when the schema is used for validation, the references are resolved
// as they are encounter -- see node.go
func (jp *JSONParser) resolveReferences() {
	for key, _ := range jp.references {
		ref := jp.refObject(jp.jsonDoc, key)
		if ref == nil {
			jp.addError("Unable to resolve reference " + key, JP_FATAL)
			OutputError(jp.references[key], "Invalid json reference " + key, JP_ERROR)
			jp.references[key] = nil
		} else {
			jp.references[key] = ref
		}
	}
	//Trace.Println("== REFERENCE TABLE ==")
	//Trace.Println(jp.references)
}

// find associated reference either in the current document or an
// external document.  If it is an external document, the document is
// retrieved and parsed.  The external document is stored for future
// reference during schema validation.
//
// if the reference cannot be resolved, it is flagged as an error
func (jp *JSONParser)refObject(doc *JSONNode, ref string) *JSONNode {
	// internal reference
	if ref[0:1] == "#" {
		if ref == "#" {
			return doc
		}

		subparts := strings.Split(ref[2:], "/")
		refObj := doc
		match := 0
		for i := 0; i < len(subparts); i++ {
			if item, found := refObj.Find(subparts[i]); found {
				item.ResetIterate()
				refObj = item.GetNext()
				match++
			}
		}
		if match == len(subparts) {
			return refObj
		}
	} else { // if strings.HasPrefix(ref, "http")  {
		var eDoc		*JSONNode
	
		found := false
		parts := strings.Split(ref, "#")

		if eDoc, found = jp.extDocs[parts[0]]; !found {
			extDoc := NewJSONParser(parts[0], 1, "default")

			extDoc.Parse()

			eDoc = extDoc.jsonDoc
			jp.extDocs[parts[0]] = eDoc
		}

		return jp.refObject(eDoc, "#" + parts[1])
	}

	return nil
}
