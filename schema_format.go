package JSONParse

import (
	"net"
	"net/url"
	"strings"
)

// 7.2.  Implementation requirements

// Implementations MAY support the "format" keyword. Should they choose to do so:
//
// they SHOULD implement validation for attributes defined below;
//
// they SHOULD offer an option to disable validation for this keyword.
//
// Implementations MAY add custom format attributes. Save for agreement between
// parties, schema authors SHALL NOT expect a peer implementation to support this
// keyword and/or custom format attributes.
//
func validFormat(stack_id string, mem *JSONNode, schema *JSONNode, parent *JSONNode, errs *SchemaErrors) bool {
	if !validateFormat {
		return true
	}

	Trace.Println(stack_id, "validFormat")
	valid := false

	if schema.GetValueType() == V_STRING {
		schemaValue := schema.GetValue().(string)
		
		Trace.Println("  validate format", schemaValue)
		memValue := mem.GetValue().(string)
		if schemaValue == "date-time" {
			valid = validDateTime(memValue)
		} else if schemaValue == "email" {
			valid = validEmail(memValue)
		} else if schemaValue == "hostname" {
			valid = validHostname(memValue)
		} else if schemaValue == "ipv4" {
			valid = validIPV4(memValue)
		} else if schemaValue == "ipv6" {
			valid = validIPV6(memValue)
		} else if schemaValue == "uri" {
			valid = validURI(memValue)
		}
	}
	Trace.Println(stack_id, "validFormat", valid)
	return valid
}
// 7.3.1.  date-time
//
// 7.3.1.1.  Applicability
//
// This attribute applies to string instances.
//
// 7.3.1.2.  Validation
//
// A string instance is valid against this attribute if it is a valid date representation
// as defined by RFC 3339, section 5.6 [RFC3339].
//
func validDateTime(value string) bool {
	return regexDateTime.MatchString(value)
}

func validEmail(value string) bool {
	return regexEmail.MatchString(value)
}

func validHostname(value string) bool {
	valid := regexHostname.MatchString(value)

	if !valid {
		return false
	}

	parts := strings.Split(value, ".")
	for comp := range parts {
		if len(parts[comp]) > 63 {
			return false
		}
	}

	parts = strings.Split(value, "-")
	for comp := range parts {
		if len(parts[comp]) > 63 {
			return false
		}
	}

	return len(value) < 254
}

func validIPV4(value string) bool {
	ip := net.ParseIP(value)

	if ip == nil {
		return false
	}

	ipStr := ip.String()
	index := strings.Index(ipStr, ".")
	if index > -1 {
		return true
	}

	return false
}

func validIPV6(value string) bool {
	ip := net.ParseIP(value)

	if ip == nil {
		return false
	}

	ipStr := ip.String()
	index := strings.Index(ipStr, ":")
	if index > -1 {
		return true
	}

	return false
}

func validURI(value string) bool {
	_, err := url.ParseRequestURI(value)
	if err == nil {
		return true
	}

	return false
}
