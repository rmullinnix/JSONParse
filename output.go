package JSONParse

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func outputInit(log_level string) {

	traceHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warningHandle := ioutil.Discard
	errorHandle := os.Stderr

	if log_level == "trace" {
		traceHandle = os.Stdout
		infoHandle = os.Stdout
		warningHandle = os.Stdout
	} else if log_level == "info" {
		infoHandle = os.Stdout
		warningHandle = os.Stdout
	} else if log_level == "warn" {
		warningHandle = os.Stdout
	}

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Lshortfile)
	Warning = log.New(warningHandle,
		"WARNING: ", 0)

	Error = log.New(errorHandle,
		"ERROR: ", 0)
}

func OutputError(node *JSONNode, errMsg string) {
	tokenIndex := node.tokenIndex
	parser := node.root.doc
	tokenStart := 0
	tokenEnd := len(parser.tokens)

	if tokenIndex > 10 {
		tokenStart = tokenIndex - 10
	}

	if tokenIndex < tokenEnd - 10 {
		tokenEnd = tokenIndex + 10
	}

	output := parser.prettyTokens(tokenStart, tokenEnd)

	Error.Fatalln(errMsg + "\n" + output)
}

// formats the json with newlines and indentation
func (jp *JSONParser) Pretty() {
	Info.Println(jp.prettyTokens(0, len(jp.tokens)))
}

func (jp *JSONParser) prettyTokens(start int, end int) string {
	indent := 0
	output := ""
	for i := start; i < end; i++ {
		if jp.tokens[i] == BEGIN_OBJECT {
			output += "{\n"
			indent++
			output += strings.Repeat("\t", indent)
		} else if jp.tokens[i] == END_OBJECT {
			if indent > 0 {
				indent--
			}
			output += "\n"
			output += strings.Repeat("\t", indent)
			output += `}`
		} else if jp.tokens[i] == BEGIN_ARRAY {
			output += "[\n"
			indent++
			output += strings.Repeat("\t", indent)
		} else if jp.tokens[i] == END_ARRAY {
			if indent > 0 {
				indent--
			}
			output += "\n"
			output += strings.Repeat("\t", indent)
			output += `]`
		} else if jp.tokens[i] == VALUE_SEPARATOR {
			output += ",\n"
			output += strings.Repeat("\t", indent)
		} else if jp.tokens[i] == J_FALSE {
			output += `false`
		} else if jp.tokens[i] == J_TRUE {
			output += `true`
		} else if jp.tokens[i] == J_NULL {
			output += `null`
		} else if jp.tokens[i] == STRING {
			i++
			output += `"` + jp.variables[jp.tokens[i]] + `"`
		} else if jp.tokens[i] == REF {
			output += `"$ref"`
		} else if jp.tokens[i] == NUMBER {
			i++
			output += jp.variables[jp.tokens[i]]
		} else if jp.tokens[i] == NAME_SEPARATOR {
			output += `: `
		}
	}

	return output
}
