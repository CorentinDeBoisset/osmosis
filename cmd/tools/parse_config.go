package tools

import (
    "errors"
    "fmt"
    "os"
    "io/ioutil"
    "path/filepath"
    "strings"
    "strconv"
    "os/user"
    "gopkg.in/yaml.v2"
)

type OsmosisServiceConfig struct {
    Src string          `yaml:"src"`
    Excludes []string   `yaml:"excludes"`
    UserId string       `yaml:"user_id"`
    GroupId string      `yaml:"group_id"`
    Image string        `yaml:"image"`
    VolumeName string   `yaml:"volume_name"`
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

    currentUser, err := user.Current()
    if err != nil {
        return errors.New("Could not read current user properties")
    }

    // Set default values for configuration
    for serviceName, serviceConf := range c.Syncs {
        if serviceConf.Image == "" {
            serviceConf.Image = "coenern/osmosis:alpha"
        }
        if serviceConf.Src == "" {
            serviceConf.Src = "."
        }
        if serviceConf.VolumeName == "" {
            serviceConf.VolumeName = serviceName
        }

        if serviceConf.UserId == "" {
            serviceConf.UserId = currentUser.Uid
        } else if _, err := strconv.Atoi(serviceConf.UserId); err != nil {
            return fmt.Errorf("%s is not valid user UID", serviceConf.UserId)
        }
        if serviceConf.GroupId == "" {
            serviceConf.GroupId = currentUser.Gid
        } else if _, err := strconv.Atoi(serviceConf.GroupId); err != nil {
            return fmt.Errorf("%s is not valid user GID", serviceConf.GroupId)
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
