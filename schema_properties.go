package JSONParse

import (
	"regexp"
	"strconv"
)

// 5.4.4.  additionalProperties, properties and patternProperties
// 
// 5.4.4.1.  Valid values
// 
// The value of "additionalProperties" MUST be a boolean or an object. 
// If it is an object, it MUST also be a valid JSON Schema.
// 
// The value of "properties" MUST be an object. Each value of this object 
// MUST be an object, and each object MUST be a valid JSON Schema.
// 
// The value of "patternProperties" MUST be an object. Each property name of 
// this object SHOULD be a valid regular expression, according to the ECMA 262
// regular expression dialect. Each property value of this object MUST be an 
// object, and each object MUST be a valid JSON Schema.
// 
// 5.4.4.2.  Conditions for successful validation
// 
// Successful validation of an object instance against these three keywords 
// depends on the value of "additionalProperties":
// 
// if its value is boolean true or a schema, validation succeeds;
// if its value is boolean false, the algorithm to determine validation 
//   success is described below.
// 
// 5.4.4.3.  Default values
// 
// If either "properties" or "patternProperties" are absent, they can be 
// considered present with an empty object as a value.
// 
// If "additionalProperties" is absent, it may be considered present with 
// an empty schema as a value.
// 
// 5.4.4.4.  If "additionalProperties" has boolean value false
// 
// In this case, validation of the instance depends on the property set of
// "properties" and "patternProperties". In this section, the property names
//  of "patternProperties" will be called regexes for convenience.
// 
// The first step is to collect the following sets:
// 
//  s - The property set of the instance to validate.
//  p - The property set from "properties".
//  pp - The property set from "patternProperties".
// Having collected these three sets, the process is as follows:
// 
// remove from "s" all elements of "p", if any;
// 
// for each regex in "pp", remove all elements of "s" which this regex matches.
// 
// Validation of the instance succeeds if, after these two steps, set "s" is empty.
// 
func validProperties(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	doc := mem
	if mem.GetType() == "member"  {
		doc = mem.GetValue().(*JSONNode)
	}

	if doc.GetState() == NODE_MUTEX {
		return true
	} else  {
		doc.SetState(NODE_MUTEX)
	}

	var item	*JSONNode
	found := false

	if item, found = parent.Find("properties"); found {
		item = item.GetValue().(*JSONNode)
	}

	var addtlProps		*JSONNode
	allowAddtl := false

	if addtlProps, allowAddtl = parent.Find("additionalProperties"); allowAddtl {
		if addtlProps.GetType() == "member" {
			addtlProps = addtlProps.GetValue().(*JSONNode)
		} else {
			if addtlProps.GetMemberType() == "boolean" {
				allowAddtl = addtlProps.GetValue().(bool)
			}
		}
	}

	hasPatterns, patterns := allowPatterns(parent)
 
	doc.ResetIterate()
	for {
		var key		string
		var mem		*JSONNode

		if key, mem = doc.GetNextMember(); mem == nil {
			break;
		}

		Trace.Println("      Match ", key)
		var schemaObj	*JSONNode
		match := false
		if item != nil {
			if schemaMem, found := item.Find(key); found {
				schemaObj = schemaMem.GetValue().(*JSONNode)
				match = true
			}
		}

		if !match && hasPatterns {
			for pattern, node := range patterns {
				regPattern := regexp.MustCompile(pattern)
				if match = regPattern.MatchString(key); match {
					schemaObj = node.GetValue().(*JSONNode)
					break;
				}
			}
		}

		if !match && allowAddtl {
			schemaObj = addtlProps
			match = true
		}

		if match {
			validMember(key, mem, schemaObj)
		} else if addtlProps != nil {
			Warning.Println("      member: ", key, " not found") 
		}
	}

	return true
}

// 5.4.1.  maxProperties
// 
// 5.4.1.1.  Valid values
// 
// The value of this keyword MUST be an integer. This integer MUST be greater than, or equal to, 0.
// 
// 5.4.1.2.  Conditions for successful validation
// 
// An object instance is valid against "maxProperties" if its number of properties is less than, or equal to, the value of this keyword.
// 
func validMaxProperties(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	value := mem.GetValue().(*JSONNode)

	propCount := value.GetMemberCount()
	maxCount := schema.GetValue().(int)

	Trace.Println("    max properties: ", maxCount, " mem count: ", propCount)

	return propCount <= maxCount
}

// 5.4.2.  minProperties
// 
// 5.4.2.1.  Valid values
// 
// The value of this keyword MUST be an integer. This integer MUST be greater than, or equal to, 0.
// 
// 5.4.2.2.  Conditions for successful validation
// 
// An object instance is valid against "minProperties" if its number of properties is greater than, or equal to, the value of this keyword.
// 
// 5.4.2.3.  Default value
// 
// If this keyword is not present, it may be considered present with a value of 0.
// 
func validMinProperties(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	value := mem.GetValue().(*JSONNode)

	propCount := value.GetMemberCount()
	minCount, _ := strconv.Atoi(schema.GetValue().(string))

	Trace.Println("    min properties: ", minCount, " mem count: ", propCount)

	return propCount >= minCount
}
