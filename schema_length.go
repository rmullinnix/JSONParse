package JSONParse

import (
	"strconv"
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
func validMaxLength(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	docStr := mem.GetValue().(string)
	strMax := schema.GetValue().(string)

	maxLen, err := strconv.Atoi(strMax)
	if err != nil {
		OutputError(mem, "Invalid integer specified in schema: " + strMax)
		return false
	}

	if len(docStr) > maxLen {
		OutputError(mem, "String <" + docStr + "> is greater than maxLength of " + strMax)
		return false
	}

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
func validMinLength(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	docStr := mem.GetValue().(string)
	strMin := schema.GetValue().(string)

	minLen, err := strconv.Atoi(strMin)
	if err != nil {
		OutputError(mem, "Invalid integer specified in schema: " + strMin)
		return false
	}

	if len(docStr) < minLen {
		OutputError(mem, "String <" + docStr + "> is less than minLength of " + strMin)
		return false
	}

	return true
}
