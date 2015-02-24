package JSONParse

import (
	"strconv"
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
		node := se.errorList[i].node
		level := se.errorList[i].level
		errMsg := se.errorList[i].message

		tokenIndex := node.tokenIndex
		parser := node.root.doc
		tokenStart := 0
		tokenEnd := len(parser.tokens)

		if tokenIndex > 15 {
			tokenStart = tokenIndex - 15
		}

		if tokenIndex < tokenEnd - 15 {
			tokenEnd = tokenIndex + 15
		}

		if level == JP_ERROR || level == JP_FATAL {
			output :=  parser.prettyTokens(tokenStart, tokenEnd)
			Error.Println(strconv.Itoa(node.lineNumber) + ":", errMsg + "\n" + output)
		} else if level == JP_WARNING {
			Error.Println(strconv.Itoa(node.lineNumber) + ":", errMsg)
		} else if level == JP_INFO {
			Error.Println(strconv.Itoa(node.lineNumber) + ":", errMsg)
		}
	}
}
