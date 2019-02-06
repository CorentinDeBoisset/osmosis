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
  -f, --file FILE
  -H, --host HOST

  --tls                       Use TLS; implied by --tlsverify
  --tlscacert CA_PATH         Trust certs signed only by this CA
  --tlscert CLIENT_CERT_PATH  Path to TLS certificate file
  --tlskey TLS_KEY_PATH       Path to TLS key file
  --tlsverify                 Use TLS and verify the remote
  --skip-hostname-check       Don't check the daemon's hostname against the
                              name specified in the client certificate
`)

    printCommandList()
}

func InvalidCommand(cmd string) {
    fmt.Println("No such command:", cmd, "\n")
    printCommandList()
}
