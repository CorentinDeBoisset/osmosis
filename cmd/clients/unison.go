package clients

import (
    "fmt"
    "os"
    "io/ioutil"
    "strconv"
    "syscall"
)

import "team-git.sancare.fr/dev/osmosis/cmd/tools"


type OsmosisUnisonInstance struct {
    Pid int
    Running bool
}

func GetUnisonInstance(serviceName string) (instance *OsmosisUnisonInstance) {
    if piddata, err := ioutil.ReadFile("/tmp/osmosis/"+serviceName+".pid"); err == nil {
        if pid, err := strconv.Atoi(string(piddata)); err == nil {
            if process, err := os.FindProcess(pid); err == nil {
                if err := process.Signal(syscall.Signal(0)); err == nil {
                    return &OsmosisUnisonInstance{Pid: pid, Running: true}
                }
            }
            return &OsmosisUnisonInstance{Pid: pid, Running: false}
        }
    }

    return &OsmosisUnisonInstance{Pid: -1, Running: false}
}


func StartUnisonInstance(serviceName string, config tools.OsmosisServiceConfig, verbose bool) (instance *OsmosisUnisonInstance, err error) {
    instance = GetUnisonInstance(serviceName)
    if instance.Pid != -1 && instance.Running {
        return nil, fmt.Errorf("A unison instance is already running for %s", serviceName)
    }

    // TODO start unison process with correct parameters
    return nil, nil
}


func StopUnisonInstance(serviceName string) (err error) {
    return nil
}
