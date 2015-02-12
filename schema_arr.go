package JSONParse

import (
	"fmt"
)

func (js *JSONSchema) validArray(arr *JSONNode, schemaObj *JSONNode) bool {
	fmt.Println("    validate array")
	sType := schemaObj.GetType()
	if sType != "array" {
		fmt.Println("      not an array")
		return false
	}

	items, hasItems := schemaObj.Find("items")

	if hasItems {
		itemObj := items.GetValue().(*JSONNode)

		typeObj, found := itemObj.Find("type")
		if found {
			itemType := typeObj.GetValue().(string)
			fmt.Println("      item type: ", itemType)
			if itemType == "object" {

			} else if itemType == "string" {

			} else if itemType == "refptr" {
			} else {
			}
		}
	}

	return true
}
