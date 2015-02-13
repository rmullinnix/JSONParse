package JSONParse

import (
)

func validMember(name string, mem *JSONNode, schemaMem *JSONNode) bool {
	nodeState := mem.GetState()
	if !(nodeState == VIRGIN || nodeState == NODE_SEMAPHORE) {
		return nodeState == VALID
	}
		
	mem.SetState(VALIDATE_IN_PROGRESS)

	valid := true

	if name == "anyOf" {
		schemaMem.dump()
	}
	Trace.Println("validate individual member: ", name, " as ", mem.GetType(), " against ", schemaMem.GetType());
	if schemaMem.GetType() == "object" {
		schemaMem.ResetIterate()
		for {
			key, item := schemaMem.GetNextMember()
			Trace.Println("  validate ", name, " against schema mem ", key)
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
		if nodeState != NODE_SEMAPHORE {
			mem.SetState(VALID)
		}
	} else {
		mem.SetState(INVALID)
	}

	return valid
}
