package commands

import (
    "os"
    "fmt"
    "text/tabwriter"
)

import (
    "team-git.sancare.fr/dev/osmosis/cmd/tools"
    "team-git.sancare.fr/dev/osmosis/cmd/clients"
)


func Status(projectName string, config tools.OsmosisFullConfig, verbose bool) (err error) {
    err = clients.DockerConnect(verbose)
    if err != nil {
        return err
    }

    writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
    fmt.Fprintln(writer, "NAME\tSTATUS")

    for serviceName, _ := range config.Syncs {
        serviceName := projectName+"_"+serviceName
        dockerInstance, err := clients.GetDockerInstance(serviceName, verbose)
        if err != nil {
            return err
        }
        unisonProcess := clients.GetUnisonInstance(serviceName)
        if dockerInstance.Status == "running" {
            if unisonProcess.Pid == -1 {
                if unisonProcess.Running {
                    fmt.Fprintln(writer, serviceName+"\t"+"running")
                } else {
                    fmt.Fprintln(writer, serviceName+"\t"+"dead")
                }
            } else {
                fmt.Fprintln(writer, serviceName+"\t"+"error")
            }
        } else {
            fmt.Fprintln(writer, serviceName+"\t"+"stopped")
        }
    }

    fmt.Fprintln(writer, "")
    writer.Flush()

    return nil
}
