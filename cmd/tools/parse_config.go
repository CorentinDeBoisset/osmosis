package tools

import "fmt"
import "strings"
import "io/ioutil"
import "gopkg.in/yaml.v2"

type ServiceConfig struct {
    Src string          `yaml:"src"`
    Excludes []string   `yaml:"excludes"`
    UserId int          `yaml:"user_id"`
    GroupId int         `yaml:"group_id"`
    Image string        `yaml:"image"`
}

type OsmosisConfig struct {
    Syncs map[string]ServiceConfig `yaml:"syncs"`
}

func (c *OsmosisConfig) ParseConfig(filePath string) (err error) {
    yamlfile, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("File %s does not exist.", filePath)
    }

    err = yaml.Unmarshal(yamlfile, c)
    if err != nil {
        if yerr, ok := err.(*yaml.TypeError); ok {
            return fmt.Errorf("Format of %s is invalid for the following reasons:\n  - %s", filePath, strings.Join(yerr.Errors, "\n  - "))
        } else {
            return fmt.Errorf("Format of %s is invalid.", filePath)
        }
    }

    return nil
}
