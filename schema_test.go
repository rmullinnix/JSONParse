//
package JSONParse

import (
	"testing"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

var globalTestScope *testing.T

func TestInit(t *testing.T) {
	globalTestScope = t
	schemaErrors = NewSchemaErrors()

	keywords = make(map[string]validator)

	keywords["type"] = validType
	keywords["enum"] = validEnum
	keywords["required"] = validRequired
	keywords["properties"] = validProperties
	keywords["additionalProperties"] = validProperties
	keywords["patternProperties"] = validProperties
	keywords["minProperties"] = validProperties
	keywords["maxProperties"] = validProperties
	keywords["items"] = validItems
	keywords["uniqueItems"] = validUnique
	keywords["maxItems"] = validMaxItems
	keywords["minItems"] = validMinItems
	keywords["additionalItems"] = validAdditionalItems
	keywords["maxLength"] = validMaxLength
	keywords["minLength"] = validMinLength
	keywords["maximum"] = validMaximum
	keywords["exclusiveMaximum"] = validMaximum
	keywords["minimum"] = validMinimum
	keywords["exclusiveMinimum"] = validMinimum
	keywords["pattern"] = validPattern
	keywords["anyOf"] = validAnyOf
	keywords["allOf"] = validAllOf
	keywords["oneOf"] = validOneOf
	keywords["multipleOf"] = validMultipleOf
	keywords["default"] = validDefault
	keywords["not"] = validNot
	keywords["dependencies"] = validDependencies
	keywords["format"] = validFormat

	regexHostname = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)` +
					`*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	regexDateTime = regexp.MustCompile(`^([0-9]{4})-([0-9]{2})-([0-9]{2})` +
					`([Tt]([0-9]{2}):([0-9]{2}):([0-9]{2})(\.[0-9]+)?)?` +
					`([Tt]([0-9]{2}):([0-9]{2}):([0-9]{2})(\\.[0-9]+)?)?` +
					`(([Zz]|([+-])([0-9]{2}):([0-9]{2})))?`)
	regexEmail = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)

	validateFormat = true
}

func TestAdditionalItems(t *testing.T) {
	genTestFile("additionalItems")
	executeTests("additionalItems", t)
}

func TestAdditionalProperties(t *testing.T) {
	genTestFile("additionalProperties")
	executeTests("additionalProperties", t)
}

func TestAllOf(t *testing.T) {
	genTestFile("allOf")
	executeTests("allOf", t)
}

func TestAnyOf(t *testing.T) {
	genTestFile("anyOf")
	executeTests("anyOf", t)
}

func TestDefault(t *testing.T) {
	genTestFile("default")
	executeTests("default", t)
}

func TestDefinitions(t *testing.T) {
	genTestFile("definitions")
	executeTests("definitions", t)
}

func TestDependencies(t *testing.T) {
	genTestFile("dependencies")
	executeTests("dependencies", t)
}

func TestEnum(t *testing.T) {
	genTestFile("enum")
	executeTests("enum", t)
}

func TestItems(t *testing.T) {
	genTestFile("items")
	executeTests("items", t)
}

func TestMaxItems(t *testing.T) {
	genTestFile("maxItems")
	executeTests("maxItems", t)
}

func TestMaxLength(t *testing.T) {
	genTestFile("maxLength")
	executeTests("maxLength", t)
}

func TestMaxProperties(t *testing.T) {
	genTestFile("maxProperties")
	executeTests("maxProperties", t)
}

func TestMaximum(t *testing.T) {
	genTestFile("maximum")
	executeTests("maximum", t)
}

func TestMinItems(t *testing.T) {
	genTestFile("minItems")
	executeTests("minItems", t)
}

func TestMinLength(t *testing.T) {
	genTestFile("minLength")
	executeTests("minLength", t)
}

func TestMinProperties(t *testing.T) {
	genTestFile("minProperties")
	executeTests("minProperties", t)
}

func TestMinimum(t *testing.T) {
	genTestFile("minimum")
	executeTests("minimum", t)
}

func TestMultipleOf(t *testing.T) {
	genTestFile("multipleOf")
	executeTests("multipleOf", t)
}

func TestNot(t *testing.T) {
	genTestFile("not")
	executeTests("not", t)
}

func TestOneOf(t *testing.T) {
	genTestFile("oneOf")
	executeTests("oneOf", t)
}

func TestPattern(t *testing.T) {
	genTestFile("pattern")
	executeTests("pattern", t)
}

func TestPatternProperties(t *testing.T) {
	genTestFile("patternProperties")
	executeTests("patternProperties", t)
}

func TestRef(t *testing.T) {
	genTestFile("ref")
	chunkSchema("ref", t)
}

//func TestRemoteRef(t *testing.T) {
//	genTestFile("refRemote")
//	executeTests("refRemote", t)
//}

func TestRequired(t *testing.T) {
	genTestFile("required")
	executeTests("required", t)
}

func TestType(t *testing.T) {
	genTestFile("type")
	executeTests("type", t)
}

func TestUniqueItems(t *testing.T) {
	genTestFile("uniqueItems")
	executeTests("uniqueItems", t)
}

// test optional items
func TestFormat(t *testing.T) {
	genTestFile("optional/format")
	executeTests("optional/format", t)
}

func genTestFile(name string) {
	stream, _ := loadDoc("https://raw.githubusercontent.com/json-schema/JSON-Schema-Test-Suite/develop/tests/draft4/" + name + ".json")

	json := `{ "testcases": ` + string(stream) + "}"

	ioutil.WriteFile("tests/" + name + ".json", []byte(json), os.ModePerm)
}

func executeTests(name string, t *testing.T) {
	jp := NewJSONParser("tests/" + name + ".json", 5, "error")

	jp.Parse()

	tree := jp.GetDoc()

	testcases, _ := tree.Find("testcases")
	testcases.ResetIterate()
	test_major := 1
	for {
		testcase := testcases.GetNext()
		if testcase == nil {
			break
		}

		testcase.ResetIterate()
		item := testcase.GetNext()
		if item == nil {
			break
		}

		schema, found := item.Find("schema")
		if found {
			schema.ResetIterate()
			schema_itm := schema.GetNext()

			tests, tstfound := item.Find("tests")
			if tstfound {
				tests.ResetIterate()
				test := tests.GetNext()

				test_minor := 1
				for {
					Mutex = NewSchemaMutex()

					if test == nil {
						break
					}

					test.ResetIterate()
					tmp := test.GetNext()
					tst_item, _ := tmp.Find("data")

					valid := validMember("*", name, tst_item, schema_itm)

					exp_item, _ := tmp.Find("valid")
					exp_result := exp_item.GetValue().(bool)
					
					test_num := strconv.Itoa(test_major) + "." + strconv.Itoa(test_minor)
					fmt.Println(test_num, "data: ", tst_item.GetJson())
					fmt.Println(test_num, "result:", valid, "vs expected:", exp_result)

					if valid != exp_result {
						t.Fail()
						t.Log(name, test_num, "failed:  expecting", exp_result, "received", valid)
					}
					test = tests.GetNext()
					test_minor++
				}
			}
		}
		test_major++
	}
}

func chunkSchema(name string, t *testing.T) {

	jp := NewJSONParser("tests/" + name + ".json", 5, "error")

	jp.ResolveRefs(false)

	jp.Parse()

	Mutex = NewSchemaMutex()

	tree := jp.GetDoc()

	testcases, _ := tree.Find("testcases")
	testcases.ResetIterate()
	index := 1
	for {
		testcase := testcases.GetNext()
		if testcase == nil {
			break
		}

		testcase.ResetIterate()
		item := testcase.GetNext()
		if item == nil {
			break
		}

		schema, _ := item.Find("schema")
		schema.ResetIterate()
		schema_doc := schema.GetNext()

		sfn := "tests/" + name + strconv.Itoa(index) + "_schema.json"
		ioutil.WriteFile(sfn, []byte(schema_doc.GetJson()), os.ModePerm)

		delete(item.namedKids, "schema")

		jp1 := NewJSONParser(sfn, 5, "error")
		jp1.Parse()
		schema_itm := jp1.GetDoc()

		tests, tstfound := item.Find("tests")
		if tstfound {
			tests.ResetIterate()
			test := tests.GetNext()

			test_minor := 1
			for {
				if test == nil {
					break
				}

				Mutex = NewSchemaMutex()

				test.ResetIterate()
				tmp := test.GetNext()
				tst_item, _ := tmp.Find("data")

				valid := validMember("*", name, tst_item, schema_itm)

				exp_item, _ := tmp.Find("valid")
				exp_result := exp_item.GetValue().(bool)
					
				test_num := strconv.Itoa(index) + "." + strconv.Itoa(test_minor)
				fmt.Println(test_num, "data: ", tst_item.GetJson())
				fmt.Println(test_num, "result:", valid, "vs expected:", exp_result)

				if valid != exp_result {
					t.Fail()
					t.Log(name, test_num, "failed:  expecting", exp_result, "received", valid)
				}
				test = tests.GetNext()
				test_minor++
			}
		}
		index++
		if index > 2 {
			break
		}
	}
}
