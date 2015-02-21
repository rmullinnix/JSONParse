package JSONParse

import (
	"strings"
)

// 5.5.2.  type
// 
// 5.5.2.1.  Valid values
// 
// The value of this keyword MUST be either a string or an array. If it is an array, elements of the array MUST be strings and MUST be unique.
// 
// String values MUST be one of the seven primitive types defined by the core specification.
// 
// 5.5.2.2.  Conditions for successful validation
// 
// An instance matches successfully if its primitive type is one of the types defined by keyword. Recall: "number" includes "integer".
//
// === validates that the type specified in the document is the same as specified in the schema
//  todo: if number, need to check if integer
//
func validType(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validType")
	valid := false

	value := valueTypeToString(mem.GetValueType())
	schemaValue := ""
	if schema.GetValueType() == V_STRING {
		schemaValue = schema.GetValue().(string)
		valid = checkType(mem, value, schemaValue)
	} else if schema.GetValueType() == V_ARRAY {
		schema.ResetIterate()
		for {
			schema_itm := schema.GetNext()
			if schema_itm == nil {
				break
			}
			schemaValue = schema_itm.GetValue().(string)
	
			if valid = checkType(mem, value, schemaValue); valid {
				break
			}
		}
	}
	if valid {
		Trace.Println("  validType() - match on ", value)
	} else {
		Trace.Println("  validType() - invalid match on ", schemaValue, " -- was ", value)
		errs.Add(mem, "invalid type: expecting - " + schemaValue + " - found - " + value + " instead", JP_ERROR)
	}

	Trace.Println(stack_id, "validType", valid)
	return valid
}

func checkType(mem *JSONNode, value string, schema_val string) bool {
	if schema_val == "integer" && value == "number" {
		strVal := mem.GetValue().(string)
		if strings.Index(strVal, ".") == -1 {
			value = "integer"
		}
	}

	return value == schema_val
}

func valueTypeToString(valType ValueType) string {
	if valType == V_OBJECT {
		return "object"
	} else if valType == V_STRING {
		return "string"
	} else if valType == V_NUMBER {
		return "number"
	} else if valType == V_BOOLEAN {
		return "boolean"
	} else if valType == V_ARRAY {
		return "array"
	} else if valType == V_NULL {
		return "null"
	} else {
		return "unknown"
	}
}
