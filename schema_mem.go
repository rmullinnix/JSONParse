package JSONParse

import (
//	"fmt"
//	"regexp"
)

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
