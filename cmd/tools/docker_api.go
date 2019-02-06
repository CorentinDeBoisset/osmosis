package tools

import "fmt"
import "errors"
import "context"
import "strconv"
import "github.com/docker/docker/client"
import "github.com/docker/docker/api/types"
import "github.com/docker/docker/api/types/container"
import "github.com/docker/docker/api/types/network"


type OsmosisDockerInstance struct {
    Id string
    Name string
    Port int
}

var cli *client.Client

func DockerConnect(verbose bool) (err error) {
    cli, err = client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        return err
    }

    _, err = cli.Ping(context.Background())
    if err != nil {
        return errors.New("Could not connect to the docker daemon.\nCheck if it is running and the DOCKER_HOST env var in case the daemon is not listening to unix:///var/run/docker.sock")
    }

    return nil
}

func GetDockerInstances(verbose bool) (instances []OsmosisDockerInstance, err error){
    if cli == nil {
        return nil, errors.New("Docker client is not initialized.")
    }

    instances = make([]OsmosisDockerInstance, 0)

    existingInstances, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
    if err != nil {
        return nil, errors.New("Could not read list of containers")
    }

    for _, existingInstance := range existingInstances {
        if existingInstance.Image == "registry.sancare.fr/base_images/unison:1.0" && len(existingInstance.Names) > 0 {
            instances = append(instances, OsmosisDockerInstance{Id: existingInstance.ID, Name: existingInstance.Names[0], Port: -1})
        }
    }

    return instances, nil
}

func DockerContainerStart(serviceName string, config OsmosisServiceConfig, verbose bool) (instance *OsmosisDockerInstance, err error) {
    if cli == nil {
        return nil, errors.New("Docker client is not initialized.")
    }

    existingInstances, err := GetDockerInstances(verbose)
    if err != nil {
        return nil, err
    }

    ctx := context.Background()

    if len(existingInstances) > 0 {
        for _, existingInstance := range existingInstances {
            if existingInstance.Name == "/"+serviceName {
                // We first inspect the container
                expandedInfo, err := cli.ContainerInspect(ctx, existingInstance.Id)
                if err != nil {
                    return nil, fmt.Errorf("Could not check status of container %s.", existingInstance.Id)
                }

                // If it is running or restarting, no problem
                if !expandedInfo.State.Running && !expandedInfo.State.Restarting {
                    if expandedInfo.State.Paused {
                        // If it was paused, we resume it
                        err = cli.ContainerUnpause(ctx, existingInstance.Id)
                        if err != nil {
                            return nil, fmt.Errorf("Container %s is paused and could not be unpaused.", existingInstance.Id)
                        }
                    } else {
                        // In other cases, we start it
                        err = cli.ContainerStart(ctx, existingInstance.Id, types.ContainerStartOptions{})
                        if err != nil {
                            return nil, fmt.Errorf("Container %s could not be started.", existingInstance.Id)
                        }
                    }
                }
                err = cli.ContainerStart(ctx, existingInstance.Id, types.ContainerStartOptions{})

                return &existingInstance, nil
            }
        }
    }

    // The container does not exist, we create and start it
    containerConfig := container.Config{Image: config.Image, Hostname: serviceName}
    hostConfig := container.HostConfig{PublishAllPorts: true}
    networkConfig := network.NetworkingConfig{}
    createdContainer, err := cli.ContainerCreate(ctx, &containerConfig, &hostConfig, &networkConfig, serviceName)
    if err != nil {
        if verbose {
            return nil, fmt.Errorf("Creation of container %s failed with the following error:\n  %s", serviceName, err)
        } else {
            return nil, fmt.Errorf("Creation of container %s failed.", serviceName)
        }
    }

    err = cli.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{})
    if err != nil {
        return nil, fmt.Errorf("Container %s was created but could not be started.", serviceName)
    }

    runningContainer, err := cli.ContainerInspect(ctx, createdContainer.ID)

    instance = &OsmosisDockerInstance{Id: createdContainer.ID, Name: serviceName, Port: -1}
    if len(runningContainer.NetworkSettings.NetworkSettingsBase.Ports) == 1 {
        // We have to iterate over PortBindings,
        // as we don't know what is the port in the container (it serves as key for the map)
        for _, portBindingList := range runningContainer.NetworkSettings.NetworkSettingsBase.Ports {
            for _, portBinding := range portBindingList {
                instance.Port, err = strconv.Atoi(portBinding.HostPort)
                if err != nil {
                    return nil, fmt.Errorf("Could not get on which port number, the container %s is listening to.", serviceName)
                }
            }
        }
    } else {
        return nil, fmt.Errorf("Container %s is not listening on only one port", serviceName)
    }

    return instance, nil
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
