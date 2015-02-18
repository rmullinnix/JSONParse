//
package JSONParse

import (
	"testing"
	"fmt"
	"io/ioutil"
	"os"
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

func TestEnum(t *testing.T) {
	genTestFile("enum")
	executeTests("enum", t)
}

func TestMaxLength(t *testing.T) {
	genTestFile("maxLength")
	executeTests("maxLength", t)
}

func TestMinLength(t *testing.T) {
	genTestFile("minLength")
	executeTests("minLength", t)
}

func TestItems(t *testing.T) {
	genTestFile("items")
	executeTests("items", t)
}

func TestMaxItems(t *testing.T) {
	genTestFile("maxItems")
	executeTests("maxItems", t)
}

func TestMinItems(t *testing.T) {
	genTestFile("minItems")
	executeTests("minItems", t)
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

				for {
					if test == nil {
						break
					}

					test.ResetIterate()
					tmp := test.GetNext()
					tst_item, _ := tmp.Find("data")

					valid := validMember(name, tst_item, schema_itm, false)

					exp_item, _ := tmp.Find("valid")
					exp_result := exp_item.GetValue().(bool)
					
					fmt.Println("data: ", tst_item.GetJson())
					fmt.Println("Test Result: ", valid)

					if valid != exp_result {
						t.Fail()
						t.Log("Test failed:  expecting", exp_result, "received", valid)
					}
					test = tests.GetNext()
				}
			}
		}
	}
}
