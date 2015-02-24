package JSONParse

import (
	"strings"
)

type NodeType int

// valid NodeType
const (
	N_OBJECT=iota
	N_MEMBER
	N_ARRAY
	N_REFERENCE
)

// valid ValueType
const (
	V_OBJECT=iota
	V_ARRAY
	V_REFERENCE
	V_STRING
	V_NUMBER
	V_BOOLEAN
	V_EMPTY
	V_NULL
)

type ValueType		int

// Node for each item in the json tree
// members are store in a map and array members are stored in an arry
type JSONNode struct {
	parent		*JSONNode
	doc		*JSONParser
	value		interface{}
	nodeType	NodeType
	valType		ValueType
	name		string
	root		*JSONNode
	namedKids	map[string]*JSONNode
	nameArray	[]string
	unnamedKids	[]*JSONNode
	curIndex	int
	tokenIndex	int
	lineNumber	int
}

// intializes a new json tree for a json document
// a parsed json document is passed in as part of initialization
func NewJSONTree(doc *JSONParser) *JSONNode {
	jn := new(JSONNode)

	jn.doc = doc
	jn.root = jn
	jn.parent  = nil
	jn.nodeType = N_OBJECT
	jn.name = "root"

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

	node.unnamedKids = make([]*JSONNode, 0)
	node.namedKids = make(map[string]*JSONNode)

	return node
}

// creates a new node of type "object"
func (jn *JSONNode) NewObject(v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = N_OBJECT

	jn.unnamedKids = append(jn.unnamedKids, node)
	jn.valType = V_OBJECT

	return node
}

// creates a new node of type "member"
func (jn *JSONNode) NewMember(name string, v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = N_MEMBER
	node.name = name

	if name == "$ref" {
		node.nodeType = N_REFERENCE
	}

	jn.namedKids[name] = node

	return node
}

// creates a new "value" node to be stored in an array
func (jn *JSONNode) NewArrayValue(v_indx int) *JSONNode {
	node := jn.newNode(v_indx)
	node.nodeType = N_ARRAY

	jn.unnamedKids = append(jn.unnamedKids, node)

	return node
}

// creates a new node of type "reference" to support json reference
//func (jn *JSONNode) NewReference(name string, v_indx int) *JSONNode {
//	node := jn.newNode(v_indx)
//	node.nodeType = V_REFERENCE
//	node.name = name

//	jn.namedKids[name] = node

//	return node
//}

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
func (jn *JSONNode) SetType(nodeType NodeType) {
	jn.nodeType = nodeType
}

// gets the type of node
func (jn *JSONNode) GetType() NodeType {
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
func (jn *JSONNode) SetValueType(memType ValueType) {
	jn.valType = memType
}

// get the value type of the members contained in the node
func (jn *JSONNode) GetValueType() ValueType {
	return jn.valType
}

// find a named member in the current node
// returns a pointer to the node if found, and nil if not found
// If json references exist in the list of members, they are resolved
// before searching. Once resolved, they become permanent links
func (jn *JSONNode) Find(name string) (*JSONNode, bool) {
	if len(jn.namedKids) == 1 {
		if refNode, found := jn.namedKids["$ref"]; found {
			refNode.collapseReference(jn)
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
func (jn *JSONNode) GetNextMember(chaseRef bool) (string, *JSONNode) {
	if jn.curIndex >= len(jn.nameArray) {
		return "", nil
	}

	// if type is ref, replace reference node with actual members
	// should only be a single member -- len(jn.nameArray) == 1
	for {
		var first	*JSONNode
		var hasRef	bool

		if first, hasRef = jn.namedKids["$ref"]; hasRef {
			if first.nodeType == N_OBJECT  {
				hasRef = false
			} else if first.nodeType == N_REFERENCE {
				if chaseRef {
					first.collapseReference(jn)
				} else {
					hasRef = false
				}
				
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

// return the current node as json (includes children)
func (jn *JSONNode) GetJson() string {
	output := ""
	if jn.nodeType == N_OBJECT {
		output = "{"
		output += jn.getKidsJson()
		output += "}"
	} else if jn.nodeType == N_ARRAY {
		output = "["
		output += jn.getKidsJson()
		output += "]"
	} else if jn.nodeType == N_MEMBER {
		output += `"` + jn.name + `":`

		output += jn.GetValueJson()
	} else if jn.nodeType == N_REFERENCE {
		output += `"$ref":`
		output += jn.GetValueJson()
	} else {
		Trace.Println("unknown type")
		jn.dump()
	}

	return output
}

func (jn *JSONNode) GetValueJson() string {
	output := ""
	if jn.valType == V_OBJECT { 
		if jn.nodeType != N_OBJECT {
			output += `{` + jn.getKidsJson() + `}`
		} else {
			output += jn.getKidsJson()
		}
	} else if jn.valType == V_ARRAY {
		output += "[" + jn.getValueKidsJson() +"]"
	} else if jn.valType == V_NUMBER {
		output += " " + jn.value.(string)
	} else if jn.valType == V_NULL {
		output += " null"
	} else if jn.valType == V_BOOLEAN {
		val := jn.value.(bool)
		if val {
			output += " true"
		} else {
			output += " false"
		}
	} else if jn.valType == V_STRING {
		output = output + `"` + jn.value.(string) + `"`

	} else {
		Trace.Println("unknown type")
		jn.dump()
	}

	return output
}

func (jn *JSONNode) getKidsJson() string {
	output := ""
	jn.ResetIterate()
	for {
		_, item := jn.GetNextMember(false)
		if item == nil {
			break
		}

		output += item.GetJson()
		output += ","
	}

	jn.ResetIterate()
	for {
		item := jn.GetNext()
		if item == nil {
			break
		}

		output += item.GetValueJson()
		output += ","
	}
	output = strings.TrimSuffix(output, ",")
	return output
}

func (jn *JSONNode) getValueKidsJson() string {
	output := ""
	jn.ResetIterate()

	for {
		item := jn.GetNext()
		if item == nil {
			break
		}

		output += item.GetValueJson()
		output += ","
	}
	output = strings.TrimSuffix(output, ",")
	return output
}

// removes the $ref tag from the list of members and
// replaces with the members of the referred to json section
// will link to internal as well as external document sections
func (jn *JSONNode) collapseReference(parent *JSONNode) {
	if jn.valType != V_STRING {
		return
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
		if key, item := refNode.GetNextMember(true); item == nil {
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
//	ptr := jn.GetValue().(*JSONNode)

	ref := jn.GetValue().(string)
	if len(ref) > 0 {
		if ptrValue, found := references[jn.GetValue().(string)]; found {
			jn.SetValue(ptrValue)
			jn.SetValueType(V_OBJECT)
		} else {
			return nil
		}
	}

	return jn.GetValue().(*JSONNode)
}

// internal troubleshooting
func (jn *JSONNode) dump() {
	Trace.Println("NodeType: ", jn.nodeType)
	Trace.Println(" valType: ", jn.valType)
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
