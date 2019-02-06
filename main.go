package main

import "fmt"
import "os"
import "github.com/spf13/pflag"
import "team-git.sancare.fr/dev/osmosis/cmd/commands"
import "team-git.sancare.fr/dev/osmosis/cmd/tools"

func main() {
    var CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
    var (
        file string
        projectName string
        verbose bool
    )

    CommandLine.StringVarP(&file, "file", "f", "osmosis.yml", "")
    CommandLine.StringVarP(&projectName, "project-name", "p", "", "")
    CommandLine.BoolVar(&verbose, "verbose", false, "")

    var err = CommandLine.Parse(os.Args[1:])
    if (err != nil) {
        commands.Help()
        os.Exit(1)
    }
    var args = CommandLine.Args()

    var osmosisConf tools.OsmosisConfig
    err = osmosisConf.ParseConfig(file)
    if err != nil {
        fmt.Println(err, "\n")
        os.Exit(1)
    }

    if (len(args) == 1) {
        switch args[0] {
        case "start":
            err = commands.Start(projectName, verbose)
        case "stop":
            err = commands.Stop(projectName, verbose)
        case "status":
            err = commands.Status(projectName, verbose)
        case "restart":
            err = commands.Stop(projectName, verbose)
            if err != nil {
                fmt.Printf("Error: %s\n\n", err);
                os.Exit(1)
            }
            err = commands.Start(projectName, verbose)
        case "clean":
            err = commands.Clean(projectName, verbose)
        case "help":
            commands.Help()
        default:
            commands.InvalidCommand(args[0])
            os.Exit(1)
        }
    } else {
        commands.Help()
    }
    if err != nil {
        fmt.Printf("Error: %s\n\n", err);
        os.Exit(1)
    }
}
