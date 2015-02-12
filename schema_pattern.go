package JSONParse

import (
//	"fmt"
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
