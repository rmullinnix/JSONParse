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

	if schemaMem.GetType() == N_OBJECT {
		schemaMem.ResetIterate()
		for {
			key, item := schemaMem.GetNextMember()
			if item == nil {
				break;
			}

			item.SetState(VALIDATE_IN_PROGRESS)
			if validator, found := keywords[key]; found {
				Trace.Println("  validMember()", name, " against schema mem ", key)
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
