package JSONParse

import (
	"fmt"
	"regexp"
//	"strconv"
	"strings"
)

const (
	UNKNOWN		= 0
	STRING          = 1
	J_TRUE          = 2
	J_FALSE		= 4
	NUMBER          = 8
	BEGIN_ARRAY     = 16
	END_ARRAY       = 32
	BEGIN_OBJECT    = 64
	END_OBJECT      = 128
	VALUE_SEPARATOR = 256
	NAME_SEPARATOR  = 512
	J_NULL          = 1024
	REF		= 2048
	END_OF_SOURCE   = 4096
)

const (
	JP_WARNING = iota
	JP_ERROR
	JP_FATAL
)

type lineItem struct {
	tokenStart int
	tokenEnd   int
	indent 	   int
}

type ParseError struct {
	Error      string
	ErrorLevel int
	LineNumber int
	Line	   string
	Offset     int
}

type JSONParser struct {
	lines         []lineItem
	lineCount     int
	curIndex      int
	variables     map[int]string
	curTokenType  int
	curTokenVar   string
	curIndent     int
	errorList     []ParseError
	errorCount    int
	maxError      int
	source        string
	raw           string
	tokens	      []int
	references    map[string]*JSONNode
	jsonDoc	      *JSONNode
	extDocs	      map[string]*JSONNode
}

var validNum		*regexp.Regexp
var jsonTree		*JSONNode

func NewJSONParser(source string, maxError int) *JSONParser {
	jp := new(JSONParser)
	jp.source = source
	jp.errorList = make([]ParseError, 0)
	jp.errorCount = 0
	jp.maxError = maxError
	jp.variables = make(map[int]string)
	jp.curTokenType = -1
	jp.curIndex = 0
	jp.lines = make([]lineItem, 1)
	var newLine		lineItem
	newLine.indent = 0
	jp.lines[0] = newLine
	jp.jsonDoc = NewJSONTree(jp)
	jp.tokens = make([]int, 0)
	jp.references = make(map[string]*JSONNode)
	jp.extDocs = make(map[string]*JSONNode)

	validNum = regexp.MustCompile(`-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?`)

	return jp
}

func (jp *JSONParser) Parse() (bool, []ParseError) {
	// read from source
	raw, ferr := loadDoc(jp.source)
	if ferr != nil {
		jp.addError("Unable to read json", JP_FATAL)
		return false, jp.errorList
	}

	jp.raw = string(raw)

	jp.tokenize()

	// retrieve tokens until entire string is parsed or max errors is
	// reached
	err := jp.expectToken(BEGIN_OBJECT)
	if err != nil {
		return false, jp.errorList
	}

	valid := jp.parseObject(jp.jsonDoc)
	if !valid {
		fmt.Println("invalid object")
		return false, jp.errorList
	}

	jp.newLine(-1)
	// check for end of source
	err = jp.expectToken(END_OF_SOURCE)
	if err != nil {
		return false, jp.errorList
	}

	jp.resolveReferences()

	return true, jp.errorList
}

func (jp *JSONParser) GetDoc() *JSONNode {
	return jp.jsonDoc
}

func (jp *JSONParser) Pretty() string {
	return jp.prettyTokens(jp.tokens)
}

func (jp *JSONParser) prettyTokens(tokens []int) string {
	indent := 0
	output := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == BEGIN_OBJECT {
			output += "{\n"
			indent++
			output += strings.Repeat("\t", indent)
		} else if tokens[i] == END_OBJECT {
			indent--
			output += "\n"
			output += strings.Repeat("\t", indent)
			output += `}`
		} else if tokens[i] == BEGIN_ARRAY {
			output += "[\n"
			indent++
			output += strings.Repeat("\t", indent)
		} else if tokens[i] == END_ARRAY {
			indent--
			output += "\n"
			output += strings.Repeat("\t", indent)
			output += `]`
		} else if tokens[i] == VALUE_SEPARATOR {
			output += ",\n"
			output += strings.Repeat("\t", indent)
		} else if tokens[i] == J_FALSE {
			output += `false`
		} else if tokens[i] == J_TRUE {
			output += `true`
		} else if tokens[i] == J_NULL {
			output += `null`
		} else if tokens[i] == STRING {
			i++
			output += `"` + jp.variables[tokens[i]] + `"`
		} else if tokens[i] == REF {
			output += `"$ref"`
		} else if tokens[i] == NUMBER {
			i++
			output += jp.variables[tokens[i]]
		} else if tokens[i] == NAME_SEPARATOR {
			output += `: `
		}
	}

	return output
}

func (jp *JSONParser) tokenize() {
	indxToken := 0
	indxVar := 0

	for {
		token := jp.getToken()
		jp.tokens = append(jp.tokens, token)
		indxToken++
		if token == STRING || token == NUMBER {
			jp.variables[indxVar] = jp.curTokenVar
			jp.tokens = append(jp.tokens, indxVar)
			indxVar++
		} else if token == END_OF_SOURCE {
			break
		}
	}

	jp.curIndex = -1
}

func (jp *JSONParser) parseArray(arr *JSONNode) bool {
	if jp.curTokenType != BEGIN_ARRAY {
		return false
	}

	for {
		val := arr.NewArrayValue(jp.curIndex)
		if !jp.parseValue(val) {
			break
		}

		err := jp.expectToken(VALUE_SEPARATOR | END_ARRAY)

		if err != nil {
			break
		}

		if jp.curTokenType == END_ARRAY {
			break
		}
	}

	if jp.curTokenType != END_ARRAY {
		return false
	}

	return true
}

func (jp *JSONParser) parseMember(mem *JSONNode) bool {
	if !(jp.curTokenType == STRING || jp.curTokenType == REF) {
		return false
	}

	curTokenType := jp.curTokenType

	err := jp.expectToken(NAME_SEPARATOR)
	if err != nil {
		return false
	}

	if curTokenType == REF {
		return jp.parseRefPtr(mem)
	}

	return jp.parseValue(mem)
}

// value = false / null / true / object / array / number / string
func (jp *JSONParser) parseValue(val *JSONNode) bool {
	err := jp.expectToken(J_TRUE | J_FALSE | J_NULL | BEGIN_OBJECT | BEGIN_ARRAY | NUMBER | STRING)
	if err != nil {
		return false
	}

	if jp.curTokenType == BEGIN_OBJECT {
		val.SetMemberType("object")
		obj := val.NewObject(jp.curIndex)

		return jp.parseObject(obj)
	} else if jp.curTokenType == BEGIN_ARRAY {
		val.SetMemberType("array")
		arr := val.NewArray(jp.curIndex)

		return jp.parseArray(arr)
	} else {
		val.SetType(tokenToString(jp.curTokenType))
		if jp.curTokenType == J_TRUE {
			val.SetMemberType("boolean")
			val.SetValue(true)
		} else if jp.curTokenType == J_FALSE {
			val.SetMemberType("boolean")
			val.SetValue(false)
		} else {
			val.SetMemberType(strings.ToLower(tokenToString(jp.curTokenType)))
			val.SetValue(jp.curTokenVar)
		}
	}

	return true
}

func (jp *JSONParser) parseRefPtr(ref *JSONNode) bool {
	err := jp.expectToken(STRING | BEGIN_OBJECT)
	if err != nil {
		return false
	}

	// definition of $ref, not actual $ref
	if jp.curTokenType == BEGIN_OBJECT {
		ref.SetType("string")
		ref.SetMemberType("object")
		obj := ref.NewObject(jp.curIndex)

		return jp.parseObject(obj)
	}

	ref.NewReference(jp.curTokenVar, jp.curIndex)

	jp.references[jp.curTokenVar] = nil

	return true
}
		
//  object = begin-object [ member *( value-separator member ) ]
//      end-object
func (jp *JSONParser) parseObject(obj *JSONNode) bool {
	if jp.curTokenType != BEGIN_OBJECT {
		return false
	}

	for {
		err := jp.expectToken(STRING | END_OBJECT | REF)
		if err != nil {
			return false
		}

		if jp.curTokenType == END_OBJECT {
			break
		}

		mem := obj.NewMember(jp.curTokenVar, jp.curIndex)

		if !jp.parseMember(mem) {
			return false
		}

		err = jp.expectToken(VALUE_SEPARATOR | END_OBJECT)
		if err != nil {
			return false
		}

		if jp.curTokenType == END_OBJECT {
			break
		}
	}

	return true
}

// Retrieve the next token in the source json
func (jp *JSONParser) getToken() int {
	tokType := jp.getWord()
	jp.curTokenType = tokType

	return tokType
}

func (jp *JSONParser) expectToken(valTokens int) *ParseError {
	token := jp.nextToken()

	if token&valTokens == token {
		return nil
	}

	return jp.addError("Unexpected token: Expecting " + tokenToString(valTokens) + " - Encountered " + tokenToString(token), JP_ERROR)
}

func (jp *JSONParser) nextToken() int {
	jp.curIndex++
	jp.curTokenType = jp.tokens[jp.curIndex]

	if jp.curTokenType == STRING || jp.curTokenType == NUMBER {
		jp.curIndex++
		jp.curTokenVar = jp.variables[jp.tokens[jp.curIndex]]
	} else if jp.curTokenType == REF {
		jp.curTokenVar = "$ref"
	}

	return jp.curTokenType
}

// Retrieve the next word in the source json
func (jp *JSONParser) getWord() int {
	var letter	string

	// skip white space
	for {
		if jp.curIndex >= len(jp.raw)  {
			return END_OF_SOURCE
		}

		letter = jp.raw[jp.curIndex:jp.curIndex+1]
		jp.curIndex++

		if strings.Contains(" \n\t", letter)  {
			continue
		}

		break
	}

	jp.curTokenVar = letter
	if letter == "\"" {
		if endQuote := strings.Index(jp.raw[jp.curIndex:], "\""); endQuote > 0 {
			jp.curTokenVar = jp.raw[jp.curIndex : jp.curIndex+endQuote]
			jp.curIndex += endQuote + 1
			if jp.curTokenVar == "$ref" {
				return REF
			}
		} else if endQuote == 0 {
			jp.curTokenVar = ""
			jp.curIndex++
		} else {
			jp.curTokenVar = "error"
			jp.addError("No matching end quote", JP_FATAL)
			return STRING
		}
		return STRING
	}

	if letter == "{" {
		return BEGIN_OBJECT
	}

	if letter == "}" {
		return END_OBJECT
	}

	if letter == "[" {
		return BEGIN_ARRAY
	}

	if letter == "]" {
		return END_ARRAY
	}

	if letter == ":" {
		return NAME_SEPARATOR
	}

	if letter == "," {
		return VALUE_SEPARATOR
	}

	if strings.Contains("-0123456789", letter) {
		// grab string to next space, check for int or float
		endWord := strings.IndexAny(jp.raw[jp.curIndex:], " ,:}]\n")
		jp.curTokenVar = jp.raw[jp.curIndex-1: jp.curIndex + endWord]

		jp.curIndex += endWord
		
		if validNum.MatchString(jp.curTokenVar) {
			return NUMBER
		} else {
			jp.addError("Invalid number format", JP_ERROR)
			return NUMBER
		}
	}

	if strings.Contains("tfn", letter) {
		// potential boolean or null - need to validate
		endWord := strings.IndexAny(jp.raw[jp.curIndex:], " ,:}]\n")
		if endWord > 0 {
			jp.curTokenVar = jp.raw[jp.curIndex - 1 : jp.curIndex + endWord]
			jp.curIndex += endWord

			if jp.curTokenVar == "true" {
				return J_TRUE
			} else if jp.curTokenVar == "false" {
				return J_FALSE
			} else if jp.curTokenVar == "null" {
				return J_NULL
			} else {
				jp.addError("Unquoted string", JP_ERROR)
				return UNKNOWN
			}
		}
	}

	jp.addError("Encountered invalid token", JP_FATAL)

	return 0
}

func (jp *JSONParser) newLine(indent int) {
	var newLine		lineItem

	jp.lines[jp.lineCount].tokenEnd = jp.curIndex - 1
	
	newLine.tokenStart = jp.curIndex
	jp.curIndent += indent
	newLine.indent = jp.curIndent

	jp.lineCount++
	jp.lines = append(jp.lines, newLine)
}

func (jp *JSONParser) addError(errText string, level int) *ParseError {
	var pError	ParseError

	pError.Error = errText
	pError.ErrorLevel = level
	pError.Offset = jp.curIndex

	jp.errorList = append(jp.errorList, pError)

	return &pError
}

func tokenToString(token int) string {
	output := ""

	if token & STRING == STRING {
		output = "string "
	}
	if token & J_TRUE == J_TRUE {
		output += "true "
	}
	if token & J_FALSE == J_FALSE {
		output += "false "
	}
	if token & NUMBER == NUMBER {
		output += "NUMBER "
	}
	if token & BEGIN_ARRAY == BEGIN_ARRAY {
		output += "BEGIN_ARRAY "
	}
	if token & END_ARRAY == END_ARRAY {
		output += "END_ARRAY "
	}
	if token & BEGIN_OBJECT == BEGIN_OBJECT {
		output += "BEGIN_OBJECT "
	}
	if token & END_OBJECT == END_OBJECT {
		output += "END_OBJECT "
	}
	if token & VALUE_SEPARATOR == VALUE_SEPARATOR {
		output += "VALUE_SEPARATOR "
	}
	if token & NAME_SEPARATOR == NAME_SEPARATOR {
		output += "NAME_SEPARATOR "
	}
	if token & J_NULL == J_NULL {
		output += "NULL "
	}
	if token & END_OF_SOURCE == END_OF_SOURCE {
		output += "END_OF_SOURCE "
	}
	if output == "" {
		output = "UNKNOWN"
	}

	output = strings.TrimSuffix(output, " ")
	output = strings.Replace(output, " ", " | ", -1)

	return output
}
