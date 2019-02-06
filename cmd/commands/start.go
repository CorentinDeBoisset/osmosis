package commands

// import "fmt"
import "team-git.sancare.fr/dev/osmosis/cmd/tools"

func Start(projectName string, verbose bool) (err error) {
    err = tools.DockerConnect(verbose)
    if err != nil {
        return err
    }

    return nil
}
