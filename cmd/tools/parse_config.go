package tools

import (
    "fmt"
    "os"
    "io/ioutil"
    "path/filepath"
    "strings"
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

func (c *OsmosisFullConfig) ParseConfig(configPath string) (err error) {
    yamlfile, err := ioutil.ReadFile(configPath)
    if err != nil {
        return fmt.Errorf("File %s does not exist.", configPath)
    }

    err = yaml.Unmarshal(yamlfile, c)
    if err != nil {
        if yerr, ok := err.(*yaml.TypeError); ok {
            return fmt.Errorf("Format of %s is invalid for the following reasons:\n  - %s", configPath, strings.Join(yerr.Errors, "\n  - "))
        } else {
            return fmt.Errorf("Format of %s is invalid.", configPath)
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
        if !filepath.IsAbs(serviceConf.Src) {
            serviceConf.Src, err = filepath.Abs(filepath.Dir(configPath) + "/" + serviceConf.Src)
            if err != nil {
                return fmt.Errorf("Could not calculate src path \"%s\" in sync \"%s\".", serviceConf.Src, serviceName)
            }
        }

        if _, err = os.Stat(serviceConf.Src); os.IsNotExist(err) {
            return fmt.Errorf("Path \"%s\" in sync \"%s\" does not exist", serviceConf.Src, serviceName)
        }
        c.Syncs[serviceName] = serviceConf
    }

    return nil
}
