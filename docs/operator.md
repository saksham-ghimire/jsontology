
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

### Validator for custom operator

When rule are parsed certain operator support parsing the values to different type altogether. e.g. 
for operator 'rgx' that checks whether or not string matches the regex, the rule value can be parsed into regex object prior to matching so
function `asRegexExpression` is used.

**Example**

```go


func isRegexMatch(ruleParam, eventParam interface{}) bool {
	v1 := reflect.ValueOf(ruleParam)
	v2 := reflect.ValueOf(eventParam)
	if v2.Kind() == reflect.String {
		re := v1.Interface().(*regexp.Regexp)
		text := v2.String()
		return re.MatchString(text)
	}
	return false
}

func asRegexExpression(value interface{}) (interface{}, error) {

	regStr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("validation failed, %v is not a string", value)
	}
	re, err := regexp.Compile(regStr)
	if err != nil {
		return nil, err
	}
	return re, nil
}

const (
   rgx Operator = "rgx"
)

RegisterNewOperator(rgx, isRegexMatch, asRegexExpression)

```

### Validator for custom operator

When rule are parsed certain operator support parsing the values to different type altogether. e.g. 
for operator 'rgx' that checks whether or not string matches the regex, the rule value can be parsed into regex object prior to matching so
function `asRegexExpression` is used.

**Example**

```go


func isRegexMatch(ruleParam, eventParam interface{}) bool {
	v1 := reflect.ValueOf(ruleParam)
	v2 := reflect.ValueOf(eventParam)
	if v2.Kind() == reflect.String {
		re := v1.Interface().(*regexp.Regexp)
		text := v2.String()
		return re.MatchString(text)
	}
	return false
}

func asRegexExpression(value interface{}) (interface{}, error) {

	regStr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("validation failed, %v is not a string", value)
	}
	re, err := regexp.Compile(regStr)
	if err != nil {
		return nil, err
	}
	return re, nil
}

const (
   rgx Operator = "rgx"
)

RegisterNewOperator(rgx, isRegexMatch, asRegexExpression)

```