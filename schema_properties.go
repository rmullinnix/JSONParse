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
func validProperties(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	doc := mem

	if doc.GetState() == NODE_MUTEX {
		return true
	} else  {
		doc.SetState(NODE_MUTEX)
	}

	if doc.GetType() == N_MEMBER || doc.GetType() == N_ARRAY {
		doc.ResetIterate()
		doc = doc.GetNext()
	}

	var item	*JSONNode
	found := false

	if item, found = parent.Find("properties"); found {
		item.ResetIterate()
		item = item.GetNext()
	}

	var addtlProps		*JSONNode
	allowAddtl := false

	if addtlProps, allowAddtl = parent.Find("additionalProperties"); allowAddtl {
		if addtlProps.GetValueType() == V_BOOLEAN {
			allowAddtl = addtlProps.GetValue().(bool)
		}
	}

	hasPatterns, patterns := allowPatterns(parent)
 
	doc.dump()
	numProps := 0
	doc.ResetIterate()
	for {
		var key		string
		var mem		*JSONNode

		if key, mem = doc.GetNextMember(); mem == nil {
			Trace.Println("  validProperties() end of members")
			break;
		}

		Trace.Println("      Match ", key)
		numProps++
		var schemaObj	*JSONNode
		match := false
		if item != nil {
			if schemaMem, found := item.Find(key); found {
				schemaMem.ResetIterate()
				schemaObj = schemaMem.GetNext()
				match = true
			}
		}

		if !match && hasPatterns {
			for pattern, node := range patterns {
				regPattern := regexp.MustCompile(pattern)
				Trace.Println("  validProperties() - match pattern <" + pattern + "> against key " + key)
				if match = regPattern.MatchString(key); match {
					node.ResetIterate()
					schemaObj = node.GetNext()
					break;
				}
			}
		}

		if !match && allowAddtl {
			schemaObj = addtlProps
			match = true
		}

		if match {
			Trace.Println("  == match successful == ")
			validMember(key, mem, schemaObj)
		} else if addtlProps != nil {
			Warning.Println("  --  member: ", key, " not found --") 
		} else {
			Trace.Println("  ++ match by addtionalProperties ++")
		}
			
	}

	if maxProps, found := parent.Find("maxProperties"); found {
		strMax := maxProps.GetValue().(string)
		maxNum, err := strconv.Atoi(strMax)
		if err != nil {
			OutputError(mem, "Non integer specefied as maxProperties in schema: " + strMax)
		}

		Trace.Println("    max properties: ", maxNum, " mem count: ", numProps)

		if numProps > maxNum {
			OutputError(mem, "Maximum <" + strMax + "> of properties exceeded")
			return false
		}
	}

	if minProps, found := parent.Find("minProperties"); found {
		strMin := minProps.GetValue().(string)
		minNum, err := strconv.Atoi(strMin)
		if err != nil {
			OutputError(mem, "Non integer specefied as minProperties in schema: " + strMin)
		}

		Trace.Println("    min properties: ", minNum, " mem count: ", numProps)

		if numProps < minNum {
			OutputError(mem, "Minimum number <" + strMin + "> of properties not supplied")
			return false
		}
	}
	return true
}

