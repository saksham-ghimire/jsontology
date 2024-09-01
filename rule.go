package jsontology

import (
	"strings"
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
func NewRule(conditions []map[string]interface{}, params map[string]interface{}, onMatch eventHandler) *Rule {

	processedConditions := parseJsonToContext(conditions)
	return &Rule{
		condition:  processedConditions,
		onMatch:    onMatch,
		extraParam: params,
	}
}

func parseJsonToContext(data []map[string]interface{}) [][]constraint {

	var returnContext [][]constraint
	for _, eachData := range data {
		// prepare internal constraint var which will hold internal array in 2d array on returnContext
		var internalContext []constraint
		for key, value := range eachData {
			splittedString := strings.Split(key, ".$") // separator
			field, operator := splittedString[0], splittedString[1]
			if operator == "nested" {
				// nested operator should always have map[string]interface{} as comparison type
				formattedValue, ok := value.(map[string]interface{})
				if !ok {
					panic(ok)
				}
				sformattedValue := []map[string]interface{}{formattedValue}
				internalContext = append(internalContext, nestedCriteria{
					path:       field,
					conditions: parseJsonToContext(sformattedValue)[0],
				})
			} else {
				internalContext = append(internalContext, criteria{
					field:    field,
					operator: Operator(operator),
					value:    value,
				})
			}
		}
		returnContext = append(returnContext, internalContext)
	}
	return returnContext
}

func (r *Rule) IsMatch(data map[string]interface{}) bool {

	var normalizedJson = transformJSON(data, "")
	normalizedJson = concatMaps(data, normalizedJson)
	or_matches := []bool{}
	for _, e := range r.condition {

		and_matches := []bool{}
		for _, eachContext := range e {
			and_matches = append(and_matches, eachContext.Evaluate(normalizedJson))
		}
		or_matches = append(or_matches, all_match(and_matches))
	}
	return any_match(or_matches)

}

func (r *Rule) Send(data map[string]interface{}) {
	if r.IsMatch(data) {
		r.onMatch.call(data, r.extraParam)
	}
}
