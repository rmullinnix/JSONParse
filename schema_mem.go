package JSONParse

import (
	"strconv"
)

var depth int

func validMember(stack_id string, name string, mem *JSONNode, schemaMem *JSONNode) bool {
	depth++
	Trace.Println(stack_id, "validMember -- depth", depth, "-- name <", name, ">")
		
	if mem.GetType() == N_MEMBER {
		name = mem.name
	}

	valid := true

	schemaMem.ResetIterate()
	index := 1
	for {
	
		key, item := schemaMem.GetNextMember(true)
		if item == nil {
			break
		}
		Trace.Println(stack_id, "validMember == key: <" + key + ">")

		if validator, found := keywords[key]; found {
			Trace.Println(stack_id, "validMember", name, "against schema mem", key)
			new_stack_id := stack_id + "." + strconv.Itoa(index)
			index++
			if suppress {
				nextValid := validator(new_stack_id, mem, item, schemaMem, suppressErrors)
				valid = valid && nextValid
			} else {
				nextValid := validator(new_stack_id, mem, item, schemaMem, schemaErrors)
				valid = valid && nextValid
			}
		}
	}

	depth--

	Trace.Println(stack_id, "validMember", valid)
	return valid
}
