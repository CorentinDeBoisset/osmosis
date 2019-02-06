package tools

import "fmt"
import "io/ioutil"
import "gopkg.in/yaml.v2"

type ServiceConfig struct {
    Src string          `yaml:"src"`
    Excludes []string   `yaml:"excludes"`
    UserId int          `yaml:"user_id"`
    GroupId int         `yaml:"group_id"`
}

type OsmosisConfig struct {
    Services map[string]ServiceConfig `yaml:"syncs"`
}

func (c *OsmosisConfig) ParseConfig(filePath string) (err error) {
    yamlfile, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("File %s does not exist.", filePath)
    }

    err = yaml.Unmarshal(yamlfile, c)
    if err != nil {
        return fmt.Errorf("Error parsing %s, format is invalid.", filePath)
    }

    return nil
}
