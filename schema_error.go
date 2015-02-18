package JSONParse

import (
)

type SchemaError struct {
	node		*JSONNode
	message		string
	level		int
}

type SchemaErrors struct {
	errorList	[]*SchemaError
}

func NewSchemaErrors() (*SchemaErrors) {
	se := new(SchemaErrors)
	se.errorList = make([]*SchemaError, 0)

	return se
}

func (se *SchemaErrors) Add(node *JSONNode, message string, level int) {
	err := new(SchemaError)
	err.node = node
	err.level = level
	err.message = message

	se.errorList = append(se.errorList, err)
}

func (se *SchemaErrors) Output() {
	for i := 0; i < len(se.errorList); i++ {
		OutputError(se.errorList[i].node, se.errorList[i].message)
	}
}
