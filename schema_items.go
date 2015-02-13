package JSONParse

import (
	"fmt"
	"strconv"
)

func validUnique(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	arr := mem.GetValue().(*JSONNode)

	duplicate := ""
	unique := schema.GetValue().(bool)
	if unique == false {
		return true
	}

	dups := make(map[string]int)
	for {
		val := arr.GetNext()
		if val == nil {
			break
		}
		token := ""

		if val.GetType() == "object" {
			token = tokenizeObject(val)
		} else if val.GetType() == "string" {
			token = val.GetValue().(string)
		}

		if _, found := dups[token]; found {
			duplicate = token
			break
		} else {
			dups[token] = 5
		}
	}

	if len(duplicate) > 0 {
		fmt.Println("        contains duplicate item ", duplicate)
	}

	return len(duplicate) == 0
}

func tokenizeObject(obj *JSONNode) string {
	token := ""
	for {
		key, item := obj.GetNextMember()
		if item == nil {
			break
		}

		token += key + item.GetValue().(string)
	}
	return token
}

func validMaxItems(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	value := mem.GetValue().(*JSONNode)

	itemCount := value.GetMemberCount()

	// may have to do a find from the schema node for "items"
	maxCount, _ := strconv.Atoi(schema.GetValue().(string))

	fmt.Println("    max items: ", maxCount, " item count: ", itemCount)

	return itemCount > maxCount
}

func validMinItems(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	value := mem.GetValue().(*JSONNode)

	itemCount := value.GetMemberCount()

	// may have to do a find from the schema node for "items"
	minCount, _ := strconv.Atoi(schema.GetValue().(string))

	fmt.Println("    min items: ", minCount, " item count: ", itemCount)

	return itemCount < minCount
}
