package jsontology

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

var eventHandlerParsingMap map[string]EventHandlerParsingFunctions

type EventHandlerParsingFunctions func(params map[string]interface{}) (eventHandler, error)

func init() {
	eventHandlerParsingMap = map[string]EventHandlerParsingFunctions{
		"CountEventHandler":          parseCountEventHandler,
		"GroupByEventHandler":        parseGroupByEventHandler,
		"TimeBasedCountEventHandler": parseTimeBasedCountEventHandler,
		"LogEventHandler":            parseLogEventHandler,
	}
}

func buildEventHandlerChain(handlerChain map[string]interface{}) (eventHandler, error) {

	if nestedHandlerData, hasNestedHandler := handlerChain["handler"]; hasNestedHandler {
		nestedHandlerMap := nestedHandlerData.(map[string]interface{})
		handlerType := nestedHandlerMap["type"].(string)
		handlerParams := nestedHandlerMap["params"].(map[string]interface{})
		return eventHandlerParsingMap[handlerType](handlerParams)
	}

	handlerType := handlerChain["type"].(string)
	handlerParams := handlerChain["params"].(map[string]interface{})
	return eventHandlerParsingMap[handlerType](handlerParams)
}

func GetEventHandlerChain(handlerChain io.Reader) (eventHandler, error) {

	var parsedHandlerChain map[string]interface{}

	handlerChainBytes, err := io.ReadAll(handlerChain)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(handlerChainBytes, &parsedHandlerChain); err != nil {
		return nil, err
	}
	return buildEventHandlerChain(parsedHandlerChain)

}

func RegisterEventHandlerParser(handlerKey string, eventHandlerParsingFunc EventHandlerParsingFunctions) {
	eventHandlerParsingMap[handlerKey] = eventHandlerParsingFunc
}

func parseCountEventHandler(params map[string]interface{}) (eventHandler, error) {
	// Validate "count" field
	count, ok := params["count"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid type for 'count': expected int, got %T", count)
	}
	// Validate "handler" field
	handlerParams, ok := params["handler"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type for 'handler': expected map[string]interface{}, got %T", handlerParams)

	}
	// Resolve handler chain
	resolvedHandler, err := buildEventHandlerChain(handlerParams)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve handler chain: %w", err)
	}

	return &CountEventHandler{
		currentCount: 0,
		count:        int(count),
		handler:      resolvedHandler,
	}, nil

}

func parseGroupByEventHandler(params map[string]interface{}) (eventHandler, error) {
	// Validate "count" field
	count, ok := params["count"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'count': expected int, got %T", params["count"])
	}

	// Validate "groupby" field
	groupBy, ok := params["groupBy"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'groupby': expected string, got %T", params["groupBy"])
	}

	// Validate "handler" field
	handlerParams, ok := params["handler"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'handler': expected map[string]interface{}, got %T", params["handler"])
	}

	// Resolve handler chain
	resolvedHandler, err := buildEventHandlerChain(handlerParams)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve handler chain: %w", err)
	}

	return &GroupByEventHandler{
		currentState: make(map[interface{}]int),
		groupBy:      groupBy,
		count:        int(count),
		handler:      resolvedHandler,
	}, nil

}

func parseTimeBasedCountEventHandler(params map[string]interface{}) (eventHandler, error) {
	// Validate "count" field
	count, ok := params["count"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'count': expected int, got %T", params["count"])
	}

	// Validate "timeLimit" field
	timeLimit, ok := params["timeLimit"].(int)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'timeLimit': expected int, got %T", params["timeLimit"])
	}

	// Validate "handler" field
	handlerParams, ok := params["handler"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'handler': expected map[string]interface{}, got %T", params["handler"])
	}

	// Resolve handler chain
	resolvedHandler, err := buildEventHandlerChain(handlerParams)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve handler chain: %w", err)
	}

	return &TimeBasedCountEventHandler{
		eventTimings: []int{},
		timeLimit:    timeLimit,
		count:        int(count),
		handler:      resolvedHandler,
	}, nil

}

func parseLogEventHandler(params map[string]interface{}) (eventHandler, error) {
	return &LogEventHandler{
		Logger: log.Default(),
	}, nil
}
