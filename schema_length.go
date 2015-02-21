package JSONParse

import (
//	"unicode/utf8"
	"strconv"
	"strings"
)
// 5.2.1.  maxLength
// 
// 5.2.1.1.  Valid values
// 
// The value of this keyword MUST be an integer. This integer MUST be greater than, or equal to, 0.
// 
// 5.2.1.2.  Conditions for successful validation
// 
// A string instance is valid against this keyword if its length is less than, or equal to, the value of this keyword.
// 
// The length of a string instance is defined as the number of its characters as defined by RFC 4627 [RFC4627].
// 
func validMaxLength(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validMaxLength")
	if mem.GetValueType() != V_STRING {
		Trace.Println("maxLength against non-string")
		Trace.Println(stack_id, "validMaxLength", true)
		return true
	}

	docLen := 0
	docStr := mem.GetValue().(string)
	if count := strings.Count(docStr, "\\u"); count > 0 {
		// need to figure out how to encode
		docLen = unicodeLen(docStr)
	} else {
		docLen = len(docStr)
	}

	strMax := schema.GetValue().(string)

	Trace.Println("  validMaxLength() - compare doc len ", docStr, " to shcema len ", strMax)
	maxLen, err := strconv.Atoi(strMax)
	if err != nil {
		errs.Add(mem, "Invalid integer specified in schema: " + strMax, JP_WARNING)
		Trace.Println(stack_id, "validMaxLength", false)
		return false
	}

	if docLen > maxLen {
		errs.Add(mem, "String <" + docStr + "> with length " + strconv.Itoa(docLen) + " is greater than maxLength of " + strMax, JP_ERROR)
		Trace.Println(stack_id, "validMaxLength", false)
		return false
	}

	Trace.Println(stack_id, "validMaxLength", true)
	return true
}

// 5.2.2.  minLength
// 
// 5.2.2.1.  Valid values
// 
// The value of this keyword MUST be an integer. This integer MUST be greater than, or equal to, 0.
// 
// 5.2.2.2.  Conditions for successful validation
// 
// A string instance is valid against this keyword if its length is greater than, or equal to, the value of this keyword.
// 
// The length of a string instance is defined as the number of its characters as defined by RFC 4627 [RFC4627].
// 
func validMinLength(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validMinLength")
	if mem.GetValueType() != V_STRING {
		Trace.Println("minLength against non-string")
		Trace.Println(stack_id, "validMinLength", true)
		return true
	}

	docLen := 0
	docStr := mem.GetValue().(string)
	if count := strings.Count(docStr, "\\u"); count > 0 {
		// need to figure out how to encode
		docLen = unicodeLen(docStr)
	} else {
		docLen = len(docStr)
	}

	strMin := schema.GetValue().(string)

	Trace.Println("  validMinLength() - compare doc len ", docStr, " to shcema len ", strMin)
	minLen, err := strconv.Atoi(strMin)
	if err != nil {
		errs.Add(mem, "Invalid integer specified in schema: " + strMin, JP_WARNING)
		Trace.Println(stack_id, "validMinLength", false)
		return false
	}

	if docLen < minLen {
		errs.Add(mem, "String <" + docStr + "> is less than minLength of " + strMin, JP_ERROR)
		Trace.Println(stack_id, "validMinLength", false)
		return false
	}

	Trace.Println(stack_id, "validMinLength", true)
	return true
}
