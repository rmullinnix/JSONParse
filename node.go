package JSONParse

import (
)

// NodeState is used to mark each node as not validated, in progress, valid or invald
// this is used to prevent validating the node more than once as json documents can
// refer to the same section multiple times
// It is also used to lock an object if multiple validators work need to work in concert
// to validate the node (e.g., properties, patternProperties, additionalProperties
type NodeState int

// Valid states for NodeState
const (
	VIRGIN=iota
	VALIDATE_IN_PROGRESS
	VALID
	INVALID
	NODE_MUTEX
	NODE_SEMAPHORE
)

// Node for each item in the json tree
// members are store in a map and array members are stored in an arry
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

// intializes a new json tree for a json document
// a parsed json document is passed in as part of initialization
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

// allocates a new node off of the current node
// the input parameter is the token index generated from the parser
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

// creates a new node of type "object"
func (jn *JSONNode) NewObject(v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = "object"

	jn.value = node

	return node
}

// creates a new node of type "member"
func (jn *JSONNode) NewMember(name string, v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = "member"

	if name == "$ref" {
		node.memNodeType = "reference"
	}

	jn.namedKids[name] = node

	return node
}

// creates a new node of type "array"
func (jn *JSONNode) NewArray(v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = "array"

	jn.value = node

	return node
}

// creates a new "value" node to be stored in an array
func (jn *JSONNode) NewArrayValue(v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = "value"

	jn.unnamedKids = append(jn.unnamedKids, node)

	return node
}

// creates a new node of type "reference" to support json reference
func (jn *JSONNode) NewReference(name string, v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = "reference"
	node.name = name

	jn.value = node

	return node
}

// returns the parent node of the current node
// if the current node is a referenced item, it refers back to the
// parent where the referred object resides
func (jn *JSONNode) GetParent() *JSONNode {
	return jn.parent
}

// returns the root node of the json object
func (jn *JSONNode) GetRoot() *JSONNode {
	return jn.root
}

// sets the value of the node
func (jn *JSONNode) SetValue(val interface{}) {
	jn.value = val
}

// gets the value of the node, knowing the type (can use reflection)
// allows the value to be cast as the appropriate value
func (jn *JSONNode) GetValue() interface{} {
	return jn.value
}

// sets the type of node
func (jn *JSONNode) SetType(nodeType string) {
	jn.nodeType = nodeType
}

// gets the type of node
func (jn *JSONNode) GetType() string {
	return jn.nodeType
}

// retrieve the number of members in the node
// members are named values
func (jn *JSONNode) GetMemberCount() int {
	return len(jn.namedKids)
}

// retrieve the number of members in the node
// members are named values
func (jn *JSONNode) GetCount() int {
	return len(jn.unnamedKids)
}

// set the value type of the members contained in the node
func (jn *JSONNode) SetMemberType(nodeType string) {
	if nodeType == "member" {
		panic("invalid node member type")
	}
	jn.memNodeType = nodeType
}

// get the value type of the members contained in the node
func (jn *JSONNode) GetMemberType() string {
	return jn.memNodeType
}

// get the current state of the node - used during validation
func (jn *JSONNode) GetState() NodeState {
	return jn.state
}

// get the current state of the node - used during validation
func (jn *JSONNode) SetState(state NodeState) {
	jn.state = state
}

// find a named member in the current node
// returns a pointer to the node if found, and nil if not found
// If json references exist in the list of members, they are resolved
// before searching. Once resolved, they become permanent links
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

// resets the node so the Member list or array can be iterated
// from the beginning
func (jn *JSONNode) ResetIterate() {
	jn.curIndex = 0
	jn.nameArray = make([]string, len(jn.namedKids))
	for key, _ := range jn.namedKids {
		jn.nameArray[jn.curIndex] = key
		jn.curIndex++
	}
	jn.curIndex = 0
}

// retrieve the next named member
// returns the member key (name) and the node
// a valid of nil for the node indicates the end of the list
// json references are resolved and become permanent members of the list
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

// retrieves the next array node from the list
func (jn *JSONNode) GetNext() (*JSONNode) {
	if jn.curIndex >= len(jn.unnamedKids) {
		return nil
	}

	item := jn.unnamedKids[jn.curIndex]
	jn.curIndex++

	return item
}

// retrieves the next array value from the list
func (jn *JSONNode) GetNextValue() interface{} {
	if jn.curIndex >= len(jn.unnamedKids) {
		return nil
	}

	item := jn.unnamedKids[jn.curIndex]
	jn.curIndex++

	return item.GetValue()
}

// removes the $ref tag from the list of members and
// replaces with the members of the referred to json section
// will link to internal as well as external document sections
func (jn *JSONNode) CollapseReference(parent *JSONNode) {
	if jn.memNodeType != "reference" {
		Error.Panicln(" invalid ref node type", jn.memNodeType)
//		return
	}

	// traverse nodes until end of reference tags
	// refNode is an object that has members (namedKids) 
	// move reference members to this node
	// results in reference being removed for future traversals
	refNode := jn.followReference(jn.root.doc.references)

	if refNode == nil {
		jn.dump()
		panic("invalid reference")
	}
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

// returns the referred node of hte reference.  If the reference has not
// been resolved yet, it access the references table to set the value
func (jn *JSONNode) followReference(references map[string]*JSONNode) *JSONNode {
	ptr := jn.GetValue().(*JSONNode)

	if ptr.GetValue() == nil {
		Trace.Println("reference ", ptr.name)
		if ptrValue, found := references[ptr.name]; found {
			ptr.SetValue(ptrValue)
		} else {
			return nil
		}
	}

	return ptr.GetValue().(*JSONNode)
}

// internal troubleshooting
func (jn *JSONNode) dump() {
	Trace.Println("NodeType: ", jn.nodeType)
	Trace.Println(" memNodeType: ", jn.memNodeType)
	Trace.Println(" name: ", jn.name)
	Trace.Println(" Named kids")

	for key, item := range jn.namedKids {
		Trace.Println("\t", key, ": ", item)
	}

	Trace.Println(" Unnamed kids")

	for i := 0; i < len(jn.unnamedKids); i++  {
		Trace.Println("\t", jn.unnamedKids[i])
	}
}
