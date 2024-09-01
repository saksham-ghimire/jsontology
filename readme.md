# Jsontology

Jsontology (JSON+Ontology) is a minimalistic JSON matching engine that supports nested matching and common aggregation handlers (e.g., count handler, groupby handler). It generates events internally and allows users to configure their own event handlers.

To get started install package using
```
go get github.com/saksham-ghimire/jsontology
```

## Documentation

- [Understanding Rule Layout](docs/rule.md)
- [Creating Your Own EventHandler](docs/event_handler.md)
- [Creating Your Own Operator](docs/operator.md)


## Future plans
The library is still in the initial phase, but here is a layout of future plans.

*  Common third party handlers like ones for sending email, message to slack will be integrated in library itself.
*  More operators built in for advanced processing of different kind of data type


#### Inspired By
* [jsonlogic](https://jsonlogic.com/)
