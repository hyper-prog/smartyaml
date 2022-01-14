/*  (C) 2021-2022 Péter Deák (hyper80@gmail.com)
 */

package smartyaml

import (
	"strings"
	"testing"
)

var sampleyaml1 string = `
---
sample:
 classes:
  one:
   name: Red
   descr: The big red group
  two:
   name: Blue
   descr: The big blue group
 members:
  - nick: Frank
    number: 43
  - nick: Dave
    number: 56
  - nick: Joe
    number: 12
`

func TestNodeExists(t *testing.T) {
	sy, error := ParseYAML([]byte(sampleyaml1))
	if error != nil {
		t.Errorf("Error at parsing: %s", error)
	}

	questions := []string{
		"$.sample.classes.two.name",
		"/sample/classes/one/name",
		"$.sample.two.weight",
		"/sample/class/two/weight",
	}

	answers := []bool{
		true,
		true,
		false,
		false,
	}

	for i := 0; i < len(questions); i++ {
		if sy.NodeExists(questions[i]) != answers[i] {
			t.Errorf("NodeExist test failed for node %s", questions[i])
		}
	}
}

func TestNodeType(t *testing.T) {
	sy, error := ParseYAML([]byte(sampleyaml1))
	if error != nil {
		t.Errorf("Error at parsing: %s", error)
	}

	questions := []string{
		"$.sample.classes.two.name",
		"/sample/classes/one",
		"$.sample.members",
		"$.sample.members[0].number",
		"$.sample.members[7].number",
	}

	answers := []string{
		"string",
		"map",
		"array",
		"int",
		"none",
	}

	for i := 0; i < len(questions); i++ {
		_, typ := sy.GetNodeByPath(questions[i])
		if typ != answers[i] {
			t.Errorf("Type query test failed for node %s", questions[i])
		}
	}
}

func TestGetStringByPath(t *testing.T) {
	sy, error := ParseYAML([]byte(sampleyaml1))
	if error != nil {
		t.Errorf("Error at parsing: %s", error)
	}

	probes := [][]string{
		{"$.sample.classes.two.name", "Blue"},
		{"/sample/classes/one/name", "Red"},
		{"$.sample.classes.one.descr", "The big red group"},
		{"/sample/members/[0]/nick", "Frank"},
		{"$.sample.members[1].nick", "Dave"},
	}

	for i := 0; i < len(probes); i++ {
		result, tp := sy.GetStringByPath(probes[i][0])
		if tp != "string" || result != probes[i][1] {
			t.Errorf("GetStringByPath test failed for node %s (%s)", probes[i][0], tp)
		}
	}
}

func TestGetStringByPathWithDefault(t *testing.T) {
	sy, error := ParseYAML([]byte(sampleyaml1))
	if error != nil {
		t.Errorf("Error at parsing: %s", error)
	}

	probes := [][]string{
		{"$.sample.classes.two.name", "Blue"},
		{"smaple/classes/three/name", "fallbackvalue"},
		{"sample/classes/one/name", "Red"},
		{"classes/one/name", "fallbackvalue"},
		{"$.sample.classes.two.descr", "The big blue group"},
		{"/sample/members/[2]/nick", "Joe"},
		{"$.sample.members[4].nick", "fallbackvalue"},
	}

	for i := 0; i < len(probes); i++ {
		result := sy.GetStringByPathWithDefault(probes[i][0], "fallbackvalue")
		if result != probes[i][1] {
			t.Errorf("GetStringByPathWithDefault test failed for node %s", probes[i][0])
		}
	}
}

func TestGetIntegerByPath(t *testing.T) {
	sy, error := ParseYAML([]byte(sampleyaml1))
	if error != nil {
		t.Errorf("Error at parsing: %s", error)
	}

	questions := []string{
		"$.sample.members[0].number",
		"/sample/members/[2]/number",
		"sample/members/[1]/number",
	}

	answers := []int{
		43,
		12,
		56,
	}

	for i := 0; i < len(questions); i++ {
		result, tp := sy.GetIntegerByPath(questions[i])
		if tp != "int" || result != answers[i] {
			t.Errorf("GetIntegerByPath test failed for node %s (%s)", questions[i], tp)
		}
	}
}

func TestGetIntegerByPathWithDefault(t *testing.T) {
	sy, error := ParseYAML([]byte(sampleyaml1))
	if error != nil {
		t.Errorf("Error at parsing: %s", error)
	}

	questions := []string{
		"$.sample.members[0].number",
		"sample/members/[4]/number",
		"sample/members/[1]/number",
		"nonexistent/some",
	}

	answers := []int{
		43,
		97,
		56,
		97,
	}

	for i := 0; i < len(questions); i++ {
		result := sy.GetIntegerByPathWithDefault(questions[i], 97)
		if result != answers[i] {
			t.Errorf("GetIntegerByPathWithDefault test failed for node %s", questions[i])
		}
	}
}

func TestGetCountDescendantsByPath(t *testing.T) {
	sy, error := ParseYAML([]byte(sampleyaml1))
	if error != nil {
		t.Errorf("Error at parsing: %s", error)
	}

	questions := []string{
		"$.sample",
		"sample/members",
		"sample/classes/one",
		"sample/nonexistent/some",
		"$.sample.members[2]",
	}

	answers := []int{
		2,
		3,
		2,
		0,
		2,
	}

	for i := 0; i < len(questions); i++ {
		result := sy.GetCountDescendantsByPath(questions[i])
		if result != answers[i] {
			t.Errorf("GetCountDescendantsByPath test failed for node %s", questions[i])
		}
	}
}

func TestConvertToJson(t *testing.T) {
	fromyaml := []string{
		`
---
top:
  sub:
    subsub: Value down
`,
		`
---
arrayitems:
- red:
    sample: apple
- green:
    sample: pear
- blue:
    sample: plum
`,
		`
---
arrayitems:
- red: apple
- green:
    subspec: 56.47
- blue: plum
`,
		`
---
first:
  second:
    third: apple
`,
		`
---
toplevel:
  arrayitems:
  - red:
    - apple: 3
    - cherry: 4
    - strawberry: 5
  - green:
    - pear: 1
    - cucumber: 6
    - grape: 7
  - blue:
    - plum: 4
    - bluegrape: 32
`,
	}

	tojson := []string{
		`{"top":{"sub":{"subsub":"Value down"}}}`,
		`{"arrayitems":[{"red":{"sample":"apple"}},{"green":{"sample":"pear"}},{"blue":{"sample":"plum"}}]}`,
		`{"arrayitems":[{"red":"apple"},{"green":{"subspec":56.47}},{"blue":"plum"}]}`,
		`{"first":{"second":{"third":"apple"}}}`,
		`{"toplevel":{"arrayitems":[{"red":[{"apple":3},{"cherry":4},{"strawberry":5}]},{"green":[{"pear":1},{"cucumber":6},{"grape":7}]},{"blue":[{"plum":4},{"bluegrape":32}]}]}}`,
	}

	for i := 0; i < len(fromyaml); i++ {
		sy, error := ParseYAML([]byte(fromyaml[i]))
		if error != nil {
			t.Errorf("Conversion test, error at parsing: %s", error)
		}
		result := sy.JsonCompacted()
		if strings.TrimSpace(result) != strings.TrimSpace(tojson[i]) {
			t.Errorf("Json conversion test failed with pattern no.%d", i+1)
		}
	}
}
