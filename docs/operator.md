
## Creating your own operator
In addition to the built-in operators, you can extend the functionality of your system by creating and registering your own custom operators. This allows you to define new comparison rules that suit your specific needs.

```go
// Custom operator
func eqStartsWith(ruleValue interface{}, fieldValue interface{}) bool {
    ruleStr, ok1 := ruleValue.(string)
    fieldStr, ok2 := fieldValue.(string)
    if !ok1 || !ok2 {
        return false
    }
    return strings.HasPrefix(fieldStr, ruleStr)
}

const (
   startsWith Operator = "startswith"
)
RegisterNewOperator(startsWith, eqStartsWith)
// and now you can use "a.$startswith" on additional rule
```