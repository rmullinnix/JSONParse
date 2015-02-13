package JSONParse

import (
	"regexp"
)

func allowPatterns(mem *JSONNode) (bool, map[string]*JSONNode) {
	var patterns	map[string]*JSONNode
	var item	*JSONNode
	var found	bool

	if item, found = mem.Find("patternProperties"); !found {
		return false, patterns
	}
		
	item = item.GetValue().(*JSONNode)
	patterns = make(map[string]*JSONNode)
	item.ResetIterate()

	for {
		key, item := item.GetNextMember()
		if item == nil {
			break
		}

		patterns[key] = item
	}
	
	return true, patterns
}

// 5.2.3.  pattern
//
// 5.2.3.1.  Valid values
//
// The value of this keyword MUST be a string. This string SHOULD be a valid
// regular expression, according to the ECMA 262 regular expression dialect.
//
// 5.2.3.2.  Conditions for successful validation
//
// A string instance is considered valid if the regular expression matches the
// instance successfully. Recall: regular expressions are not implicitly anchored.
//
func validPattern(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	Trace.Println("  validPattern")

	var value		*JSONNode
	if schema.GetMemberType() == "object" {
		value = schema.GetValue().(*JSONNode)
	} else {	
		value = schema
	}

	pattern := value.GetValue().(string)

	key := mem.GetValue().(string)
	Trace.Println("    match pattern <" + pattern + "> against " + key)

	patternReg := regexp.MustCompile(pattern)
	if match := patternReg.MatchString(key); match {
		return true
	}

	return false
}

