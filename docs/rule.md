
## Understanding Rule Layout

* Rule Structure : A Rule defines a set of conditions that must be evaluated to determine if an event handler should be triggered. The structure of a Rule is designed to support complex conditional logic and event handling. Here's a breakdown:

```go
type Rule struct {
	condition  [][]constraint // list of constraint that needs to be matched for rule to call eventHandler
	onMatch    eventHandler // eventHandler to call when rule matches
	extraParam map[string]interface{} // any information about rule that gets passed on to eventHandler
}
```

### Condition Structure
Conditions are specified using a nested approach to handle logical operations like AND and OR:

* __OR Conditions__: Conditions that are placed in separate objects within the list are evaluated with an "OR" logic. This means that if any of these conditions are true, the rule is considered to match.

* __AND Conditions__: Conditions that are grouped together within a single object in the list are evaluated with an "AND" logic. This means that all conditions within this object must be true for the rule to match.

* __Operators__: Operators are used to specify the type of comparison to be performed. In the key of a condition, the part after .$ indicates the operation. For example:
"a.$eq": 1 specifies that the condition is checking if the value of a is equal to 1.
"b.$gt": 10 specifies that the condition is checking if the value of b is greater than 10.



Consider the following object

```
{"a": "b" , "c" : "d", "e": [1,2,3]}
```

* __OR Condition__

For the rule (a == b) or (c == d), you want to match if either a is "b" or c is "d". To write this rule, you need to create separate objects for each condition within the same list

```
[{"a.$eq" : "b"}, {"c.$eq" : "d"}]
```
* __AND Condition__

For the rule (a == b) and (c == d), both conditions must be true simultaneously. To write this rule, put all conditions into one object
```
[{"a.$eq" : "b", "c.$eq" : "d"}]
```

*NOTE : Please refer to test cases for more advance example.*

### Rule Building Process

To simplify a complex logical expression, follow these steps:

1. **Original Expression:** `(((a or b) and c) or d)`

2. **Replace Operators:**
   - Replace `or` with `+`
   - Replace `and` with `*`

3. **Converted Expression:**
   - `(((a + b) * c) + d)`

4. **Expand and Simplify:**
   - Distribute multiplication over addition:
     - `((a * c) + (b * c) + d)`
   - **Meaning:** `a and c` or `b and c` or `d`
   - **Expressed:** `[{"a.$":"", "c.$":""},{"b.$":"", "c.$":""},{"d,$":""}]`



#### Examples

* **Rule example** 

```go

func main() {
   condition := `[{"a.a.$eq":1 , "g.$eq":"h"}]`
   data := `{"a":[{"a":1},{"a":2}],"g":"h"}`
   parsedData := make(map[string]interface{})
   json.Unmarshal([]byte(data), &parsedData)
   r, _ := jsontology.NewRule(strings.NewReader(condition), map[string]interface{}{}, &jsontology.LogEventHandler{
      Logger: log.Default(),
   })
   // to verify is rule matches with data 
   // Note that when JSON is unmarshalled 
	// into a map in Go, numeric fields like `{"a": 4}` may be represented as `float64` instead of `int`. 
	// Be mindful of this when passing data parsed from JSON to ensure compatibility. (Alternatively use Send)
   fmt.Println(r.IsMatch(parsedData))
   // to trigger the event chain in case of match
   r.Send(strings.NewReader(data))
}

```

* **Dealing with nested Json**

```go
func main() {
	condition := `[{"a.$nested":{"a.$eq": 1, "b.$eq": 3}}]`
	data := `{"a":[{"a":1},{"a":1, "b": 3}],"g":"h"}`
	parsedData := make(map[string]interface{})
	json.Unmarshal([]byte(data), &parsedData)
	r, _ := jsontology.NewRule(strings.NewReader(condition), map[string]interface{}{}, &jsontology.LogEventHandler{
		Logger: log.Default(),
	})
	// to verify is rule matches with data
	fmt.Println(r.IsMatch(parsedData))
	// to trigger the event chain in case of match
	r.Send(strings.NewReader(data))
}
```

* **Parse event handler chain from json**

```go
func main() {
	eventHandlerChain := `{"handler":{"type":"CountEventHandler","params":{"count":3,"handler":{"type":"LogEventHandler","params":{}}}}}`
	conditions := `[{"name.$eq": "someName"}]`
	jsonEvents := []string{`{"name": "someName"}`, `{"name": "someName"}`, `{"name": "someName"}`, `{"name": "someName"}`}

	eventHandler, err := jsontology.GetEventHandlerChain(strings.NewReader(eventHandlerChain))
	if err != nil {
		log.Fatalf("Unable to get rule handler error %v", err)
	}

	r, err := jsontology.NewRule(strings.NewReader(conditions), map[string]interface{}{"rule_id": 1}, eventHandler)

	if err != nil {
		log.Fatalf("Unable to get rule received error %v", err)
	}
	for _, eachEvent := range jsonEvents {
		r.Send(strings.NewReader(eachEvent))
	}
}
```

**Output**

```
2024/09/15 20:55:29 Event Matched 
 Event : map[name:someName] 
 Rule Params : map[rule_id:1]
```

_Note : Pre configured eventhandlers are CountEventHandler, GroupByEventHandler, TimeBasedCountEventHandler, LogEventHandler ..._
