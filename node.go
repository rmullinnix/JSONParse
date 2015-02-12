package JSONParse

import (
	"fmt"
//	"strings"
)

type NodeState int

const (
	VIRGIN=iota
	VALIDATE_IN_PROGRESS
	VALID
	INVALID
	NODE_MUTEX
)

type JSONNode struct {
	parent		*JSONNode
	doc		*JSONParser
	value		interface{}
	state		NodeState
	nodeType	string
	memNodeType	string
	name		string
	root		*JSONNode
	namedKids	map[string]*JSONNode
	nameArray	[]string
	unnamedKids	[]*JSONNode
	curIndex	int
	tokenIndex	int
}

func NewJSONTree(doc *JSONParser) *JSONNode {
	jn := new(JSONNode)

	jn.doc = doc
	jn.root = jn
	jn.parent  = nil
	jn.nodeType = "object"
	jn.state = VIRGIN

	jn.unnamedKids = make([]*JSONNode, 0)
	jn.namedKids = make(map[string]*JSONNode)

	return jn
}

func (jn *JSONNode) newNode(v_indx int) *JSONNode {
	node := new(JSONNode)

	node.parent = jn
	node.root = jn.root
	node.tokenIndex = v_indx
	node.state = VIRGIN

	node.unnamedKids = make([]*JSONNode, 0)
	node.namedKids = make(map[string]*JSONNode)

	return node
}

func (jn *JSONNode) NewObject(v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = "object"

	jn.value = node

	return node
}

func (jn *JSONNode) NewMember(name string, v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = "member"

	if name == "$ref" {
		node.memNodeType = "reference"
	}

	jn.namedKids[name] = node

	return node
}

func (jn *JSONNode) NewArray(v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = "array"

	jn.value = node

	return node
}

func (jn *JSONNode) NewArrayValue(v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = "value"

	jn.unnamedKids = append(jn.unnamedKids, node)

	return node
}

func (jn *JSONNode) NewReference(name string, v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = "reference"
	node.name = name

	jn.value = node

	return node
}

func (jn *JSONNode) GetParent() *JSONNode {
	return jn.parent
}

func (jn *JSONNode) GetRoot() *JSONNode {
	return jn.root
}

func (jn *JSONNode) SetValue(val interface{}) {
	jn.value = val
}

func (jn *JSONNode) GetValue() interface{} {
	return jn.value
}

func (jn *JSONNode) SetType(nodeType string) {
	jn.nodeType = nodeType
}

func (jn *JSONNode) GetType() string {
	return jn.nodeType
}

func (jn *JSONNode) GetMemberCount() int {
	return len(jn.namedKids)
}

func (jn *JSONNode) SetMemberType(nodeType string) {
	if nodeType == "member" {
		panic("invalid node member type")
	}
	jn.memNodeType = nodeType
}

func (jn *JSONNode) GetMemberType() string {
	return jn.memNodeType
}

func (jn *JSONNode) GetState() NodeState {
	return jn.state
}

func (jn *JSONNode) SetState(state NodeState) {
	jn.state = state
}

func (jn *JSONNode) Find(name string) (*JSONNode, bool) {
	if len(jn.namedKids) == 1 {
		if refNode, found := jn.namedKids["$ref"]; found {
			refNode.CollapseReference(jn)
		}
	}

	if item, found := jn.namedKids[name]; found {
		return item, true
	} else {
		return nil, false
	}
}

func (jn *JSONNode) ResetIterate() {
	jn.curIndex = 0
	jn.nameArray = make([]string, len(jn.namedKids))
	for key, _ := range jn.namedKids {
		jn.nameArray[jn.curIndex] = key
		jn.curIndex++
	}
	jn.curIndex = 0
}

func (jn *JSONNode) GetNextMember() (string, *JSONNode) {
	if jn.curIndex >= len(jn.nameArray) {
		return "", nil
	}

	// if type is ref, replace reference node with actual members
	// should only be a single member -- len(jn.nameArray) == 1
	for {
		var first	*JSONNode
		var hasRef	bool

		if first, hasRef = jn.namedKids["$ref"]; hasRef {
			if first.memNodeType == "object"  {
				hasRef = false
			} else if first.memNodeType == "reference" {
				first.CollapseReference(jn)
			} else {
				hasRef = false
			}
		}

		if !hasRef {
			break
		}
	}

	if jn.curIndex >= len(jn.nameArray) {
		return "", nil
	}

	item := jn.namedKids[jn.nameArray[jn.curIndex]]

	key := jn.nameArray[jn.curIndex]
	jn.curIndex ++

	return key, item
}

func (jn *JSONNode) GetNext() (*JSONNode) {
	if jn.curIndex >= len(jn.unnamedKids) {
		return nil
	}

	item := jn.unnamedKids[jn.curIndex]
	jn.curIndex++

	return item
}

func (jn *JSONNode) GetNextValue() interface{} {
	if jn.curIndex >= len(jn.unnamedKids) {
		return nil
	}

	item := jn.unnamedKids[jn.curIndex]
	jn.curIndex++

	return item.GetValue()
}

func (jn *JSONNode) CollapseReference(parent *JSONNode) {
	if jn.memNodeType != "reference" {
		fmt.Println(" invalid ref node type", jn.memNodeType)
//		return
	}

	// traverse nodes until end of reference tags
	// refNode is an object that has members (namedKids) 
	// move reference members to this node
	// results in reference being removed for future traversals
	refNode := jn.FollowReference(jn.root.doc.references)

	delete(parent.namedKids, "$ref")

	// calling GetNextMember() will chase references until valid object
	refNode.ResetIterate()
	for {
		if key, item := refNode.GetNextMember(); item == nil {
			break
		} else {
			parent.namedKids[key] = item
		}
	}
	parent.ResetIterate()
}

func (jn *JSONNode) FollowReference(references map[string]*JSONNode) *JSONNode {
	ptr := jn.GetValue().(*JSONNode)

	if ptr.GetValue() == nil {
		if ptrValue, found := references[ptr.name]; found {
			ptr.SetValue(ptrValue)
		} else {
			return nil
		}
	}

	return ptr.GetValue().(*JSONNode)
}

func (jn *JSONNode) dump() {
	fmt.Println("NodeType: ", jn.nodeType)
	fmt.Println(" memNodeType: ", jn.memNodeType)
	fmt.Println(" name: ", jn.name)
	fmt.Println(" Named kids")

	for key, item := range jn.namedKids {
		fmt.Println("\t", key, ": ", item)
	}

	fmt.Println(" Unnamed kids")

	for i := 0; i < len(jn.unnamedKids); i++  {
		fmt.Println("\t", jn.unnamedKids[i])
	}
}
