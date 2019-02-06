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
        host string
        tls bool
        tlscacert string
        tlscert string
        tlskey string
        tlsverify bool
        skipHostnameCheck bool
    )

    CommandLine.StringVarP(&file, "file", "f", "osmosis.yml", "")
    CommandLine.StringVarP(&host, "host", "H", "/var/run/docker.sock", "")
    CommandLine.BoolVar(&tls, "tls", false, "")
    CommandLine.StringVar(&tlscacert, "tlscacert", "", "")
    CommandLine.StringVar(&tlscert, "tlscert", "", "")
    CommandLine.StringVar(&tlskey, "tlskey", "", "")
    CommandLine.BoolVar(&tlsverify, "tlsverify", false, "")
    CommandLine.BoolVar(&skipHostnameCheck, "skip-hostname-check", false, "")

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
            commands.Start()
        case "stop":
            commands.Stop()
        case "status":
            commands.Status()
        case "restart":
            commands.Stop()
            commands.Start()
        case "clean":
            commands.Clean()
        case "help":
            commands.Help()
        default:
            commands.InvalidCommand(args[0])
            os.Exit(1)
        }
    } else {
        commands.Help()
    }
}
