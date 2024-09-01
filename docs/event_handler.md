
## Creating Your Own EventHandler to Handle Events

In Json Eventizer, custom event handlers allow for flexible event processing. Here's how you can define and implement your own event handlers.

### The EventHandler Interface

The `eventHandler` interface specifies the method that must be implemented:

```go
type eventHandler interface {
	call(eventJson, extraParams map[string]interface{})
}
```

So any struct adhering to the interface is a eventHandler.
Create event handlers to perform specific actions. For example, a CountEventHandler counts events and triggers another handler when a count is reached:

```go
type CountEventHandler struct {
	currentCount int
	count        int
	handler      eventHandler
}

func (c *CountEventHandler) call(eventJson, extraParams map[string]interface{}) {
	c.currentCount += 1
	if c.currentCount == c.count {
		c.handler.call(eventJson, extraParams)
		c.currentCount = 0
	}
}
```

### Chaining EventHandlers

Event handlers can be chained to create complex processing pipelines. Each handler can invoke another, allowing for flexible event handling.

__Example Usage__
```go
// Final handler
finalHandler := &MyFinalEventHandler{}

customHandler := &CustomHandler{
   handler : finalHandler
}
func (c *customHandler) call(){
   // Do anything here
   // eg, log.info("event detected sending to elasticsearch")
   c.handler.call()
}
// CountEventHandler triggers the customHandler after 5 events which does action defined in `call` method
// and based on implementation will trigger finalHandler 
countHandler := &CountEventHandler{
	count:   5,
	handler: customHandler,
}

rule := NewRule(ruleCondition, map[string]interface{}{}, countHandler)

```