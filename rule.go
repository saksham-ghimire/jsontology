package jsontology

import (
	"encoding/json"
	"io"
)

type Rule struct {
	condition  [][]constraint
	onMatch    eventHandler
	extraParam map[string]interface{}
}

// NewRule creates a new rule with given conditions and event handler.
//
// conditions: A slice of maps, where each map represents a condition.
// Each condition map contains field-operator-value pairs.
// For nested conditions, use ".$nested" as the operator and provide a map as the value.
//
// params: A map of parameters that can be used in the conditions.
//
// onMatch: An event handler function that will be called when the rule matches.
//
// Returns a pointer to a new Rule instance.
func NewRule(conditions io.Reader, params map[string]interface{}, onMatch eventHandler) (*Rule, error) {

	var parsedConditions []map[string]interface{}

	conditionBytes, err := io.ReadAll(conditions)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(conditionBytes, &parsedConditions); err != nil {
		return nil, err
	}

	processedConditions, err := parseJsonToContext(parsedConditions)
	if err != nil {
		return nil, err
	}
	return &Rule{
		condition:  processedConditions,
		onMatch:    onMatch,
		extraParam: params,
	}, nil
}

// IsMatch checks if the provided data meets the rule's conditions.
//
// The `data` is first normalized and merged with itself for evaluation. Note that when JSON is unmarshalled
// into a map in Go, numeric fields like `{"a": 4}` may be represented as `float64` instead of `int`.
// Be mindful of this when passing data parsed from JSON to ensure compatibility.
//
// Parameters:
// - data: A map representing the input data to be checked.
//
// Returns:
// - bool: true if the data matches the conditions, false otherwise.
func (r *Rule) IsMatch(data map[string]interface{}) bool {

	var normalizedJson = transformJSON(data, "")
	normalizedJson = concatMaps(data, normalizedJson)
	orMatches := []bool{}
	for _, e := range r.condition {

		andMatches := []bool{}
		for _, eachContext := range e {
			andMatches = append(andMatches, eachContext.Evaluate(normalizedJson))
		}
		orMatches = append(orMatches, allMatch(andMatches))
	}
	return anyMatch(orMatches)

}

func (r *Rule) Send(data io.Reader) error {
	var parsedData map[string]interface{}

	dataBytes, err := io.ReadAll(data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dataBytes, &parsedData); err != nil {
		return err
	}
	if isMatch := r.IsMatch(parsedData); isMatch {
		r.onMatch.call(parsedData, r.extraParam)
	}
	return nil
}
