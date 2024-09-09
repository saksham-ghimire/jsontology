package jsontology

import (
	"errors"
	"strings"
)

func parseJsonToContext(data []map[string]interface{}) ([][]constraint, error) {

	var returnContext [][]constraint
	for _, eachData := range data {
		// prepare internal constraint var which will hold internal array in 2d array on returnContext
		var internalContext []constraint
		for key, value := range eachData {
			splittedString := strings.Split(key, ".$") // separator
			field, operator := splittedString[0], Operator(splittedString[1])

			if operator == nested {
				formattedValue, ok := value.(map[string]interface{})
				if !ok {
					return nil, errors.New("parsing error, value for nested operator is not map[string]interface{}")
				}
				mapFormattedValue := []map[string]interface{}{formattedValue}
				nestedContext, err := parseJsonToContext(mapFormattedValue)
				if err != nil {
					return nil, err
				}
				internalContext = append(internalContext, nestedCriteria{
					path:       field,
					conditions: nestedContext[0],
				})
			} else {
				if transformer, ok := operatorTypeHandlerMapping[operator]; ok {
					transformedValue, err := transformer(value)
					if err != nil {
						return nil, err
					}
					value = transformedValue
				}
				internalContext = append(internalContext, criteria{
					field:    field,
					operator: operator,
					value:    value,
				})
			}
		}
		returnContext = append(returnContext, internalContext)
	}
	return returnContext, nil
}
