package jsontology

import "testing"

func TestIsEquals(t *testing.T) {
	tests := []struct {
		ruleParam  interface{}
		eventParam interface{}
		expected   bool
		name       string
	}{
		{true, true, true, "Bool Equality"},
		{true, false, false, "Bool Inequality"},
		{42, 42, true, "Int Equality"},
		{42, 24, false, "Int Inequality"},
		{3.14, 3.14, true, "Float Equality"},
		{3.14, 2.71, false, "Float Inequality"},
		{"hello", "hello", true, "String Equality"},
		{"hello", "world", false, "String Inequality"},
		{[]int{1, 2, 3}, []int{1, 2, 3}, true, "Slice Equality"},
		{[]int{1, 2, 3}, []int{1, 2, 4}, false, "Slice Inequality"},
		{map[string]int{"a": 1}, map[string]int{"a": 1}, true, "Map Equality"},
		{map[string]int{"a": 1}, map[string]int{"b": 1}, false, "Map Key Difference"},
		{map[string]int{"a": 1}, map[string]int{"a": 2}, false, "Map Value Difference"},
		{42, []interface{}{1, 2, 42}, true, "Fallback Array Match"},
		{99, []interface{}{1, 2, 42}, false, "Fallback Array NoMatch"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isEquals(test.ruleParam, test.eventParam)
			if result != test.expected {
				t.Errorf("isEquals(%v, %v) = %v; want %v", test.ruleParam, test.eventParam, result, test.expected)
			}
		})
	}
}

func TestIsNotEquals(t *testing.T) {
	tests := []struct {
		ruleParam  interface{}
		eventParam interface{}
		expected   bool
		name       string
	}{
		{true, true, false, "Not Equals Bool Equality"},
		{true, false, true, "Not Equals Bool Inequality"},
		{42, 42, false, "Not Equals Int Equality"},
		{42, 24, true, "Not Equals Int Inequality"},
		{3.14, 3.14, false, "Not Equals Float Equality"},
		{3.14, 2.71, true, "No tEquals Float Inequality"},
		{"hello", "hello", false, "Not Equals String Equality"},
		{"hello", "world", true, "Not Equals String Inequality"},
		{[]int{1, 2, 3}, []int{1, 2, 3}, false, "Not Equals Slice Equality"},
		{[]int{1, 2, 3}, []int{1, 2, 4}, true, "Not Equals Slice Inequality"},
		{map[string]int{"a": 1}, map[string]int{"a": 1}, false, "Not Equals Map Equality"},
		{map[string]int{"a": 1}, map[string]int{"b": 1}, true, "Not Equals Map Key Difference"},
		{map[string]int{"a": 1}, map[string]int{"a": 2}, true, "Not Equals Map Value Difference"},
		{42, []interface{}{1, 2, 42}, false, "Not Equals Fallback ArrayMatch"},
		{99, []interface{}{1, 2, 42}, true, "Not Equals FallbackArray NoMatch"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isNotEquals(test.ruleParam, test.eventParam)
			if result != test.expected {
				t.Errorf("isNotEquals(%v, %v) = %v; want %v", test.ruleParam, test.eventParam, result, test.expected)
			}
		})
	}
}

func TestIsGreaterThan(t *testing.T) {
	tests := []struct {
		ruleParam  interface{}
		eventParam interface{}
		expected   bool
		name       string
	}{
		{10, 20, true, "Greater Than Int"},
		{20, 10, false, "Greater Than Int Inequality"},
		{10.5, 20.5, true, "Greater Than Float"},
		{20.5, 10.5, false, "Greater Than Float Inequality"},
		{42, []interface{}{10, 20, 30}, false, "Greater Than Fallback Array No Match"},
		{42, []interface{}{50, 60}, true, "Greater Than Fallback Array Match"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isGreaterThan(test.ruleParam, test.eventParam)
			if result != test.expected {
				t.Errorf("isGreaterThan(%v, %v) = %v; want %v", test.ruleParam, test.eventParam, result, test.expected)
			}
		})
	}
}

func TestIsLessThan(t *testing.T) {
	tests := []struct {
		ruleParam  interface{}
		eventParam interface{}
		expected   bool
		name       string
	}{
		{20, 10, true, "Less Than Int"},
		{10, 20, false, "Less Than Int Inequality"},
		{20.5, 10.5, true, "Less Than Float"},
		{10.5, 20.5, false, "Less Than Float Inequality"},
		{42, []interface{}{50, 60, 70}, false, "Less Than Fallback Array No Match"},
		{42, []interface{}{10, 20}, true, "Less Than Fallback Array Match"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isLessThan(test.ruleParam, test.eventParam)
			if result != test.expected {
				t.Errorf("isLessThan(%v, %v) = %v; want %v", test.ruleParam, test.eventParam, result, test.expected)
			}
		})
	}
}
