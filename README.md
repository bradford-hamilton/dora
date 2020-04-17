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
    <img alt="Using go version 1.14" src="https://img.shields.io/badge/go-1.14-9cf.svg" />
  </a>
  <a href="https://travis-ci.com/bradford-hamilton/dora" target="_blank">
    <img alt="Using go version 1.14" src="https://travis-ci.com/bradford-hamilton/dora.svg?branch=master" />
  </a>
  <a href="https://goreportcard.com/report/github.com/bradford-hamilton/dora" target="_blank">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/bradford-hamilton/dora/pkg" />
  </a>
  <a href="https://godoc.org/github.com/bradford-hamilton/dora/pkg" target="_blank">
    <img alt="godoc" src="https://godoc.org/github.com/bradford-hamilton/dora/pkg?status.svg" />
  </a>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> Dora makes exploring JSON fast, painless, and elegant. (total lie right now, it's 2-3x slower than encoding/json, painful because it's way unfinished, and ugly because there is no thought out API yet)

### NOTE: Please don't use yet, it's currently a WIP. Right now dora lexes/scans and parses JSON into an AST. Spiked out a `GetByPath` call which uses a query syntax similar to javascript or JSONPath. This one feature so far is getting close to functional/stable. May also expand project to have an AST explorer of some sort.

## Install

```sh
go get github.com/bradford-hamilton/dora/pkg/dora
```

## Usage
```go
c, err := dora.NewFromString("{\"someObj\": [1, { \"neatKey\": \"neatVal\" }, true] }")
if err != nil {
  fmt.Printf("\nError creating client: %v\n", err)
}

result, err := c.GetByPath("$.someObj[1].neatKey") // result == "neatVal"
if err != nil {
  fmt.Printf("\nError executing GetByPath query: %v\n", err)
}
```

## Query Syntax

1. All queries start with `$`.
2. Access objects with `.` only, no support for object access with bracket notation `[]`
    - This is intentional, as you can interpolate at the call site, so there is no reason to offer two syntaxes that do the same thing.
3. Access arrays by index with bracket notation `[]`

 Example with a JSON object as root value:
```js
{
  "name": "bradford",
  "someArray": ["some", "values"]
  "obj": {
    "innerKey": {
      "innerKey2": "innerValue",
      "innerKey3": [{ "kindOfStuff": "neatStuff" }]
    }
  }
}

$.name                                  == "bradford"
$.someArray                             == "[\"array\", \"values\"]"
$.someArray[0]                          == "some"
$.someArray[1]                          == "values"
$.someArray[2]                          == error
$.obj.innerKey.innerKey2                == "innerValue"
$.obj.innerKey.innerKey3[0].kindOfStuff == "neatStuff"
```

 Example with a JSON array as root value:
```js
[
  "some",
  "values",
  {
    "objKey": "objValue",
    "objKey2": [{ "catstack": "lampcat" }]
  }
]

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

* Website: https://www.bradfordhamilton.io
* Twitter: [@lamsonscribner](https://twitter.com/lamsonscribner)
* Github: [@bradford-hamilton](https://github.com/bradford-hamilton)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/bradford-hamilton/dora/pkg/issues). 

## Show your support

Give a ⭐️ if this project helped you!

## Credit
Started playing with this idea after finding:
1. [jsonparser](https://github.com/buger/jsonparser)
2. [ast-explorer](https://astexplorer.net)

***
##### _This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_
