package tools

import "fmt"
import "errors"
import "context"
import "github.com/docker/docker/client"
// import "github.com/docker/docker/api/types"

type OsmosisDockerInstance struct {
    Id string
    Name string
    Status string
    Port int
}

var cli *client.Client

func DockerConnect(verbose bool) (err error) {
    cli, err = client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        fmt.Println(err)
        return err
    }

    ping, err := cli.Ping(context.Background())
    if err != nil {
        return errors.New("Could not connect to the docker daemon. Check it is running and the DOCKER_HOST env var if it is not listening to unix:///var/run/docker.sock")
    }

    if verbose {
        fmt.Printf("Docker is running: \n  - Api version: %s\n  - OS type: %s\n\n", ping.APIVersion, ping.OSType)
    }

    return nil
}

func DockerStatus(packageName string, verbose bool) (statuses []OsmosisDockerInstance, err error){
    if cli == nil {
        return nil, errors.New("Docker client is not initialized.")
    }

    containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
    if err != nil {
        return nil, errors.New("Could not read list of containers")
    }

    instances = []OsmosisDockerInstance
    for _, container := range containers {
        if container.Image == "registry.sancare.fr/base_images/unison:1.0" && len(container.Ports) == 1 && len(container.Names) > 0 {
            instances = append(instances, OsmosisDockerInstance{Id: container.ID, Name: container.Names[0], Status: container.Status, Port: container.Ports[0]})
        }
    }

    return instances, nil
}

func DockerContainerStart(packageName string, containerName string, verbose bool) (instance OsmosisDockerInstance, err error) {
    if cli == nil {
        return nil, errors.New("Docker client is not initialized.")
    }

    instances, err := DockerStatus(packageName, verbose)
    if err != nil {
        return nil, err
    }

    ctx := context.Background()

    fullName := packageName+"_"+containerName
    if len(instances) > 0 {
        for _, instance := range instances {
            if instance.ContainerName == fullName {
                // TODO Handle containerStatus correctly
                if instance.ContainerStatus != "running" {
                    // TODO fix start options
                    err = cli.ContainerStart(ctx, instance.ContainerId, types.ContainerStartOptions{})
                    if err != nil {
                        return nil, fmt.Errorf("Container %s could not be started.", container.containerName)
                    }
                    instance.containerStatus = "running"
                }

                return instance, nil
            }
        }
    }

    // TODO create & start the container

    return nil, instance
}

func DockerContainerStop() (err error) {
    if cli == nil {
        return errors.New("Docker client is not initialized.")
    }

    return nil
}

func DockerContainerRemove() (err error) {
    if cli == nil {
        return errors.New("Docker client is not initialized.")
    }

    return nil
}

func DockerVolumeCreate() (err error) {
    if cli == nil {
        return errors.New("Docker client is not initialized.")
    }

    return nil
}

func DockerVolumeRemove() (err error) {
    if cli == nil {
        return errors.New("Docker client is not initialized.")
    }

    return nil
}

func DockerVolumeStatus() (err error) {
    if cli == nil {
        return errors.New("Docker client is not initialized.")
    }

    return nil
}
