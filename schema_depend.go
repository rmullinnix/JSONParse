package JSONParse

import (
)

// 5.4.5.  dependencies
//
// 5.4.5.1.  Valid values
//
// This keyword's value MUST be an object. Each value of this object MUST be
// either an object or an array.
//
// If the value is an object, it MUST be a valid JSON Schema. This is called a
// schema dependency.
//
// If the value is an array, it MUST have at least one element. Each element
// MUST be a string, and elements in the array MUST be unique. This is called
// a property dependency.
//
// 5.4.5.2.  Conditions for successful validation
//
// 5.4.5.2.1.  Schema dependencies
//
// For all (name, schema) pair of schema dependencies, if the instance has a 
// property by this name, then it must also validate successfully against the schema.
//
// Note that this is the instance itself which must validate successfully,
// not the value associated with the property name.
//
// 5.4.5.2.2.  Property dependencies
//
// For each (name, propertyset) pair of property dependencies, if the instance
// has a property by this name, then it must also have properties with the same
// names as propertyset.
// 
func validDependencies(mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println("  validDependencies()")
	if schema.GetValueType() != V_OBJECT {
		return false
	}

	schema.ResetIterate()
	schema_itm := schema.GetNext()

	schema_itm.ResetIterate()
	key, prop := schema_itm.GetNextMember(true)

	Trace.Println(" dep key", key)

	if mem.GetValueType() != V_OBJECT {
		return true
	}

	mem.ResetIterate()
	mem = mem.GetNext()

	_, found := mem.Find(key)
	if !found {
		return true
	}

	valid := true
	if prop.GetValueType() == V_ARRAY {
		prop.ResetIterate()
		missing := ""
		match := true
		for {
			val := prop.GetNextValue()
			if val == nil {
				Trace.Println("end list")
				break
			}
			Trace.Println(" dep element:", val)

			if _, found := mem.Find(val.(string)); !found  {
				match = false
				missing = val.(string)
				break
			} else {
				Trace.Println("  match element")
			}
		}

		if !match {
			errs.Add(mem, "document missing dependent property " + missing, JP_ERROR)
			valid = false
		} else {
			valid = true
		}
	} else if prop.GetValueType() == V_OBJECT {
		prop.ResetIterate()
		schema_sub := prop.GetNext()

		valid = validMember("dependencies", mem, schema_sub)
	}

	return valid
}
