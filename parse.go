package JSONParse

import (
	"regexp"
//	"strconv"
	"strings"
)

// A JSON text is a sequence of tokens.  The set of tokens includes six
//   structural characters, strings, numbers, and three literal names.
//   A JSON text is a serialized object or array.
//      JSON-text = object / array
//   These are the six structural characters:
//
//      begin-array     = ws %x5B ws  ; [ left square bracket
//      begin-object    = ws %x7B ws  ; { left curly bracket
//      end-array       = ws %x5D ws  ; ] right square bracket
//      end-object      = ws %x7D ws  ; } right curly bracket
//      name-separator  = ws %x3A ws  ; : colon
//      value-separator = ws %x2C ws  ; , comma
//
//  A JSON value MUST be an object, array, number, or string, or one of
//   the following three literal names:
//
//      false null true

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
	JP_INFO = iota
	JP_WARNING
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
	Offset     int
}

type JSONParser struct {
	lines         []lineItem
	lineCount     int
	curIndex      int
	ltrIndex      int
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
	resolveRef    bool
	references    map[string]*JSONNode
	jsonDoc	      *JSONNode
	extDocs	      map[string]*JSONNode
}

var validNum		*regexp.Regexp
var jsonTree		*JSONNode

// Creates a new JSON Parser
func NewJSONParser(source string, maxError int, level string) *JSONParser {
	jp := new(JSONParser)
	jp.source = source
	jp.errorList = make([]ParseError, 0)
	jp.errorCount = 0
	jp.maxError = maxError
	jp.variables = make(map[int]string)
	jp.curTokenType = -1
	jp.curIndex = 0
	jp.ltrIndex = 0
	jp.lines = make([]lineItem, 1)
	var newLine		lineItem
	newLine.indent = 0
	jp.lines[0] = newLine
	jp.jsonDoc = NewJSONTree(jp)
	jp.tokens = make([]int, 0)
	jp.resolveRef = true
	jp.references = make(map[string]*JSONNode)
	jp.extDocs = make(map[string]*JSONNode)

	if level != "default" {
		outputInit(level)
	}

	validNum = regexp.MustCompile(`-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?`)

	return jp
}

// Parses the json by tokenizing the stream and parsing the object
//  References addresses are stored in a map structure and
//  are solved when traversing the json tree
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
		Error.Println("invalid object")
		return false, jp.errorList
	}

	// check for end of source
	err = jp.expectToken(END_OF_SOURCE)
	if err != nil {
		return false, jp.errorList
	}

	if jp.resolveRef {
		jp.resolveReferences()
	}

	return true, jp.errorList
}

// returns a pointer to the json document
func (jp *JSONParser) GetDoc() *JSONNode {
	return jp.jsonDoc
}

func (jp *JSONParser) ResolveRefs(val bool) {
	jp.resolveRef = val
}

// converts the json stream into constants representing items
//   names are giving an index and stored in a reference map
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

// 2.3.  Arrays
// 
//    An array structure is represented as square brackets surrounding zero
//    or more values (or elements).  Elements are separated by commas.
// 
//       array = begin-array [ value *( value-separator value ) ] end-array
func (jp *JSONParser) parseArray(arr *JSONNode) bool {
	if jp.curTokenType != BEGIN_ARRAY {
		return false
	}

	for {
		if jp.tokens[jp.curIndex + 1] != END_ARRAY {
			val := arr.NewArrayValue(jp.curIndex)
			val.lineNumber = len(jp.lines)
			if !jp.parseValue(val) {
				break
			}
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

//       member = string name-separator value
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
		mem.SetValueType(V_REFERENCE)
		return jp.parseRefPtr(mem)
	}

	return jp.parseValue(mem)
}

// value = false / null / true / object / array / number / string
func (jp *JSONParser) parseValue(val *JSONNode) bool {
	err := jp.expectToken(J_TRUE | J_FALSE | J_NULL | BEGIN_OBJECT | BEGIN_ARRAY | NUMBER | STRING | REF)
	if err != nil {
		return false
	}

	if jp.curTokenType == BEGIN_OBJECT {
		val.SetValueType(V_OBJECT)
		obj := val.NewObject(jp.curIndex)
		obj.lineNumber = len(jp.lines)

		return jp.parseObject(obj)
	} else if jp.curTokenType == BEGIN_ARRAY {
		val.SetValueType(V_ARRAY)

		return jp.parseArray(val)
	} else {
		if jp.curTokenType == J_TRUE {
			val.SetValueType(V_BOOLEAN)
			val.SetValue(true)
		} else if jp.curTokenType == J_FALSE {
			val.SetValueType(V_BOOLEAN)
			val.SetValue(false)
		} else if jp.curTokenType == J_NULL {
			val.SetValueType(V_NULL)
			val.SetValue(jp.curTokenVar)
		} else if jp.curTokenType == NUMBER {
			val.SetValueType(V_NUMBER)
			// numbers stored as strings, don't know if int or float yet
			val.SetValue(jp.curTokenVar)
		} else if jp.curTokenType == STRING {
			val.SetValueType(V_STRING)
			val.SetValue(jp.curTokenVar)
		} else if jp.curTokenType == REF {
			val.SetValueType(V_STRING)
			val.SetValue("$ref")
		}
	}

	return true
}

//   A JSON Reference is a JSON object, which contains a member named
//   "$ref", which has a JSON string value.  Example:
//
//   { "$ref": "http://example.com/example.json#/foo/bar" }
//
//   If a JSON value does not have these characteristics, then it SHOULD
//   NOT be interpreted as a JSON Reference.
//
//   The "$ref" string value contains a URI [RFC3986], which identifies
//   the location of the JSON value being referenced.  It is an error
//   condition if the string value does not conform to URI syntax rules.
//   Any members other than "$ref" in a JSON Reference object SHALL be
//   ignored.
func (jp *JSONParser) parseRefPtr(ref *JSONNode) bool {
	err := jp.expectToken(STRING | BEGIN_OBJECT)
	if err != nil {
		return false
	}

	// definition of $ref, not actual $ref
	if jp.curTokenType == BEGIN_OBJECT {
		ref.SetType(N_OBJECT)
		ref.SetValueType(V_OBJECT)
		obj := ref.NewObject(jp.curIndex)
		obj.lineNumber = len(jp.lines)

		return jp.parseObject(obj)
	}

	ref.SetValueType(V_STRING)
	ref.SetValue(jp.curTokenVar)

	jp.references[jp.curTokenVar] = ref

	return true
}

// 2.2.  Objects
// 
//    An object structure is represented as a pair of curly brackets
//    surrounding zero or more name/value pairs (or members).  A name is a
//    string.  A single colon comes after each name, separating the name
//    from the value.  A single comma separates a value from a following
//    name.  The names within an object SHOULD be unique.
// 
//       object = begin-object [ member *( value-separator member ) ]
//       end-object
// 
func (jp *JSONParser) parseObject(obj *JSONNode) bool {
	if jp.curTokenType != BEGIN_OBJECT {
		Trace.Println("begin object error")
		panic("begin object error")
		return false
	}

	for {
		err := jp.expectToken(STRING | END_OBJECT | REF)
		if err != nil {
			Trace.Println("expect token error")
			return false
		}

		if jp.curTokenType == END_OBJECT {
			break
		}

		mem := obj.NewMember(jp.curTokenVar, jp.curIndex)
		mem.lineNumber = len(jp.lines)

		if !jp.parseMember(mem) {
			Trace.Println("parse mem failed")
			return false
		}

		err = jp.expectToken(VALUE_SEPARATOR | END_OBJECT)
		if err != nil {
			Trace.Println("expect token error 2")
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

// expect one or more tokens as the next item in the json stream
func (jp *JSONParser) expectToken(valTokens int) *ParseError {
	token := jp.nextToken()

	if token&valTokens == token {
		return nil
	}

	return jp.addError("Unexpected token: Expecting " + tokenToString(valTokens) + " - Encountered " + tokenToString(token), JP_ERROR)
}

// retrieve the next token in the tokenized json stream
func (jp *JSONParser) nextToken() int {
	jp.curIndex++
	jp.curTokenType = jp.tokens[jp.curIndex]

	if (jp.curTokenType == STRING) || (jp.curTokenType == NUMBER) {
		jp.curIndex++
		jp.curTokenVar = jp.variables[jp.tokens[jp.curIndex]]
	} else if jp.curTokenType == REF {
		jp.curTokenVar = "$ref"
	}

	if jp.curTokenType == BEGIN_ARRAY || jp.curTokenType == BEGIN_OBJECT || jp.curTokenType == VALUE_SEPARATOR {
		jp.newLine(1, true)
	} else if jp.curTokenType == END_ARRAY || jp.curTokenType == END_OBJECT {
		jp.newLine(-1, false)
	}

	return jp.curTokenType
}

// Retrieve the next word in the source json source
func (jp *JSONParser) getWord() int {
	var letter	string

	// skip white space
	for {
		if jp.ltrIndex >= len(jp.raw)  {
			return END_OF_SOURCE
		}

		letter = jp.raw[jp.ltrIndex:jp.ltrIndex+1]
		jp.ltrIndex++

		if strings.Contains(" \n\t", letter)  {
			continue
		}

		break
	}

	jp.curTokenVar = letter
	if letter == "\"" {
		endQuote := 0
		foundEnd := false
		for i := jp.ltrIndex; i < len(jp.raw); i++ {
			letter = jp.raw[i:i+1]
			// skip escpaped characters
			if letter == `\` {
				i++
			} else if letter == "\"" {
				endQuote = i
				foundEnd = true
				break
			}
		}

		jp.curTokenVar = jp.raw[jp.ltrIndex : endQuote]
		jp.ltrIndex = endQuote + 1
		if jp.curTokenVar == "$ref" {
			return REF
		} else if !foundEnd {
			jp.curTokenVar = "error"
			jp.addError("No matching end quote", JP_FATAL)
			return END_OF_SOURCE
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
		endWord := strings.IndexAny(jp.raw[jp.ltrIndex:], " ,:}]\n")
		jp.curTokenVar = jp.raw[jp.ltrIndex-1: jp.ltrIndex + endWord]

		jp.ltrIndex += endWord
		
		if validNum.MatchString(jp.curTokenVar) {
			return NUMBER
		} else {
			jp.addError("Invalid number format", JP_ERROR)
			return NUMBER
		}
	}

	if strings.Contains("tfn", letter) {
		// potential boolean or null - need to validate
		endWord := strings.IndexAny(jp.raw[jp.ltrIndex:], " ,:}]\n")
		if endWord > 0 {
			jp.curTokenVar = jp.raw[jp.ltrIndex - 1 : jp.ltrIndex + endWord]
			jp.ltrIndex += endWord

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

// keep track of the lines for the json tokens used to reproduce
// json for error display
func (jp *JSONParser) newLine(indent int, tokenOnLine bool) {
	var newLine		lineItem

	if tokenOnLine {
		jp.lines[jp.lineCount].tokenEnd = jp.curIndex
		newLine.tokenStart = jp.curIndex + 1
	} else {
		jp.lines[jp.lineCount].tokenEnd = jp.curIndex - 1
		newLine.tokenStart = jp.curIndex
	}

	jp.curIndent += indent
	newLine.indent = jp.curIndent

	jp.lineCount++
	jp.lines = append(jp.lines, newLine)
}

// add an error to the list of errors encountered during parsing
func (jp *JSONParser) addError(errText string, level int) *ParseError {
	var pError	ParseError

	pError.Error = errText
	pError.ErrorLevel = level
	pError.Offset = jp.curIndex
	pError.LineNumber = len(jp.lines)

	Trace.Println(errText)
	jp.errorList = append(jp.errorList, pError)

	return &pError
}

// convert one or more tokens to a string representation separated by |
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
