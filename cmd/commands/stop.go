package commands

// import "fmt"

import (
    "team-git.sancare.fr/dev/osmosis/cmd/clients"
)

func Stop(projectName string, verbose bool) (err error) {
    if err = clients.DockerConnect(verbose); err != nil {
        return err
    }

    if err = clients.StopUnisonInstance(projectName); err != nil {
        return err
    }

    if err = clients.DockerContainerStop(projectName, verbose); err != nil {
        return err
    }

    return nil
}
