package main

import (
	"fmt"
	"github.com/rmullinnix/JSONParse"
	"flag"
//	"strings"
)

func main() {
	filePtr := flag.String("file", "security.json", "file to parse")
	errorPtr := flag.Int("maxerror", 5, "maximum number of errors before abort")
	schemaPtr := flag.String("schema", "", "schema to use for validation")
	levelPtr := flag.String("level", "error", "level of logging (trace, info, warning, error)")
//	schemaPtr := flag.String("schema", "http://swagger.io/v2/schema.json#", "schema to use for validation")

	flag.Parse()

	file := *filePtr
	maxError := *errorPtr
	schemaFile := *schemaPtr
	level := *levelPtr

	fmt.Println("file", file)
	var parser	*JSONParse.JSONParser
	if len(file) > 0 {
		parser = JSONParse.NewJSONParser(file, maxError, level)
		valid, errors := parser.Parse()
		if !valid {
			fmt.Println(errors)
		}
	}

	if len(schemaFile) > 0 {
		fmt.Println("parse schema")
		schema := JSONParse.NewJSONSchema(schemaFile, level)
		fmt.Println("validate file")
		schema.ValidateDocument(file)
	}
	
	if flag.Arg(0) == "pretty" {
		fmt.Println(parser.Pretty())
	}
//	schema := JSONParse.NewJSONParser(schemaFile, maxError, level)
//	valid, errors = schema.Parse()

//	JSONParse.ValidateAgainstSchema(parser.GetDoc(), schema.GetDoc())

//	fmt.Println(parser.Pretty())
}
