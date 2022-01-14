/*  Smart YAML functions - Helper functions to work with YAML
    (C) 2021-2022 Péter Deák (hyper80@gmail.com)
    License: Apache 2.0
*/

// smartyaml is a go package to handle parsed yaml files more confortable.
// It provides helper functions to parse, query and covert the yaml data.
package smartyaml

import (
	"github.com/hyper-prog/smartjsonyamlstub"

	"gopkg.in/yaml.v3"
)

// SmartYAML holds the parsed data. You can call the smartyaml functions on this structure
type SmartYAML struct {
	smartjsonyamlstub.SmartJsonYamlBase
}

// ParseYAML parse the raw yaml data (read from file) and returns a SmartYAML structure
func ParseYAML(rawdata []byte) (SmartYAML, error) {
	s := SmartYAML{}
	s.Config.InitConfig()
	err := yaml.Unmarshal(rawdata, &s.ParsedData)
	s.ParsedFrom = "yaml"
	return s, err
}

// String generate a yaml string
func (smartYaml SmartYAML) String() string {
	return smartYaml.Yaml()
}

// GetSubyamlByPath returns an another SmartYAML struct which holds a part of the parsed yaml specified by path
func (smartYaml SmartYAML) GetSubyamlByPath(path string) (SmartYAML, string) {
	cd, str := smartYaml.GetSubtreeByPath(path)
	rd := SmartYAML{}
	rd.Config = cd.Config
	rd.ParsedData = cd.ParsedData
	return rd, str
}
