package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coreoas/osmosis/cmd/commands"
	"github.com/coreoas/osmosis/cmd/tools"
	"github.com/spf13/pflag"
)

func main() {
	var CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	var (
		file        string
		projectName string
		verbose     bool
	)

	CommandLine.StringVarP(&file, "file", "f", "osmosis.yml", "")
	CommandLine.StringVarP(&projectName, "project-name", "p", "", "")
	CommandLine.BoolVar(&verbose, "verbose", false, "")

	var err = CommandLine.Parse(os.Args[1:])
	if err != nil {
		commands.Help()
		os.Exit(1)
	}
	var args = CommandLine.Args()

	var osmosisConf tools.OsmosisFullConfig
	err = osmosisConf.ParseConfig(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if projectName == "" {
		fullPath, err := filepath.Abs(file)
		if err != nil {
			fmt.Printf("Could not determine %s absolute directory.\n\n", file)
			os.Exit(1)
		}
		projectName = filepath.Base(filepath.Dir(fullPath))
	}

	if len(args) == 1 {
		switch args[0] {
		case "start":
			for serviceName, serviceConfig := range osmosisConf.Syncs {
				fullName := projectName + "_" + serviceName
				fmt.Printf("Starting service %s... ", serviceName)
				err = commands.Start(fullName, serviceConfig, verbose)
				if err != nil {
					fmt.Printf("\nError: %s\n\n", err)
					os.Exit(1)
				}
				fmt.Println("Done")
			}
		case "stop":
			for serviceName := range osmosisConf.Syncs {
				fullName := projectName + "_" + serviceName
				fmt.Printf("Stopping service %s...", serviceName)
				err = commands.Stop(fullName, verbose)
				if err != nil {
					fmt.Printf("\nError: %s\n\n", err)
					os.Exit(1)
				}
				fmt.Println("Done")
			}
		case "status":
			err = commands.Status(projectName, osmosisConf, verbose)
		case "clean":
			for serviceName, serviceConfig := range osmosisConf.Syncs {
				fullName := projectName + "_" + serviceName
				fmt.Printf("Cleaning service %s...", serviceName)
				err = commands.Clean(fullName, serviceConfig, verbose)
				if err != nil {
					fmt.Printf("\nError: %s\n\n", err)
					os.Exit(1)
				}
				fmt.Println("Done")
			}
		case "help":
			commands.Help()
		case "version":
			// TODO: improve version management
			fmt.Println("osmosis 0.1.0")
		default:
			commands.InvalidCommand(args[0])
			os.Exit(1)
		}
	} else {
		commands.Help()
	}
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		os.Exit(1)
	}
}
