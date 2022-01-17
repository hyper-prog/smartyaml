[![Go Reference](https://pkg.go.dev/badge/github.com/hyper-prog/smartyaml.svg)](https://pkg.go.dev/github.com/hyper-prog/smartyaml)

SmartYAML - Go package to handle YAML
======================================

The smartyaml is a go package to handle parsed YAML files more confortable.
This package is not a parser, it uses the "gopkg.in/yaml.v3" package as parser.
It gives you some convenient helper functions to query/view/convert parsed YAML data structures.

In order to use it you have to get the package

	go get github.com/hyper-prog/smartyaml

*It seems that "go get ..." does not always work above go version 1.16. 
To workaround this set the "GO111MODULE=auto" environment variable during go get.*

Check out the https://github.com/hyper-prog/smartjson package to work with JSON

Sample code:

```go
package main

import (
	"fmt"
	"github.com/hyper-prog/smartyaml"
)

var sampleyaml string = `
firstName: John
lastName: Smith
isAlive: Yes
age: 27
address:
  streetAddress: 21 2nd Street
  city: New York
  state: NY
  postalCode: "10021-3100"
phoneNumbers:
- type: home
  number: "212 555-1234"
- type: office
  number: "646 555-4567"
children:
spouse:
`

func main() {
	/* // Yaml load from file:
	yamlData, flerr := ioutil.ReadFile("test.yml")
	if flerr != nil {
		fmt.Println("Error, cannot read file: ", flerr.Error())
	}
	sj, parsererror := smartyaml.ParseYAML(yamlData)
	*/

	sy, parsererror := smartyaml.ParseYAML([]byte(sampleyaml))
	if parsererror != nil {
		fmt.Println("Error, not valid YAML: ", parsererror.Error())
	}

	fmt.Println("City: ", sy.GetStringByPathWithDefault("$.address.city", "Unknown"))
	fmt.Println("First phone number: ", sy.GetStringByPathWithDefault("/phoneNumbers/[0]/number", "Not available"))
	_, agetype := sy.GetNodeByPath("age")
	fmt.Println("Type of Age: ", agetype)
	fmt.Println("Json:")
	fmt.Println()
	fmt.Println(sy.JsonIndented())
}
```

Query path types
----------------

You can query the nodes like directories, separated with `/` where the root sign is optional.

You can also use JsonPath to query a value.
JsonPath mode is automatic if the query string starts with "$." or "JsonPath:" prefix

	/members/[0]/name - Name of the first member
	$.members[2]/name - Name of the third member
	/config/items/[]/description  - The description of the last config item
	$.config.server.address - The address of the configured server

Available functions
-------------------

| Function                                         | Description                                 |
| ------------------------------------------------ | ------------------------------------------- |
| `ParseYAML(rawdata []byte)`                      | Parse the YAML data                         |
| `Yaml()`                                         | Generates a YAML string                     |
| `JsonIndented()`                                 | Generates an indented JSON string           |
| `JsonCompacted()`                                | Generates a compacted JSON string           |
| `NodeExists(path string)`                        | True if the given node exists               |
| `GetCountDescendantsByPath(path string)`         | Gives the number of descendants of the node |
| `GetNodeByPath(path string)`                     | Gives the node and node type by path        |
| `GetStringByPath(path string)`                   | String value of the requested node          |
| `GetNumberByPath(path string)`                   | Float or int value as float of the requested node  |
| `GetFloat64ByPath(path string)`                  | Float value of the requested node           |
| `GetIntegerByPath(path string)`                  | Integer value of the requested node         |
| `GetTimeByPath(path string)`                     | Time value of the requested node            |
| `GetBoolByPath(path string)`                     | Bool value of the requested node            |
| `GetStringByPathWithDefault(path string, def string)`   | String value of the requested node with fallback value  |
| `GetNumberByPathWithDefault(path string, def float64)`  | Float or int value as float of the requested node with fallback value |
| `GetFloat64ByPathWithDefault(path string, def float64)` | Float value of the requested node with fallback value   |
| `GetIntegerByPathWithDefault(path string, def int)`     | Integer value of the requested node with fallback value |
| `GetTimeByPathWithDefault(path string, def time.Time)`  | Time value of the requested node with fallback value    |
| `GetBoolByPathWithDefault(path string, def bool)`       | Bool value of the requested node with fallback value    |
| `GetMapByPath(path string)`                      | The map on the path                         |
| `GetArrayByPath(path string)`                    | The array on the path                       |
| `GetSubyamlByPath(path string)`                  | Gives the sub yaml specified by path        |


Author, License
---------------

The package is written by Peter Deak (C) hyper80@gmail.com under Apache 2.0 license
