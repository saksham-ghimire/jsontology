package jsontology

import (
	"reflect"
)

type Operator string

const (
	equals      Operator = "eq"
	notEquals   Operator = "neq"
	greaterThan Operator = "gt"
	lessThan    Operator = "lt"
)

var operatorMapping map[Operator]func(ruleValue interface{}, fieldValue interface{}) bool = map[Operator]func(ruleValue interface{}, fieldValue interface{}) bool{
	equals:      isEquals,
	notEquals:   isNotEquals,
	greaterThan: isGreaterThan,
	lessThan:    isLessThan,
}

func RegisterNewOperator(op Operator, opFunc func(ruleValue interface{}, fieldValue interface{}) bool) {
	operatorMapping[op] = opFunc
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
		return existsInArray(ruleParam, eventParam)
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
		for _, each_element := range eventParam {
			if isGreaterThan(ruleParam, each_element) {
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
		for _, each_element := range eventParam {
			if isLessThan(ruleParam, each_element) {
				return true
			}
		}

	}
	return false
}

func existsInArray(value interface{}, array []interface{}) bool {
	for _, element := range array {
		if isEquals(value, element) {
			return true
		}
	}
	return false
}
