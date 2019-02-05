package main

import "os"
import "team-git.sancare.fr/dev/osmosis/commands"

func main() {
    if (len(os.Args) == 1) {
        commands.Help()
    } else {
        switch os.Args[1] {
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
            commands.InvalidCommand(os.Args[1])
        }
    }
}
