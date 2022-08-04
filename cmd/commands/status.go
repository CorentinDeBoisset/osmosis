package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/coreoas/osmosis/cmd/clients"
	"github.com/coreoas/osmosis/cmd/tools"
)

func Status(projectName string, config tools.OsmosisFullConfig, verbose bool) (err error) {
	err = clients.DockerConnect(verbose)
	if err != nil {
		return err
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(writer, "NAME\tSTATUS")

	for serviceName := range config.Syncs {
		serviceName := projectName + "_" + serviceName
		dockerInstance, err := clients.GetDockerInstance(serviceName, verbose)
		if err != nil {
			return err
		}
		unisonProcess := clients.GetUnisonInstance(serviceName)
		if dockerInstance == nil {
			fmt.Fprintln(writer, serviceName+"\t"+"uninitialized")
		} else if dockerInstance.Status == "running" {
			if unisonProcess.Pid != -1 {
				if unisonProcess.Running {
					fmt.Fprintln(writer, serviceName+"\t"+"running")
				} else {
					fmt.Fprintln(writer, serviceName+"\t"+"dead")
				}
			} else {
				fmt.Fprintln(writer, serviceName+"\t"+"error (container is running but not unison)")
			}
		} else {
			fmt.Fprintln(writer, serviceName+"\t"+"stopped")
		}
	}

	fmt.Fprintln(writer, "")
	writer.Flush()

	return nil
}
