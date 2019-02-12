package commands

import (
    "team-git.sancare.fr/dev/osmosis/cmd/tools"
    "team-git.sancare.fr/dev/osmosis/cmd/clients"
)

func Start(serviceName string, config tools.OsmosisServiceConfig, verbose bool) (err error) {
    err = clients.DockerConnect(verbose)
    if err != nil {
        return err
    }

    _, err = clients.DockerContainerStart(serviceName, config, verbose)
    if err != nil {
        return err
    }
    // TODO start unison instance

    return nil
}
