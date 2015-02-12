package JSONParse

import (
	"fmt"
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
