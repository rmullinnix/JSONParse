package JSONParse

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

// 5.4.3.2.  Conditions for successful validation
// 
// An object instance is valid against this keyword if its property set contains all elements in this keyword's array value.
// 
func validMinLength(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	return true
}
