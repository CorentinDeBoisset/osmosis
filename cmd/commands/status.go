package commands

import "team-git.sancare.fr/dev/osmosis/cmd/tools"

func Status(projectName string, config tools.OsmosisFullConfig, verbose bool) (err error) {
    err = tools.DockerConnect(verbose)
    if err != nil {
        return err
    }

    for serviceName, _ := range config.Syncs {
        _, err := tools.GetDockerInstance(projectName+"_"+serviceName, verbose)
        if err != nil {
            return err
        }
        // TODO check running instances of unison, and check if they are OK
    }

    return nil
}
