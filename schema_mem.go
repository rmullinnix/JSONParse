package JSONParse

import (
)

var depth int

func validMember(name string, mem *JSONNode, schemaMem *JSONNode, suppress bool) bool {
	nodeState := mem.GetState()
//	if !(nodeState == VIRGIN || nodeState == NODE_SEMAPHORE) {
//		return nodeState == VALID
//	}
		
	depth++
	mem.SetState(VALIDATE_IN_PROGRESS)
	if mem.GetType() == N_MEMBER {
		name = mem.name
	}

	valid := true

//	if schemaMem.GetType() == N_OBJECT || schemaMem.GetType() == N_MEMBER {
		schemaMem.ResetIterate()
		for {
			key, item := schemaMem.GetNextMember()
			if item == nil {
				break;
			}

			item.SetState(VALIDATE_IN_PROGRESS)
			if validator, found := keywords[key]; found {
				Trace.Println("  validMember()", depth, " ", name, " against schema mem ", key)
				if suppress {
					errs := NewSchemaErrors()
					valid = valid && validator(mem, item, schemaMem, errs)
				} else {
					valid = valid && validator(mem, item, schemaMem, schemaErrors)
				}
			}
		}
//	}

	depth--
	if valid {
		if nodeState != NODE_SEMAPHORE {
			mem.SetState(VALID)
		}
	} else {
		mem.SetState(INVALID)
	}

	return valid
}
