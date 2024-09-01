package jsontology

import "strings"

func all_match(slice []bool) bool {
	for _, v := range slice {
		if !v {
			return false
		}
	}
	return true
}

func any_match(slice []bool) bool {
	for _, v := range slice {
		if v {
			return true
		}
	}
	return false
}

// concatMaps merges two maps into a new map. If a key exists in both maps,
// the value from the second map is appended to the value in the first map.
// If the value in the first map is not a slice, it is converted to a slice
// before appending the value from the second map.
//
// Parameters:
// - m1: The first map to merge.
// - m2: The second map to merge.
//
// Returns:
// - A new map containing the merged key-value pairs.
func concatMaps(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {

		if value, ok := result[k]; ok {

			switch vv := value.(type) {
			case []interface{}:
				result[k] = append(vv, v)
			default:
				result[k] = []interface{}{value, v}
			}

		} else {
			result[k] = v
		}
	}
	return result
}

// transformJSON is a recursive function that transforms a nested JSON-like data structure into a flat map.
// It handles bool, int, float64, string, map[string]interface{}, and []interface{} types.
//
// Parameters:
// - data: The input data to be transformed. It can be of any type mentioned above.
// - currentPath: The current path of the data in the nested structure. It is used to construct the keys in the flat map.
//
// Returns:
// - A map[string]interface{} representing the flat map of the transformed data.
//
// Note:
// - If the input data is a map[string]interface{}, it recursively calls itself for each value with the updated currentPath.
// - If the input data is a []interface{}, it iterates over each value and handles it accordingly.
// - If the input data is a bool, int, float64, or string, it checks if the currentPath contains a dot and if it does, it adds the currentPath and the value to the result map.
func transformJSON(data interface{}, currentPath string) map[string]interface{} {
	var result = make(map[string]interface{})
	switch v := data.(type) {
	case bool, int, float64, string:
		if strings.Contains(currentPath, ".") {
			result[currentPath] = v
		}
	case map[string]interface{}:
		for key, value := range v {
			if currentPath != "" {
				key = currentPath + "." + key
			}
			result = concatMaps(result, transformJSON(value, key))
		}
	case []interface{}:
		for _, each_value := range v {
			switch each_value.(type) {
			case map[string]interface{}:
				result = concatMaps(result, transformJSON(each_value, currentPath))
			default:
				if _, ok := result[currentPath]; ok {
					d := result[currentPath].([]interface{})
					result[currentPath] = append(d, each_value)
				} else {
					if strings.Contains(currentPath, ".") {
						result[currentPath] = []interface{}{each_value}
					}
				}
			}
		}
	}
	return result
}

func filter[T any](iterable []T, function func(T) bool) (ret []T) {

	for _, s := range iterable {
		if function(s) {
			ret = append(ret, s)
		}
	}
	return
}
