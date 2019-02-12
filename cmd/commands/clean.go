package commands

// import "fmt"
import (
    "team-git.sancare.fr/dev/osmosis/cmd/clients"
)

func Clean(projectName string, verbose bool) (err error) {
    err = clients.DockerConnect(verbose)
    if err != nil {
        return err
    }

    return nil
}
