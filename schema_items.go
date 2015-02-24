package JSONParse

import (
	"strconv"
//	"strings"
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
func validItems(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	arrNode := mem

	Trace.Println(stack_id, "vaildItems")

	// ignore non-arrays
	if mem.GetValueType() != V_ARRAY {
		Trace.Println("  items section is not an array")
		return true 
	}

	iterateItems := false
	if schema.GetValueType() == V_ARRAY {
		if schema.GetCount() != 1 {
			if schema.GetCount() != arrNode.GetCount() {
				Trace.Println("  items count does not match member count")
				Trace.Println(stack_id, "vaildItems", false)
				return false
			}
			schema.ResetIterate()
			iterateItems = true
		}
	} else {
		schema.ResetIterate()
		schema = schema.GetNext()
	}

	arrNode.ResetIterate()
	schem_item := schema
	index := 1
	for {
		item := arrNode.GetNext()
		if item == nil {
			break
		}

		if iterateItems {
			schem_item = schema.GetNext()
			schem_item.ResetIterate()
			schem_item = schem_item.GetNext()
		}

		new_stack_id := stack_id + "." + strconv.Itoa(index)
		index++
		valid := validMember(new_stack_id, "items", item, schem_item)

		if !valid {
			Trace.Println(stack_id, "vaildItems", false)
			return false
		}
	}
	
	Trace.Println(stack_id, "vaildItems", true)
	return true
}

func validAdditionalItems(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validAddtionalItems")
//	if _, found := mem.Find("items"); !found {
//		return true
//	}

	if schema.GetValueType() == V_OBJECT {
		schema.ResetIterate()
		schema_itm := schema.GetNext()

		mem.ResetIterate()
		valid := true
		index := 1
		for {
			mem_itm := mem.GetNext()
			if mem_itm == nil {
				break
			}

			if mem_itm.GetValueType() != V_NULL {
				new_stack_id := stack_id + "." + strconv.Itoa(index)
				index++
				nextValid := validMember(new_stack_id, "additionalItems", mem_itm, schema_itm)
				valid = valid && nextValid
			}
		}

		Trace.Println(stack_id, "validAddtionalItems", valid)
		return valid

	} else if schema.GetValueType() == V_BOOLEAN {
		if addtlItems := schema.GetValue().(bool); addtlItems {
			Trace.Println(stack_id, "validAddtionalItems", true)
			return true
		}

		// no idea if this is the  proper validation
		items, found := parent.Find("items")
		schemaCount := 0
		if found {
			if items.GetValueType() == V_ARRAY {
				schemaCount = items.GetCount()
			} else if items.GetValueType() == V_OBJECT {
				Trace.Println(stack_id, "validAddtionalItems", true)
				return true
			}
		} else {
			errs.Add(mem, "additionalItems specified but no items member in schema", JP_ERROR)
		}

		memCount := 0
		if found {
			if mem.GetValueType() == V_ARRAY {
				memCount = mem.GetCount()
			} else if items.GetType() == N_OBJECT {
				memCount = mem.GetMemberCount()
			}
		} else {
			errs.Add(mem, "additionalItems specified but no items member in document", JP_ERROR)
		}

		Trace.Println("memCount", memCount, "schemaCount", schemaCount)
		Trace.Println(stack_id, "validAddtionalItems", memCount <= schemaCount)
		return memCount <= schemaCount
	} else {
		Trace.Println("Unforeseen additionalItems type", schema.GetValueType())
		Trace.Println(stack_id, "validAddtionalItems", false)
		return false
	}
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
func validMaxItems(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validMaxItems")
	if mem.GetValueType() != V_ARRAY {
		Trace.Println("minItems() - items is not an array")
		Trace.Println(stack_id, "validMaxItems", true)
		return true
	}

	itemCount := mem.GetCount()

	strMax := schema.GetValue().(string)
	maxCount, err := strconv.Atoi(strMax)
	if err != nil {
		errs.Add(mem, "Invalid number for maxItems in Schema <" + strMax + ">", JP_WARNING)
	}

	Trace.Println("validMaxItems: max items: ", maxCount, " item count: ", itemCount)

	if itemCount > maxCount {
		errs.Add(mem, "Number of items provided is larger than maxItems <" + strMax + ">", JP_ERROR)
		Trace.Println(stack_id, "validMaxItems", false)
		return false
	}

	Trace.Println(stack_id, "validMaxItems", true)
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
func validMinItems(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validMinItems")
	if mem.GetValueType() != V_ARRAY {
		Trace.Println("minItems() - items is not an array")
		Trace.Println(stack_id, "validMinItems", true)
		return true
	}

	itemCount := mem.GetCount()

	strMin := schema.GetValue().(string)
	minCount, err := strconv.Atoi(strMin)
	if err != nil {
		errs.Add(mem, "Invalid number for minItems in Schema <" + strMin + ">", JP_WARNING)
	}

	Trace.Println("validMinItems: min items: ", minCount, " item count: ", itemCount)

	if itemCount < minCount {
		errs.Add(mem, "Number of items provided is less than minItems <" + strMin + ">", JP_ERROR)
		Trace.Println(stack_id, "validMinItems", false)
		return false
	}

	Trace.Println(stack_id, "validMinItems", true)
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
func validUnique(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validUnique")
	arr := mem

	duplicate := ""
	unique := schema.GetValue().(bool)
	if unique == false {
		Trace.Println(stack_id, "validUnique", true)
		return true
	}

	arr.ResetIterate()
	dups := make(map[string]int)
	for {
		val := arr.GetNext()
		if val == nil {
			break
		}
		token := ""

		if val.GetValueType() == V_OBJECT || val.GetValueType() == V_ARRAY {
			token = tokenizeObject(val)
		} else if val.GetValueType() == V_STRING {
			token = val.GetValue().(string)
		} else if val.GetValueType() == V_NUMBER {
			token = val.GetValue().(string)
			fVal, _ := strconv.ParseFloat(token, 64)
			token = strconv.FormatFloat(fVal, 'e', 6, 64)
		} else if val.GetValueType() == V_BOOLEAN {
			token = strconv.FormatBool(val.GetValue().(bool))
		} else if val.GetValueType() == V_NULL {
			token = "null"
		}

		if _, found := dups[token]; found {
			duplicate = token
			break
		} else {
			dups[token] = 5
		}
	}

	if len(duplicate) > 0 {
		errs.Add(mem, "Non unique items: document contains duplicate item " + duplicate, JP_ERROR)
	}

	Trace.Println(stack_id, "validUnique", len(duplicate) == 0)
	return len(duplicate) == 0
}

// poor mans attempt to tokenize an array object so it can be compared for uniqueness
// arrays of objects with nested objects will break this
//  todo: determine better way to perform object comparison
func tokenizeObject(obj *JSONNode) string {
	return obj.GetJson()
}
