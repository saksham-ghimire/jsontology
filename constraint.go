package jsontology

type constraint interface {
	Evaluate(data map[string]interface{}) bool
}

const (
	keyNotFound string = "keyNotFound"
)

type criteria struct {
	field    string
	operator Operator
	value    interface{}
}

type nestedCriteria struct {
	path       string
	conditions []constraint
}

func (c criteria) Evaluate(data map[string]interface{}) bool {

	if value, ok := data[c.field]; ok {
		return operatorFuncMapping[c.operator](c.value, value)
	}
	return operatorFuncMapping[c.operator](c.value, keyNotFound)
}

func (c nestedCriteria) Evaluate(data map[string]interface{}) bool {

	// check if given path is array
	arrayData, ok := data[c.path].([]interface{})
	if ok {
		for _, eachData := range arrayData {
			if formattedData, ok := eachData.(map[string]interface{}); ok {
				andMatches := []bool{}
				for _, e := range c.conditions {
					andMatches = append(andMatches, e.Evaluate(formattedData))
				}
				if allMatch(andMatches) {
					return true
				}
			}
		}
	}
	return false
}
