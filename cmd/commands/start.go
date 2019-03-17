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

    if err := clients.DockerVolumeCreate(config.VolumeName, verbose); err != nil {
        return err
    }

    dockerInstance, err := clients.DockerContainerStart(serviceName, config, verbose)
    if err != nil {
        return err
    }

    if _, err = clients.StartUnisonInstance(serviceName, config, dockerInstance.Port, verbose); err != nil {
        return err
    }

    return nil
}
