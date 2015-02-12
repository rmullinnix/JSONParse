package JSONParse

import (
//	"fmt"
	"strings"
)

func (jp *JSONParser) resolveReferences() {
	for key, _ := range jp.references {
		jp.references[key] = jp.refObject(jp.jsonDoc, key)
	}
}

func (jp *JSONParser)refObject(doc *JSONNode, ref string) *JSONNode {
	// internal reference
	if ref[0:1] == "#" {
		if ref == "#" {
			return doc
		}

		subparts := strings.Split(ref[1:], "/")
		refObj := doc
		match := 0
		for i := 0; i < len(subparts); i++ {
			if item, found := refObj.Find(subparts[i]); found {
				refObj = item.GetValue().(*JSONNode)
				match++
			}
		}
		if match == len(subparts) - 1 {
			return refObj
		}
	} else if strings.HasPrefix(ref, "http")  {
		var eDoc		*JSONNode
	
		found := false
		parts := strings.Split(ref, "#")

		if eDoc, found = jp.extDocs[parts[0]]; !found {
			extDoc := NewJSONParser(parts[0], 1)

			extDoc.Parse()

			eDoc = extDoc.jsonDoc
			jp.extDocs[parts[0]] = eDoc
		}

		return jp.refObject(eDoc, "#" + parts[1])
	}
	jp.addError("Unable to resolve reference " + ref, JP_FATAL)

	return nil
}
