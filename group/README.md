# Group Finder
Group finder is a package that adds another DSL to improve the maintainability 
of the searched patterns and enables searches on specific fields of structured documents.

## Finder
The finder is used to manage multiple rules. It will use the DSL along with the gofindthem finder 
to verify if the tags where found on the specified fields.

### Usage
You will need 2 sets of expressions, one to define the patterns that are needed to 
be searched with its tag and the second one with the expressions to define the tags
 and field relations.

```golang
	gofindthemRules := map[string][]string{
		"tag1": {
			`"string1"`,
			`"string2"`,
		},
		"tag2": {
			`"string3"`,
			`"string4"`,
		},
		"tag3": {
			`"string5"`,
			`"string6"`,
		},
		"tag4": {
			`"string7"`,
			`"string8"`,
		},
	}

	rules := map[string][]string{
		"rule1": {`"tag1" or "tag2" and not "tag3"`},
		"rule2": {`"tag3:Field3.SomeField1" or "tag4"`},
		"rule3": {`"tag3:Field3" or "tag4"`},
	}
```

With the 2 sets of expressions ready you will first need to create the 
gofinthem finder and the group finder:

```golang
	gft, err := finder.NewFinderWithExpressions(
		&finder.CloudflareForkEngine{},
		&finder.RegexpEngine{},
		false,
		gofindthemRules,
	)

    gftg, err := gfinder.NewFinderWithRules(gft, rules)
	if err != nil {
		panic(err)
	}
```

Now its possible check which rules where evaluated as true on a text
or on a structured document:

```golang
    // searching on a struct
	res, err := gftg.ProcessObject(someObject, gftg.GetFieldNames(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessObject: ", res)

    // searching on a raw json
    res3, err := gftg.ProcessJson(rawJson, gftg.GetFieldNames(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessJson: ", res3)
```
The full example can be found at `/examples/group/finder/main.go`

## Group Finder DSL
### Definition
The DSL uses 3 operators (AND, OR, NOT), Tag (defined by "tag:(field)"), 
where the field is optional, and parentheses to form expressions.
A valid expression can be:

- A single rule with or without a specific field. Eg: `"tag1"` `"tag1:field1"`
- The result of an operation. `"tag1" OR "tag2:field1"` 
- An expression enclosed by parentheses `("tag1" OR "tag2:field1")`

Each operator functions as the following:

- **AND** - Uses the expression before and after it to solve them as a logical `AND` operator. 
    > (valid expression) AND (valid expression) eg: `"tag1" AND "tag2"` 

- **OR** - Uses the expression before and after it to solve them as a logical `OR` operator.
    > \<valid expression\> OR \<valid expression\> eg: `"tag1" OR "tag2"` 

- **NOT** - Uses the expression after it to solve them as a logical `NOT` operator.
    > NOT \<valid expression\> eg: `NOT "tag1"`
