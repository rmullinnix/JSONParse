package JSONParse

import (
	"strings"
	"strconv"
)

// multipleOf
//
// 5.1.1.1.  Valid values
//
// The value of "multipleOf" MUST be a JSON number. This number MUST be 
// strictly greater than 0.
//
// 5.1.1.2.  Conditions for successful validation
//
// A numeric instance is valid against "multipleOf" if the result of the 
// division of the instance by this keyword's value is an integer.
//
func validMultipleOf(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validMultipleOf")
	if mem.GetValueType() != V_NUMBER {
		Trace.Println(stack_id, "validMultipleOf", true)
		return true
	}

	strDocNum := mem.GetValue().(string)
	strSchemaNum := schema.GetValue().(string)

	if strings.Index(strDocNum, ".") > -1 || strings.Index(strSchemaNum, ".") > -1 {
		Trace.Println("  validMutlipleOf() - with float")
		_, dErr := strconv.ParseFloat(strDocNum, 64)
		if dErr != nil {
			errs.Add(mem, "Invalid number in document " + strDocNum, JP_ERROR)
		}
		_, iErr := strconv.ParseFloat(strSchemaNum, 64)
		if iErr != nil {
			errs.Add(mem, "Invalid number in schema " + strSchemaNum, JP_WARNING)
		}

		docLen := len(strDocNum) - strings.Index(strDocNum, ".")
		schemaLen := len(strSchemaNum) - strings.Index(strSchemaNum, ".")

		var newDocNum		string
		var newSchemaNum	string

		if schemaLen > docLen {
			newDocNum = strings.Replace(strDocNum, ".", "", 1) + strings.Repeat("0", schemaLen-docLen)
			newSchemaNum = strings.Replace(strDocNum, ".", "", 1)
		} else {
			newDocNum = strings.Replace(strDocNum, ".", "", 1)
			newSchemaNum = strings.Replace(strDocNum, ".", "", 1) + strings.Repeat("0", docLen-schemaLen)
		}

		iDocNum, _ := strconv.Atoi(newDocNum)
		iSchemaNum, _ := strconv.Atoi(newSchemaNum)

		if iDocNum == 0 {
			Trace.Println(stack_id, "validMultipleOf", true)
			return true
		}

		rem := iDocNum % iSchemaNum
		if rem != 0 {
			Trace.Println("Number", strDocNum, "is not multipleOf", strSchemaNum, rem)
			errs.Add(mem, "Number " + strDocNum + " is not multipleOf " + strSchemaNum, JP_ERROR)
			Trace.Println(stack_id, "validMultipleOf", false)
			return false
		}
	 } else {
		Trace.Println("  validMutlipleOf() - with int")
		iDocNum, dErr := strconv.Atoi(strDocNum)
		if dErr != nil {
			errs.Add(mem, "Invalid number in document " + strDocNum, JP_ERROR)
		}
		iSchemaNum, iErr := strconv.Atoi(strSchemaNum)
		if iErr != nil {
			errs.Add(mem, "Invalid number in schema " + strSchemaNum, JP_WARNING)
		}

		if iDocNum == 0 {
			Trace.Println(stack_id, "validMultipleOf", true)
			return true
		}

		rem := iDocNum % iSchemaNum

		if rem != 0 {
			errs.Add(mem, "Number " + strDocNum + " is not multipleOf " + strSchemaNum, JP_ERROR)
			Trace.Println(stack_id, "validMultipleOf", false)
			return false
		}
	}

	Trace.Println(stack_id, "validMultipleOf", true)
	return true
}

// 5.1.2.  maximum and exclusiveMaximum
//
// 5.1.2.1.  Valid values
//
// The value of "maximum" MUST be a JSON number. The value of 
// "exclusiveMaximum" MUST be a boolean.
//
// If "exclusiveMaximum" is present, "maximum" MUST also be present.
//
// 5.1.2.2.  Conditions for successful validation
//
// Successful validation depends on the presence and value of "exclusiveMaximum":
//
// if "exclusiveMaximum" is not present, or has boolean value false, then 
// the instance is valid if it is lower than, or equal to, the value of "maximum";
//
// if "exclusiveMaximum" has boolean value true, the instance is valid if it 
// is strictly lower than the value of "maximum".
//
// 5.1.2.3.  Default value
//
// "exclusiveMaximum", if absent, may be considered as being present with 
// boolean value false.
//
func validMaximum(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validMaximum")

	if val, found := Mutex.Find(stack_id[:strings.LastIndex(stack_id, ".")] + "maximum"); found {
		Trace.Println(stack_id, "validProperties -- mutex --", val)
		return val
	} else {
		Mutex.Add(stack_id[:strings.LastIndex(stack_id, ".")] + "maximum")
	}

	doc := mem
	if doc.GetValueType() != V_NUMBER {
		Warning.Println("valid max against non number")
		Trace.Println(stack_id, "validMaximum", true)

		Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "maximum", true)

		return true
	}

	strDocNum := mem.GetValue().(string)
	strSchemaMax := ""

	hasMax := false

       	if item, found := parent.Find("maximum"); found {
		strSchemaMax = item.GetValue().(string)
		hasMax = true
	}

	eMax := false
       	if item, found := parent.Find("exclusiveMaximum"); found {
		if !hasMax {
			errs.Add(mem, "exclusiveMaximum is present without correspoding maximum", JP_ERROR)
			Trace.Println(stack_id, "validMaximum", false)
			Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "maximum", false)
			return false
		}
		eMax = item.GetValue().(bool)
	}

	if strings.Index(strDocNum, ".") > -1 {
		fDocNum, dErr := strconv.ParseFloat(strDocNum, 64)
		if dErr != nil {
			errs.Add(mem, "Invalid number in document " + strDocNum, JP_ERROR)
		}
		fSchemaMax, iErr := strconv.ParseFloat(strSchemaMax, 64)
		if iErr != nil {
			errs.Add(mem, "Invalid number in schema " + strSchemaMax, JP_WARNING)
		}

		if eMax {
			if fDocNum < fSchemaMax {
				return true
			} else {
				errs.Add(mem, "Document number " + strDocNum + " is not less than maximum " + strSchemaMax, JP_ERROR)
			}
		} else if fDocNum <= fSchemaMax {
			Trace.Println(stack_id, "validMaximum", true)
			Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "maximum", true)
			return true
		} else {
			errs.Add(mem, "Document number " + strDocNum + " is not less than or equal to maximum " + strSchemaMax, JP_ERROR)
		}
	 } else {
		iDocNum, dErr := strconv.Atoi(strDocNum)
		if dErr != nil {
			errs.Add(mem, "Invalid number in document " + strDocNum, JP_ERROR)
		}
		iSchemaMax, iErr := strconv.Atoi(strSchemaMax)
		if iErr != nil {
			errs.Add(mem, "Invalid number in schema " + strSchemaMax, JP_WARNING)
		}

		if eMax {
			if iDocNum < iSchemaMax {
				Trace.Println(stack_id, "validMaximum", true)
				Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "maximum", true)
				return true
			} else {
				errs.Add(mem, "Document number " + strDocNum + " is not less than maximum " + strSchemaMax, JP_ERROR)
			}
		} else if iDocNum <= iSchemaMax {
			Trace.Println(stack_id, "validMaximum", true)
			Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "maximum", true)
			return true
		} else {
			errs.Add(mem, "Document number " + strDocNum + " is not less than or equal to maximum " + strSchemaMax, JP_ERROR)
		}
	}

	Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "maximum", false)

	Trace.Println(stack_id, "validMaximum", false)
	return false
}

// 5.1.3.  minimum and exclusiveMinimum
//
// 5.1.3.1.  Valid values
//
// The value of "minimum" MUST be a JSON number. The value of "exclusiveMinimum" MUST be a boolean.
//
// If "exclusiveMinimum" is present, "minimum" MUST also be present.
//
// 5.1.3.2.  Conditions for successful validation
//
// Successful validation depends on the presence and value of "exclusiveMinimum":
//
// if "exclusiveMinimum" is not present, or has boolean value false, then the instance is valid if it is greater than, or equal to, the value of "minimum";
//
// if "exclusiveMinimum" is present and has boolean value true, the instance is valid if it is strictly greater than the value of "minimum".
//
// 5.1.3.3.  Default value
//
// "exclusiveMinimum", if absent, may be considered as being present with boolean value false.
//
func validMinimum(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	Trace.Println(stack_id, "validMinimum")

	if val, found := Mutex.Find(stack_id[:strings.LastIndex(stack_id, ".")] + "minimum"); found {
		Trace.Println(stack_id, "validProperties -- mutex --", val)
		return val
	} else {
		Mutex.Add(stack_id[:strings.LastIndex(stack_id, ".")] + "minimum")
	}

	doc := mem
	if doc.GetValueType() != V_NUMBER {
		Warning.Println("valid max against non number")
		Trace.Println(stack_id, "validMinimum", true)

		Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "minimum", true)

		return true
	}
	
	strDocNum := mem.GetValue().(string)
	strSchemaMin := ""

	hasMin := false

       	if item, found := parent.Find("minimum"); found {
		strSchemaMin = item.GetValue().(string)
		Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "minimum", true)
		hasMin = true
	}

	eMin := false
       	if item, found := parent.Find("exclusiveMinimum"); found {
		if !hasMin {
			errs.Add(mem, "exclusiveMinium is present without correspoding minimum", JP_ERROR)
			Trace.Println(stack_id, "validMinimum", false)
			Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "minimum", false)
			return false
		}
		eMin = item.GetValue().(bool)
	}

	if strings.Index(strDocNum, ".") > -1 {
		fDocNum, dErr := strconv.ParseFloat(strDocNum, 64)
		if dErr != nil {
			errs.Add(mem, "Invalid number in document " + strDocNum, JP_ERROR)
		}
		fSchemaMin, iErr := strconv.ParseFloat(strSchemaMin, 64)
		if iErr != nil {
			errs.Add(mem, "Invalid number in schema " + strSchemaMin, JP_WARNING)
		}

		if eMin {
			if fDocNum > fSchemaMin {
				Trace.Println(stack_id, "validMinimum", true)
				Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "minimum", true)
				return true
			} else {
				errs.Add(mem, "Document number " + strDocNum + " is not greater than minimum " + strSchemaMin, JP_ERROR)
			}
		} else if fDocNum >= fSchemaMin {
			Trace.Println(stack_id, "validMinimum", true)
			Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "minimum", true)
			return true
		} else {
			errs.Add(mem, "Document number " + strDocNum + " is not greater than or equal to minimum " + strSchemaMin, JP_ERROR)
		}
	 } else {
		iDocNum, dErr := strconv.Atoi(strDocNum)
		if dErr != nil {
			errs.Add(mem, "Invalid number in document " + strDocNum, JP_ERROR)
		}
		iSchemaMin, iErr := strconv.Atoi(strSchemaMin)
		if iErr != nil {
			errs.Add(mem, "Invalid number in schema " + strSchemaMin, JP_WARNING)
		}

		if eMin {
			if iDocNum > iSchemaMin {
				Trace.Println(stack_id, "validMinimum", true)
				Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "minimum", true)
				return true
			} else {
				errs.Add(mem, "Document number " + strDocNum + " is not greater than minimum " + strSchemaMin, JP_ERROR)
			}
		} else if iDocNum >= iSchemaMin {
			Trace.Println(stack_id, "validMinimum", true)
			Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "minimum", true)
			return true
		} else {
			errs.Add(mem, "Document number " + strDocNum + " is not greater than or equal to minimum " + strSchemaMin, JP_ERROR)
		}
	}

	Mutex.Set(stack_id[:strings.LastIndex(stack_id, ".")] + "minimum", false)

	Trace.Println(stack_id, "validMinimum", false)
	return false
}
