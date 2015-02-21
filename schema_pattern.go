package JSONParse

import (
	"regexp"
	"strings"
)

func allowPatterns(mem *JSONNode) (bool, map[string]*JSONNode) {
	var patterns	map[string]*JSONNode
	var item	*JSONNode
	var found	bool

	if item, found = mem.Find("patternProperties"); !found {
		Trace.Println("  allowPatterns() - no")
		return false, patterns
	}
		
	item.ResetIterate()
	item = item.GetNext()
	patterns = make(map[string]*JSONNode)

	item.ResetIterate()

	for {
		key, item := item.GetNextMember(true)
		if item == nil {
			break
		}

		patterns[key] = item
	}
	
	Trace.Println("  allowPatterns() - ", patterns)
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
func validPattern(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validPattern")

	if mem.GetValueType() != V_STRING {
		Trace.Println(stack_id, "validPattern", true)
		return true
	}

	value := schema
	pattern := value.GetValue().(string)

	key := mem.GetValue().(string)

	pattern = strings.Replace(pattern, `\\`, `\`, -1)
	patternReg := regexp.MustCompile(pattern)
	if match := patternReg.MatchString(key); match {
		Trace.Println("    match pattern <" + pattern + "> against " + key + " - true")
		Trace.Println(stack_id, "validPattern", true)
		return true
	}

	Trace.Println("    match pattern <" + pattern + "> against " + key + " - false")

	Trace.Println(stack_id, "validPattern", false)
	return false
}

