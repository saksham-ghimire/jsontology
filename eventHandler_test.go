package jsontology

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/mock"
)

type eventHandlerMock struct {
	mock.Mock
}

func (m *eventHandlerMock) call(eventJson, extraParams map[string]interface{}) {
	m.Called()
}

func TestCountEventHandler(t *testing.T) {
	table := []struct {
		name       string
		jsonEvents []string
		condition  string
		handler    eventHandler
		matches    int
	}{
		{ // Generate 1 alert if {"c":1} event is detected in stream 2 times
			name:       "simple count event handler",
			jsonEvents: []string{`{"a":[{"a":1}, {"a":2}], "c" : 1}`, `{"c": 1}`},
			condition:  `[{"c.$eq":1}]`,
			handler: &CountEventHandler{
				count:   2,
				handler: &eventHandlerMock{},
			}, matches: 1,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.handler.(*CountEventHandler).handler.(*eventHandlerMock)
			h.On("call").Times(tt.matches)
			h.On("call").Return(nil)
			ruleCondition := []map[string]interface{}{}
			err := json.Unmarshal([]byte(tt.condition), &ruleCondition)
			if err != nil {
				t.Fatal("Invalid data", err)
			}
			rule := NewRule(ruleCondition, map[string]interface{}{}, tt.handler)
			for _, eachJson := range tt.jsonEvents {
				jsonData := map[string]interface{}{}
				err = json.Unmarshal([]byte(eachJson), &jsonData)
				if err != nil {
					t.Fatal("Invalid data", err)
				}
				rule.Send(jsonData)
			}
			h.AssertExpectations(t)
		})
	}
}

func TestGroupByEventHandler(t *testing.T) {
	table := []struct {
		name       string
		jsonEvents []string
		condition  string
		handler    eventHandler
		matches    int
	}{
		{ // Generate 1 alert if same value for key [name] is in the stream for 2 times
			name:       "simple group by event handler",
			jsonEvents: []string{`{"name": "someName"}`, `{"name":"someName"}`, `{"name":"someName"}`, `{"name":"someName"}`},
			condition:  `[{"name.$eq": "someName"}]`,
			handler: &GroupByEventHandler{
				currentState: map[interface{}]int{},
				count:        2,
				groupBy:      "name",
				handler:      &eventHandlerMock{},
			},
			matches: 2,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			h := tt.handler.(*GroupByEventHandler).handler.(*eventHandlerMock)
			h.On("call").Times(tt.matches)
			h.On("call").Return(nil)

			ruleCondition := []map[string]interface{}{}
			err := json.Unmarshal([]byte(tt.condition), &ruleCondition)
			if err != nil {
				t.Fatal("Invalid data", err)
			}

			rule := NewRule(ruleCondition, map[string]interface{}{}, tt.handler)

			for _, eachJson := range tt.jsonEvents {
				jsonData := map[string]interface{}{}
				err = json.Unmarshal([]byte(eachJson), &jsonData)
				if err != nil {
					t.Fatal("Invalid data", err)
				}
				rule.Send(jsonData)
			}

			h.AssertExpectations(t)
		})
	}
}
