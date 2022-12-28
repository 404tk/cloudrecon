package asset_check

import (
	"embed"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:embed pocs
var Pocs embed.FS

var defaultPocs []*Poc

type Poc struct {
	Name  string `yaml:"name"`
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Tag        string   `yaml:"tag"`
	Regions    []string `yaml:"regions"`
	Traversal  bool     `yaml:"region_traversal"`
	Method     string   `yaml:"method"`
	Format     []string `yaml:"format"`
	Expression []int    `yaml:"expression"`
}

func InitDefaultPoc() {
	entries, err := Pocs.ReadDir("pocs")
	if err != nil {
		log.Printf("[-] init defaultPoc error: %v\n", err)
		return
	}
	for _, one := range entries {
		path := one.Name()
		if strings.HasSuffix(path, ".yaml") {
			if poc, _ := loadPoc(path, Pocs); poc != nil {
				defaultPocs = append(defaultPocs, poc)
			}
		}
	}
}

func loadPoc(fileName string, Pocs embed.FS) (*Poc, error) {
	p := &Poc{}
	yamlFile, err := Pocs.ReadFile("pocs/" + fileName)
	if err != nil {
		log.Printf("[-] load poc %s error: %v", fileName, err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		log.Printf("[-] load poc %s error: %v", fileName, err)
		return nil, err
	}
	return p, err
}
