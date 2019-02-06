package commands

import "fmt"

func printCommandList() {
    fmt.Printf(`Commands:
  start         Create and start sync services
  ps            List sync services
  stop          Stop services
  restart       Restart services
  clean         Remove all cache and volumes related to a service
  help          Get help about Osmosis
  version       Show the Osmosis version
`)
}

func Help() {
    fmt.Printf(`Define and manage file synchronisation services with docker containers.

Usage:
  osmosis <command> [<args>]

Options:
  -f, --file FILE           Specify an alternate osmosis file
                            (default: osmosis.yml)
  -p, --project-name NAME   Specify an alternate project name
                            (default: directory name)
  --verbose                 Show more output
`)

    printCommandList()
}

func InvalidCommand(cmd string) {
    fmt.Println("No such command:", cmd, "\n")
    printCommandList()
}
