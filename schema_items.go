package JSONParse

import (
	"strconv"
)

// 5.3.1.  additionalItems and items
//
// 5.3.1.1.  Valid values
//
// The value of "additionalItems" MUST be either a boolean or an object. If it is
// an object, this object MUST be a valid JSON Schema.
//
// The value of "items" MUST be either an object or an array. If it is an object,
// this object MUST be a valid JSON Schema. If it is an array, items of this array
// MUST be objects, and each of these objects MUST be a valid JSON Schema.
//
// 5.3.1.2.  Conditions for successful validation
//
// Successful validation of an array instance with regards to these two keywords
// is determined as follows:
//
//   if "items" is not present, or its value is an object, validation of the instance
//     always succeeds, regardless of the value of "additionalItems";
//   if the value of "additionalItems" is boolean value true or an object, validation
//     of the instance always succeeds;
//   if the value of "additionalItems" is boolean value false and the value of "items"
//     is an array, the instance is valid if its size is less than, or equal to, the
//     size of "items".
//
func validItems(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	schema.ResetIterate()
	items := schema.GetNext() // items is of type object
	
	arrNode := mem

	Trace.Println("  vaildItems()")
	arrNode.ResetIterate()
	for {
		item := arrNode.GetNext()
		if item == nil {
			break
		}

		Trace.Println("  items: call validMember()")
		valid := validMember("items", item, items)
		if !valid {
			return false
		}
	}
	
	return true
}

func validAdditionalItems(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	Trace.Println("  validAddtionalItems()")
	if _, found := mem.Find("items"); !found {
		return true
	}

	addtlItems := schema.GetValue().(bool)
	if addtlItems {
		return true
	}

	// no idea if this is the  proper validation
	items, found := parent.Find("items")
	schemaCount := 0
	if found {
		if items.GetType() == N_ARRAY {
			schemaCount = items.GetCount()
		} else if items.GetType() == N_OBJECT {
			schemaCount = items.GetMemberCount()
		}
	} else {
		OutputError(mem, "additionItems specified but no items member in schema")
	}

	memCount := 0
	items, found = mem.Find("items")
	if found {
		if items.GetType() == N_ARRAY {
			memCount = items.GetCount()
		} else if items.GetType() == N_OBJECT {
			memCount = items.GetMemberCount()
		}
	} else {
		OutputError(mem, "additionItems specified but no items member in document")
	}

	return memCount <= schemaCount
}


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
	itemCount := mem.GetCount()

	strMax := schema.GetValue().(string)
	maxCount, err := strconv.Atoi(strMax)
	if err != nil {
		OutputError(mem, "Invalid number for maxItems in Schema <" + strMax + ">")
	}

	Trace.Println("validMaxItems: max items: ", maxCount, " item count: ", itemCount)

	if itemCount > maxCount {
		OutputError(mem, "Number of items provided is larger than maxItems <" + strMax + ">")
		return false
	}

	return true
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
	itemCount := mem.GetCount()

	strMin := schema.GetValue().(string)
	minCount, err := strconv.Atoi(strMin)
	if err != nil {
		OutputError(mem, "Invalid number for minItems in Schema <" + strMin + ">")
	}

	Trace.Println("validMinItems: min items: ", minCount, " item count: ", itemCount)

	if itemCount < minCount {
		OutputError(mem, "Number of items provided is less than minItems <" + strMin + ">")
		return false
	}

	return true
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
	arr := mem

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

		if val.GetValueType() == V_OBJECT {
			token = tokenizeObject(val)
		} else if val.GetValueType() == V_STRING {
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
		OutputError(mem, "Non unique items: document contains duplicate item " + duplicate)
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
