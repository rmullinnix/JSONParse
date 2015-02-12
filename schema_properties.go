package JSONParse

import (
	"fmt"
	"regexp"
)

// covers properties, additonal properties and pattern properties
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

		fmt.Println("      Match ", key)
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

		if !match && (addtlProps != nil) {
			if schemaMem, found := addtlProps.Find(key); found {
				schemaObj = schemaMem.GetValue().(*JSONNode)
				match = true
			}
		}

		if match {
			validMember(key, mem, schemaObj)
		} else if addtlProps != nil {
			fmt.Println("      member: ", key, " not found") 
		}
	}

	return true
}
