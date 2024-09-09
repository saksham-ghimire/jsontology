package jsontology

import (
	"net"
	"reflect"
	"regexp"
	"strings"
)

type Operator string
type OperatorFunc func(ruleValue interface{}, fieldValue interface{}) bool
type OperatorTypeHandlerFunc func(value interface{}) (interface{}, error)

const (
	equals        Operator = "eq"
	notEquals     Operator = "neq"
	greaterThan   Operator = "gt"
	lessThan      Operator = "lt"
	startsWith    Operator = "sw"
	endsWith      Operator = "ew"
	regexMatch    Operator = "rgx"
	notRegexMatch Operator = "nrgx"
	nested        Operator = "nested"
	ipInRange     Operator = "ipInRange"
)

var operatorFuncMapping map[Operator]OperatorFunc = map[Operator]OperatorFunc{
	equals:        isEquals,
	notEquals:     isNotEquals,
	greaterThan:   isGreaterThan,
	lessThan:      isLessThan,
	startsWith:    isStartingWith,
	endsWith:      isEndingWith,
	regexMatch:    isRegexMatch,
	notRegexMatch: isNotRegexMatch,
	ipInRange:     isIPInRange,
}

var operatorTypeHandlerMapping map[Operator]OperatorTypeHandlerFunc = map[Operator]OperatorTypeHandlerFunc{
	greaterThan:   isNumber,
	lessThan:      isNumber,
	regexMatch:    asRegexExpression,
	notRegexMatch: asRegexExpression,
	ipInRange:     asIpNet,
}

// RegisterNewOperator registers a new operator with the system.
//
// The `op` parameter specifies the operator to be registered.
// The `opFunc` parameter is the function that will be called when the operator is executed.
// The `opTypeHandler` parameter is an optional function that can be used to pre-process data per operator's type.
func RegisterNewOperator(op Operator, opFunc OperatorFunc, opTypeHandler OperatorTypeHandlerFunc) {
	operatorFuncMapping[op] = opFunc
	if opTypeHandler != nil {
		operatorTypeHandlerMapping[op] = opTypeHandler
	}
}

func isEquals(ruleParam, eventParam interface{}) bool {

	v1 := reflect.ValueOf(ruleParam)
	v2 := reflect.ValueOf(eventParam)
	if v1.Kind() == v2.Kind() {
		switch v1.Kind() {

		case reflect.Bool:
			return v1.Bool() == v2.Bool()

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v1.Int() == v2.Int()

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return v1.Uint() == v2.Uint()

		case reflect.Float32, reflect.Float64:
			return v1.Float() == v2.Float()

		case reflect.Complex64, reflect.Complex128:
			return v1.Complex() == v2.Complex()

		case reflect.String:
			return v1.String() == v2.String()

		case reflect.Slice:
			for i := 0; i < v1.Len(); i++ {
				if !isEquals(v1.Index(i).Interface(), v2.Index(i).Interface()) {
					return false
				}
			}
			return true

		case reflect.Map:
			if v1.Len() != v2.Len() {
				return false
			}
			for _, key := range v1.MapKeys() {
				v1Value := v1.MapIndex(key)
				v2Value := v2.MapIndex(key)
				if !v2Value.IsValid() || !isEquals(v1Value.Interface(), v2Value.Interface()) {
					return false
				}
			}
			return true
		}
	}

	// if type is not equal only one fallback is supported
	switch eventParam := eventParam.(type) {
	case []interface{}:
		return isInArray(ruleParam, eventParam)
	}

	return false
}

func isNotEquals(ruleParam, eventParam interface{}) bool {
	return !isEquals(ruleParam, eventParam)
}

func isGreaterThan(ruleParam, eventParam interface{}) bool {

	v1 := reflect.ValueOf(ruleParam)
	v2 := reflect.ValueOf(eventParam)

	if v1.Kind() == v2.Kind() {
		switch v1.Kind() {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v2.Int() > v1.Int()

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return v2.Uint() > v1.Uint()

		case reflect.Float32, reflect.Float64:
			return v2.Float() > v1.Float()

		}
	}
	switch eventParam := eventParam.(type) {
	case []interface{}:
		for _, eachElement := range eventParam {
			if isGreaterThan(ruleParam, eachElement) {
				return true
			}
		}

	}
	return false
}

func isLessThan(ruleParam, eventParam interface{}) bool {
	v1 := reflect.ValueOf(ruleParam)
	v2 := reflect.ValueOf(eventParam)

	if v1.Kind() == v2.Kind() {
		switch v1.Kind() {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v2.Int() < v1.Int()

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return v2.Uint() < v1.Uint()

		case reflect.Float32, reflect.Float64:
			return v2.Float() < v1.Float()

		}
	}
	switch eventParam := eventParam.(type) {
	case []interface{}:
		for _, eachElement := range eventParam {
			if isLessThan(ruleParam, eachElement) {
				return true
			}
		}

	}
	return false
}

func isInArray(value interface{}, array []interface{}) bool {
	for _, element := range array {
		if isEquals(value, element) {
			return true
		}
	}
	return false
}

func isStartingWith(ruleParam, eventParam interface{}) bool {
	v1 := reflect.ValueOf(ruleParam)
	v2 := reflect.ValueOf(eventParam)
	if v1.Kind() == v2.Kind() {
		switch v1.Kind() {
		case reflect.String:
			return strings.HasPrefix(v2.String(), v1.String())
		}
	}
	return false
}

func isEndingWith(ruleParam, eventParam interface{}) bool {
	v1 := reflect.ValueOf(ruleParam)
	v2 := reflect.ValueOf(eventParam)
	if v1.Kind() == v2.Kind() {
		switch v1.Kind() {
		case reflect.String:
			return strings.HasSuffix(v2.String(), v1.String())
		}
	}
	return false
}

func isRegexMatch(ruleParam, eventParam interface{}) bool {
	v1 := reflect.ValueOf(ruleParam)
	v2 := reflect.ValueOf(eventParam)
	if v2.Kind() == reflect.String {
		re := v1.Interface().(*regexp.Regexp)
		text := v2.String()
		return re.MatchString(text)
	}
	return false
}

func isNotRegexMatch(ruleParam, eventParam interface{}) bool {
	return !isRegexMatch(ruleParam, eventParam)
}

func isIPInRange(ruleParam, eventParam interface{}) bool {

	v2 := reflect.ValueOf(eventParam)
	v1 := reflect.ValueOf(ruleParam)

	if v2.Kind() == reflect.String {
		ip := net.ParseIP(v2.String())
		if ip == nil {
			return false
		}
		subnetStr := v1.Interface().(*net.IPNet)
		return subnetStr.Contains(ip)
	}
	return false
}
