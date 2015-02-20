package JSONParse

import (
)

var depth int

func validMember(name string, mem *JSONNode, schemaMem *JSONNode) bool {
	depth++
	Trace.Println("  validMember()", depth, name)
	nodeState := mem.GetState()
//	if !(nodeState == VIRGIN || nodeState == NODE_SEMAPHORE) {
//		return nodeState == VALID
//	}
		
	mem.SetState(VALIDATE_IN_PROGRESS)
	if mem.GetType() == N_MEMBER {
		name = mem.name
	}

	valid := true

//	if schemaMem.GetType() == N_OBJECT || schemaMem.GetType() == N_MEMBER {
		schemaMem.ResetIterate()
		for {
	
			key, item := schemaMem.GetNextMember(true)
			if item == nil {
				Trace.Println("  end of schema section")
				break
			}
			Trace.Println("  validMember()", depth, "  key:", key)

			item.SetState(VALIDATE_IN_PROGRESS)
			if validator, found := keywords[key]; found {
				Trace.Println("  validMember()", depth, "   ", name, "against schema mem", key)
				if suppress {
					errs := NewSchemaErrors()
					nextValid := validator(mem, item, schemaMem, errs)
					valid = valid && nextValid
					Trace.Println("  validMember()", depth, valid, nextValid)
				} else {
					nextValid := validator(mem, item, schemaMem, schemaErrors)
					valid = valid && nextValid
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
