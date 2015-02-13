package JSONParse

import (
	"fmt"
	"strconv"
)

// 5.3.2.  maxItems
// 
// 5.3.2.1.  Valid values
// 
// The value of this keyword MUST be an integer. This integer MUST be greater than, or equal to, 0.
// 
// 5.3.2.2.  Conditions for successful validation
// 
// An array instance is valid against "maxItems" if its size is less than, or equal to, the value of this keyword.
// 
func validMaxItems(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	value := mem.GetValue().(*JSONNode)

	itemCount := value.GetMemberCount()

	// may have to do a find from the schema node for "items"
	maxCount, _ := strconv.Atoi(schema.GetValue().(string))

	fmt.Println("    max items: ", maxCount, " item count: ", itemCount)

	return itemCount <= maxCount
}


// 5.3.3.  minItems
// 
// 5.3.3.1.  Valid values
// 
// The value of this keyword MUST be an integer. This integer MUST be greater than, or equal to, 0.
// 
// 5.3.3.2.  Conditions for successful validation
// 
// An array instance is valid against "minItems" if its size is greater than, or equal to, the value of this keyword.
// 
// 5.3.3.3.  Default value
// 
// If this keyword is not present, it may be considered present with a value of 0.
// 
func validMinItems(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	value := mem.GetValue().(*JSONNode)

	itemCount := value.GetMemberCount()

	// may have to do a find from the schema node for "items"
	minCount, _ := strconv.Atoi(schema.GetValue().(string))

	fmt.Println("    min items: ", minCount, " item count: ", itemCount)

	return itemCount >= minCount
}

// 5.3.4.  uniqueItems
// 
// 5.3.4.1.  Valid values
// 
// The value of this keyword MUST be a boolean.
// 
// 5.3.4.2.  Conditions for successful validation
// 
// If this keyword has boolean value false, the instance validates successfully. 
// If it has boolean value true, the instance validates successfully if all of its elements are unique.
// 
// 5.3.4.3.  Default value
// 
// If not present, this keyword may be considered present with boolean value false.
//
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

// poor mans attempt to tokenize an array object so it can be compared for uniqueness
// arrays of objects with nested objects will break this
//  todo: determine better way to perform object comparison
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
