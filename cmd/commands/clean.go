package commands

import (
	"github.com/coreoas/osmosis/cmd/clients"
	"github.com/coreoas/osmosis/cmd/tools"
)

func Clean(projectName string, config tools.OsmosisServiceConfig, verbose bool) (err error) {
	err = clients.DockerConnect(verbose)
	if err != nil {
		return err
	}

	if err = clients.StopUnisonInstance(projectName); err != nil {
		return err
	}

	if err = clients.DockerContainerStop(projectName, verbose); err != nil {
		return err
	}

	if err = clients.DockerContainerRemove(projectName, verbose); err != nil {
		return err
	}

	if err := clients.DockerVolumeRemove(config.VolumeName, verbose); err != nil {
		return err
	}

	return nil
}
