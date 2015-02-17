package JSONParse

import (
	"math"
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
func validMultipleOf(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	strDocNum := mem.GetValue().(string)
	strSchemaNum := schema.GetValue().(string)

	if strings.Index(strDocNum, ".") > -1 {
		Trace.Println("  validMutlipleOf() - with float")
		fDocNum, dErr := strconv.ParseFloat(strDocNum, 64)
		if dErr != nil {
			OutputError(mem, "Invalid number in document " + strDocNum)
		}
		fSchemaNum, iErr := strconv.ParseFloat(strSchemaNum, 64)
		if iErr != nil {
			OutputError(mem, "Invalid number in schema " + strSchemaNum)
		}

		rem := math.Remainder(fDocNum, fSchemaNum)
		if rem != 0 {
			OutputError(mem, "Number " + strDocNum + " is not multipleOf " + strSchemaNum)
			return false
		}
	 } else {
		Trace.Println("  validMutlipleOf() - with int")
		iDocNum, dErr := strconv.Atoi(strDocNum)
		if dErr != nil {
			OutputError(mem, "Invalid number in document " + strDocNum)
		}
		iSchemaNum, iErr := strconv.Atoi(strSchemaNum)
		if iErr != nil {
			OutputError(mem, "Invalid number in schema " + strSchemaNum)
		}

		rem := iDocNum % iSchemaNum

		if rem != 0 {
			OutputError(mem, "Number " + strDocNum + " is not multipleOf " + strSchemaNum)
			return false
		}
	}

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
func validMaximum(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	doc := mem
	if doc.GetValueType() != V_NUMBER {
		Warning.Println("valid max against non number")
		return true
	}

	Trace.Println("  validMaximum()")
	if doc.GetState() == NODE_MUTEX {
		return true
	} else  {
		doc.SetState(NODE_MUTEX)
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
			OutputError(mem, "exclusiveMaximum is present without correspoding maximum")
			return false
		}
		eMax = item.GetValue().(bool)
	}

	if strings.Index(strDocNum, ".") > -1 {
		fDocNum, dErr := strconv.ParseFloat(strDocNum, 64)
		if dErr != nil {
			OutputError(mem, "Invalid number in document " + strDocNum)
		}
		fSchemaMax, iErr := strconv.ParseFloat(strSchemaMax, 64)
		if iErr != nil {
			OutputError(mem, "Invalid number in schema " + strSchemaMax)
		}

		if eMax {
			if fDocNum < fSchemaMax {
				return true
			} else {
				OutputError(mem, "Document number " + strDocNum + " is not less than maximum " + strSchemaMax)
			}
		} else if fDocNum <= fSchemaMax {
			return true
		} else {
			OutputError(mem, "Document number " + strDocNum + " is not less than or equal to maximum " + strSchemaMax)
		}
	 } else {
		iDocNum, dErr := strconv.Atoi(strDocNum)
		if dErr != nil {
			OutputError(mem, "Invalid number in document " + strDocNum)
		}
		iSchemaMax, iErr := strconv.Atoi(strSchemaMax)
		if iErr != nil {
			OutputError(mem, "Invalid number in schema " + strSchemaMax)
		}

		if eMax {
			if iDocNum < iSchemaMax {
				return true
			} else {
				OutputError(mem, "Document number " + strDocNum + " is not less than maximum " + strSchemaMax)
			}
		} else if iDocNum <= iSchemaMax {
			return true
		} else {
			OutputError(mem, "Document number " + strDocNum + " is not less than or equal to maximum " + strSchemaMax)
		}
	}

	return true
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
func validMinimum(mem *JSONNode, schema *JSONNode, parent *JSONNode) bool {
	doc := mem
	if doc.GetValueType() != V_NUMBER {
		Warning.Println("valid max against non number")
		return true
	}
	

	Trace.Println("  validMinimum()")
	if doc.GetState() == NODE_MUTEX {
		return true
	} else  {
		doc.SetState(NODE_MUTEX)
	}

	strDocNum := mem.GetValue().(string)
	strSchemaMin := ""

	hasMin := false

       	if item, found := parent.Find("minimum"); found {
		strSchemaMin = item.GetValue().(string)
		hasMin = true
	}

	eMin := false
       	if item, found := parent.Find("exclusiveMinimum"); found {
		if !hasMin {
			OutputError(mem, "exclusiveMinium is present without correspoding minimum")
			return false
		}
		eMin = item.GetValue().(bool)
	}

	if strings.Index(strDocNum, ".") > -1 {
		fDocNum, dErr := strconv.ParseFloat(strDocNum, 64)
		if dErr != nil {
			OutputError(mem, "Invalid number in document " + strDocNum)
		}
		fSchemaMin, iErr := strconv.ParseFloat(strSchemaMin, 64)
		if iErr != nil {
			OutputError(mem, "Invalid number in schema " + strSchemaMin)
		}

		if eMin {
			if fDocNum > fSchemaMin {
				return true
			} else {
				OutputError(mem, "Document number " + strDocNum + " is not greater than minimum " + strSchemaMin)
			}
		} else if fDocNum >= fSchemaMin {
			return true
		} else {
			OutputError(mem, "Document number " + strDocNum + " is not greater than or equal to minimum " + strSchemaMin)
		}
	 } else {
		iDocNum, dErr := strconv.Atoi(strDocNum)
		if dErr != nil {
			OutputError(mem, "Invalid number in document " + strDocNum)
		}
		iSchemaMin, iErr := strconv.Atoi(strSchemaMin)
		if iErr != nil {
			OutputError(mem, "Invalid number in schema " + strSchemaMin)
		}

		if eMin {
			if iDocNum > iSchemaMin {
				return true
			} else {
				OutputError(mem, "Document number " + strDocNum + " is not greater than minimum " + strSchemaMin)
			}
		} else if iDocNum >= iSchemaMin {
			return true
		} else {
			OutputError(mem, "Document number " + strDocNum + " is not greater than or equal to minimum " + strSchemaMin)
		}
	}

	return true
}
