package jsontology

import (
	"encoding/json"
	"testing"
)

func TestRule(t *testing.T) {

	table := []struct {
		name      string
		condition string
		data      string
		expected  bool
	}{
		{
			name:      "simple nested logic",
			condition: `[{"a.a.$eq":3,"c.$eq":"1","a.$nested":{"a.$eq":1}}]`,
			data:      `{"a":[{"a":1},{"a":2}],"c":"1"}`,
			expected:  false,
		},
		{
			name:      "array or condition",
			condition: `[{"a.a.$eq":1,"c.$eq":"1"}]`,
			data:      `{"a":[{"a":1},{"a":2}],"c":"1"}`,
			expected:  true,
		},
		{
			name:      "complex nested logic",
			condition: `[{"a.a.$eq":1,"c.$eq":"1","a.$nested":{"a.$eq":2,"c.$nested":{"g.$eq":"h","a.$eq":3}}}]`,
			data:      `{"a":[{"a":1}, {"a":2, "c": [{"a":2, "e":"f"}, {"g":"h", "a":3}]}], "c" : "1"}`,
			expected:  true,
		},

		{
			name:      "simple or logic",
			condition: `[{"a.a.$eq":1}, {"g.$eq":"h"}]`,
			data:      `{"a":[{"a":1},{"a":2}],"c":"1"}`,
			expected:  true,
		},

		{
			name:      "simple and logic",
			condition: `[{"a.a.$eq":1 , "g.$eq":"h"}]`,
			data:      `{"a":[{"a":1},{"a":2}],"g":"h"}`,
			expected:  true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			data := make(map[string]interface{})
			err := json.Unmarshal([]byte(tt.data), &data)
			if err != nil {
				t.Fatal("Invalid data", err)
			}
			conditions := []map[string]interface{}{}
			err = json.Unmarshal([]byte(tt.condition), &conditions)
			if err != nil {
				t.Fatal("Invalid JSON for rule", err)
			}
			r := NewRule(conditions, map[string]interface{}{}, &LogEventHandler{})
			if got := r.IsMatch(data); got != tt.expected {
				t.Errorf("IsMatch() = %v, want %v", got, tt.expected)
			}
		})
	}
}
