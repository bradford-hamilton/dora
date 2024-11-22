<div align="center">
  <img
    alt="Dora backpack with JSON"
    src="./dora.png"
    height="300px"
  />
</div>
<h1 align="center">Welcome to dora the JSON explorer 👋</h1>
<p align="center">
  <a href="https://golang.org/dl" target="_blank">
    <img alt="Using go version 1.23.1" src="https://img.shields.io/badge/go-1.23.1-9cf.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/bradford-hamilton/dora" target="_blank">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/bradford-hamilton/dora" />
  </a>
  <a href="https://godoc.org/github.com/bradford-hamilton/dora/pkg" target="_blank">
    <img alt="godoc" src="https://godoc.org/github.com/bradford-hamilton/dora/pkg?status.svg" />
  </a>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> Dora makes exploring JSON fast, painless, and elegant.

### **_NOTE_**:
- dora's main focus was for learning and teaching the content through a [medium blog post](https://medium.com/@bradford_hamilton/building-a-json-parser-and-query-tool-with-go-8790beee239a). In other words, dora is not currently a stable tool.

## Install

```sh
go get github.com/bradford-hamilton/dora/pkg/dora
```

## Usage
```go
var exampleJSON = `{ "string": "a neat string", "bool": true, "PI": 3.14159 }`

c, err := dora.NewFromString(exampleJSON)
if err != nil {
  return err
}

str, err := c.GetString("$.string")
if err != nil {
  return err
}

boolean, err := c.GetBool("$.bool")
if err != nil {
  return err
}

float, err := c.GetFloat64("$.PI")
if err != nil {
  return err
}

fmt.Println(str)     // a neat string
fmt.Println(boolean) // true
fmt.Println(float)   // 3.14159
```

## Query Syntax

1. All queries start with `$`.

2. Access objects with `.` only, no support for object access with bracket notation `[]`.
    - This is intentional, as you can interpolate at the call site, so there is no reason to offer two syntaxes that do the same thing.

3. Access arrays by index with bracket notation `[]`.

4. **New**: Fetch by type to allow caller to ask for the proper Go type. For the time being, asking for objects or arrays in their entirety must be done through `GetString` which will return the chunk of JSON.
    
    Current API:
    - `GetString`
    - `GetFloat64`
    - `GetBool`

5. Next feature will be approaching this either with some sort of serialization option maybe similar to stdlib or a simpler one with no options that returns a map or something? Will think about that some.

 Example with a JSON object as root value:
```js
JSON:

{
  "name": "bradford",
  "someArray": ["some", "values"]
  "obj": {
    "innerKey": {
      "innerKey2": "innerValue",
      "innerKey3": [{ "kindOfStuff": "neatStuff" }]
    }
  },
  "someBool": true,
  "PI": 3.14159
}

Query:                                  Result:

$.name                                  == "bradford"
$.someArray                             == "[\"array\", \"values\"]"
$.someArray[0]                          == "some"
$.someArray[1]                          == "values"
$.someArray[2]                          == error
$.obj.innerKey.innerKey2                == "innerValue"
$.obj.innerKey.innerKey3[0].kindOfStuff == "neatStuff"
$.someBool                              == true
$.PI                                    == 3.14159
```

 Example with a JSON array as root value:
```js
JSON:

[
  "some",
  "values",
  {
    "objKey": "objValue",
    "objKey2": [{ "catstack": "lampcat" }]
  }
]

Query:                   Result:

$[0]                     == "some"
$[1]                     == "values"
$[2]                     == "{ \"objKey\": \"objValue\", \"objKey2\": [{ \"catstack\": \"lampcat\" }] }"
$[2].objKey              == "objValue"
$[2].objKey2[0]          == "{ \"catstack\": \"lampcat\" }"
$[2].objKey2[0].catstack == "lampcat"
```

## Run tests

```shs
go test ./...
```

## Author

👤 **Bradford Lamson-Scribner**

* Website: bradfordhamilton.io
* Twitter: [@lamsonscribner](https://twitter.com/lamsonscribner)
* Github: [@bradford-hamilton](https://github.com/bradford-hamilton)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/bradford-hamilton/dora/issues). 

## Show your support

Give a ⭐️ if this project helped you!

