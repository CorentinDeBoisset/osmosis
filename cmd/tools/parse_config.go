package tools

import (
    "fmt"
    "strings"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type OsmosisServiceConfig struct {
    Src string          `yaml:"src"`
    Excludes []string   `yaml:"excludes"`
    UserId int          `yaml:"user_id"`
    GroupId int         `yaml:"group_id"`
    Image string        `yaml:"image"`
}

type OsmosisFullConfig struct {
    Syncs map[string]OsmosisServiceConfig `yaml:"syncs"`
}

func (c *OsmosisFullConfig) ParseConfig(filePath string) (err error) {
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

    // Set default values for configuration
    for serviceName, serviceConf := range c.Syncs {
        if serviceConf.Image == "" {
            serviceConf.Image = "coenern/osmosis:alpha"
        }
        if serviceConf.Src == "" {
            serviceConf.Src = "."
        }
        c.Syncs[serviceName] = serviceConf
    }

    return nil
}
