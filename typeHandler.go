package jsontology

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
)

func isNumber(value interface{}) (interface{}, error) {
	v1 := reflect.ValueOf(value)

	switch v1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32,
		reflect.Float64:
		return value, nil
	}
	return nil, fmt.Errorf("validation failed, %v is not a number", value)
}

func asRegexExpression(value interface{}) (interface{}, error) {

	regStr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("validation failed, %v is not a string", value)
	}
	re, err := regexp.Compile(regStr)
	if err != nil {
		return nil, err
	}
	return re, nil
}

func asIpNet(value interface{}) (interface{}, error) {
	str, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("validation failed, %v is not a string", value)
	}

	ipv4Regex := regexp.MustCompile(`^((25[0-5]|2[0-4]\d|1\d{2}|[1-9]\d|[1-9])\.){3}(25[0-5]|2[0-4]\d|1\d{2}|[1-9]\d|[1-9])$`)
	ipv6Regex := regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
	
	if !(ipv4Regex.MatchString(str) || ipv6Regex.MatchString(str)) {
		return false, fmt.Errorf("validation failed, invalid IP address: %s", str)
	}

	// If networkStr is an IP with a subnet (CIDR notation)
	_, subnetStr, err := net.ParseCIDR(str)
	if err != nil {
		return false, fmt.Errorf("validation failed, invalid CIDR subnet: %s", str)
	}
	return subnetStr, nil 
}
