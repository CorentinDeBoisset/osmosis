package commands

import "team-git.sancare.fr/dev/osmosis/cmd/tools"

func Status(projectName string, verbose bool) (err error) {
    err = tools.DockerConnect(verbose)
    if err != nil {
        return err
    }

    instances, err := tools.GetDockerInstances(verbose)
    if err != nil {
        return err
    }

    // TODO check running instances of unison, and check if they are OK

    return nil
}
