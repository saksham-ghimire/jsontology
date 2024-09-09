package jsontology

import (
	"log"
	"time"
)

type eventHandler interface {
	call(eventJson, extraParams map[string]interface{})
}

type LogEventHandler struct {
	Logger *log.Logger
}

type CountEventHandler struct {
	currentCount int
	count        int
	handler      eventHandler
}

type GroupByEventHandler struct {
	currentState map[interface{}]int
	count        int
	groupBy      string
	handler      eventHandler
}

type TimeBasedCountEventHandler struct {
	eventTimings []int
	count        int
	timeLimit    int
	handler      eventHandler
}

func NewCountEventHandler(count int, handler eventHandler) *CountEventHandler {
	return &CountEventHandler{
		currentCount: 0,
		count:        count,
		handler:      handler,
	}
}

func NewGroupByEventHandler(count int, groupBy string, handler eventHandler) *GroupByEventHandler {
	return &GroupByEventHandler{
		currentState: make(map[interface{}]int),
		count:        count,
		groupBy:      groupBy,
		handler:      handler,
	}
}
func NewTimeBasedCountEventHandler(count int, timeLimit int, handler eventHandler) *TimeBasedCountEventHandler {
	return &TimeBasedCountEventHandler{
		eventTimings: []int{},
		count:        count,
		timeLimit:    timeLimit,
	}
}

func (l *LogEventHandler) call(eventJson, extraParams map[string]interface{}) {
	l.Logger.Printf("Event Matched \n Event : %v \n Rule Params : %v", eventJson, extraParams)
}

func (c *CountEventHandler) call(eventJson, extraParams map[string]interface{}) {
	c.currentCount += 1
	if c.currentCount == c.count {
		c.handler.call(eventJson, extraParams)
		// Reset the previous counter once handler is called
		c.currentCount = 0
	}
}

func (c *GroupByEventHandler) call(eventJson, extraParams map[string]interface{}) {
	c.currentState[eventJson[c.groupBy]] += 1
	for key, value := range c.currentState {
		if value == c.count {
			c.handler.call(eventJson, extraParams)
			// Reset the previous counter once handler is called
			delete(c.currentState, key)
		}
	}
}

func (c *TimeBasedCountEventHandler) call(eventJson, extraParams map[string]interface{}) {
	// clean up expired timing
	var currentTimestamp int64 = time.Now().Unix()
	c.eventTimings = filter(c.eventTimings, func(x int) bool { return (x + c.timeLimit) >= int(currentTimestamp) })
	// append new timings
	c.eventTimings = append(c.eventTimings, int(currentTimestamp))

	if len(c.eventTimings) == c.count {
		c.handler.call(eventJson, extraParams)

		// Reset the previous counter once handler is called
		c.eventTimings = []int{}
	}

}
