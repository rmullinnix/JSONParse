package JSONParse

import (
	"fmt"
//	"regexp"
)

func (js *JSONSchema) validMembers(doc *JSONNode, schema *JSONNode) {
	var item	*JSONNode
	var found	bool

	//fmt.Println("    validate members")

	if item, found = schema.Find("properties"); found {
		item = item.GetValue().(*JSONNode)
	} else {
		panic("invalid schema")
	}

	validMember("root", doc, schema)
}

func (js *JSONSchema) requiredMembers(mem *JSONNode, reqProps *JSONNode) {
	if reqProps == nil {
		return
	}

	reqArr := reqProps.GetValue().(*JSONNode)

	reqArr.ResetIterate()
	for {
		reqItem := reqArr.GetNextValue()
		if reqItem == nil {
			break
		}
		reqValue := reqItem.(string)

		if _, found := mem.Find(reqValue); !found {
			fmt.Println("      required member missing: ", reqValue)
		} else {
			//fmt.Println("      required member found: ", reqValue)
		}
	}
}

func validMember(name string, mem *JSONNode, schemaMem *JSONNode) bool {
	nodeState := mem.GetState()
	if nodeState != VIRGIN {
		return nodeState == VALID
	}
		
	mem.SetState(VALIDATE_IN_PROGRESS)

	valid := true

	//fmt.Println("      validate individual member: ", name, " as ", mem.GetType(), " against ", schemaMem.GetType());
	if schemaMem.GetType() == "object" {
		schemaMem.ResetIterate()
		for {
			key, item := schemaMem.GetNextMember()
	//		fmt.Println("        validate ", name, " against schema mem ", key)
			if item == nil {
				break;
			}

			item.SetState(VALIDATE_IN_PROGRESS)
			if validator, found := keywords[key]; found {
				valid = validator(mem, item, schemaMem)
			}
		}
	}

	if valid {
		mem.SetState(VALID)
	} else {
		mem.SetState(INVALID)
	}

	return valid
}
