
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

