package commands

import "team-git.sancare.fr/dev/osmosis/cmd/tools"

func Start(serviceName string, config tools.OsmosisServiceConfig, verbose bool) (err error) {
    err = tools.DockerConnect(verbose)
    if err != nil {
        return err
    }

    _, err = tools.DockerContainerStart(serviceName, config, verbose)
    if err != nil {
        return err
    }

    return nil
}
