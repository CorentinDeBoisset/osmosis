package clients

import (
    "fmt"
    "errors"
    "context"
    "strconv"
    "strings"
    dockerClient "github.com/docker/docker/client"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/network"
    "github.com/docker/docker/api/types/volume"
    "github.com/docker/docker/api/types/mount"
    "github.com/docker/docker/api/types/filters"
)

import "team-git.sancare.fr/dev/osmosis/cmd/tools"


type OsmosisDockerInstance struct {
    Id string
    Image string
    Name string
    Port int
    Status string
}

var cli *dockerClient.Client

func DockerConnect(verbose bool) (err error) {
    cli, err = dockerClient.NewClientWithOpts(dockerClient.FromEnv)
    if err != nil {
        return err
    }

    _, err = cli.Ping(context.Background())
    if err != nil {
        return errors.New("Could not connect to the docker daemon.\nCheck if it is running and the DOCKER_HOST env var in case the daemon is not listening to unix:///var/run/docker.sock")
    }

    return nil
}

func getContainerInfo(containerId string) (status string, listeningPort int, err error) {
    expandedInfo, err := cli.ContainerInspect(context.Background(), containerId)
    if err != nil {
        return "", -1, fmt.Errorf("Could not check status of container %s.", containerId)
    }

    if len(expandedInfo.NetworkSettings.NetworkSettingsBase.Ports) == 1 {
        // We have to iterate over PortBindings,
        // as we don't know what is the port in the container (it serves as key for the map)
        for _, portBindingList := range expandedInfo.NetworkSettings.NetworkSettingsBase.Ports {
            for _, portBinding := range portBindingList {
                listeningPort, err = strconv.Atoi(portBinding.HostPort)
                if err != nil {
                    return "", -1, fmt.Errorf("Could not read the port on which container %s is listening.", containerId)
                }
                return expandedInfo.State.Status, listeningPort, nil
            }
        }
    }

    return expandedInfo.State.Status, -1, nil
}

func GetDockerInstance(serviceName string, verbose bool) (instance *OsmosisDockerInstance, err error){
    if cli == nil {
        return nil, errors.New("Docker client is not initialized.")
    }

    existingInstances, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
    if err != nil {
        return nil, errors.New("Could not read list of containers")
    }

    for _, existingInstance := range existingInstances {
        if len(existingInstance.Names) > 0 && existingInstance.Names[0] == "/"+serviceName {
            status, portNb, err := getContainerInfo(existingInstance.ID)
            if err != nil {
                return nil, err
            }

            instance = &OsmosisDockerInstance{
                Id: existingInstance.ID,
                Image: existingInstance.Image,
                Name: existingInstance.Names[0],
                Port: portNb,
                Status: status,
            }
            return instance, nil
        }
    }

    return nil, nil
}

func DockerContainerStart(serviceName string, config tools.OsmosisServiceConfig, verbose bool) (instance *OsmosisDockerInstance, err error) {
    if cli == nil {
        return nil, errors.New("Docker client is not initialized.")
    }

    instance, err = GetDockerInstance(serviceName, verbose)
    if err != nil {
        return nil, err
    }

    ctx := context.Background()

    if instance != nil {
        if instance.Image != config.Image {
            return nil, fmt.Errorf("There is already a container named %s, but it is based on the image \"%s\".\nRun this command to remove any old containers:\n\n  osmosis clean", serviceName, instance.Image)
        }

        // If it is running or restarting, no problem
        if instance.Status != "running" && instance.Status != "restarting" {
            if instance.Status == "paused" {
                // If it was paused, we resume it
                err = cli.ContainerUnpause(ctx, instance.Id)
                if err != nil {
                    return nil, fmt.Errorf("Container %s is paused and could not be unpaused.", instance.Id)
                }
            } else {
                // In other cases, we start it
                err = cli.ContainerStart(ctx, instance.Id, types.ContainerStartOptions{})
                if err != nil {
                    return nil, fmt.Errorf("Container %s could not be started.", instance.Id)
                }
                status, portNb, err := getContainerInfo(instance.Id)
                if err != nil {
                    return nil, fmt.Errorf("Container %s could not be inspected.", instance.Id)
                }
                if status != "running" || portNb == -1 {
                    return nil, fmt.Errorf("Container %s was started but it could not be used.", serviceName)
                }
                instance.Port = portNb
                instance.Status = status
            }
        } else if instance.Port == -1 {
            return nil, fmt.Errorf("Container %s is running but not listening on any port.", instance.Id)
        }

        return instance, nil
    }

    containerConfig := container.Config{
        Image: config.Image,
        Hostname: serviceName,
        Env: []string{
            fmt.Sprintf("OWNER_UID=%s", config.UserId),
            fmt.Sprintf("OWNER_GID=%s", config.GroupId),
            fmt.Sprintf("IGNORED_DIRS=%s", strings.Join(config.Excludes, ":")),
        },
    }
    hostConfig := container.HostConfig{
        PublishAllPorts: true,
        Mounts: []mount.Mount{{ Type: "volume", Source: config.VolumeName, Target: "/sync" }},
    }
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

    status, portNb, err := getContainerInfo(createdContainer.ID)
    if err != nil {
        return nil, err
    }
    if status != "running" || portNb == -1 {
        return nil, fmt.Errorf("Container %s was created but it could not be used.", serviceName)
    }

    instance = &OsmosisDockerInstance{
        Id: createdContainer.ID,
        Image: config.Image,
        Name: serviceName,
        Port: portNb,
        Status: status,
    }

    return instance, nil
}

func DockerContainerStop(serviceName string, verbose bool) (err error) {
    if cli == nil {
        return errors.New("Docker client is not initialized.")
    }

    instance, err := GetDockerInstance(serviceName, verbose)
    if err != nil {
        return err
    }
    if instance == nil {
        return nil
    }

    err = cli.ContainerStop(context.Background(), instance.Id, nil)
    if err != nil {
        return fmt.Errorf("Container %s could not be stopped", instance.Id)
    }

    return nil
}

func DockerContainerRemove(serviceName string, verbose bool) (err error) {
    if cli == nil {
        return errors.New("Docker client is not initialized.")
    }

    instance, err := GetDockerInstance(serviceName, verbose)
    if err != nil {
        return err
    }
    if instance == nil {
        return nil
    }

    err = cli.ContainerRemove(context.Background(), instance.Id, types.ContainerRemoveOptions{})
    if err != nil {
        return fmt.Errorf("Container %s could not be removed", instance.Id)
    }

    return nil
}

func DockerVolumeCreate(volumeName string, verbose bool) (err error) {
    if cli == nil {
        return errors.New("Docker client is not initialized.")
    }

    // First, check if the volume already exists
    volumeExists, err := DockerVolumeStatus(volumeName, verbose)
    if err != nil {
        return err
    }
    if volumeExists {
        return nil
    }

    // If it does not exist, create it
    body := volume.VolumeCreateBody{
        Driver: "local",
        DriverOpts: nil,
        Labels: nil,
        Name: volumeName,
    }

    if _, err = cli.VolumeCreate(context.Background(), body); err != nil {
        return fmt.Errorf("Could not create volume %s.", volumeName)
    }

    return nil
}

func DockerVolumeRemove(volumeName string, verbose bool) (err error) {
    if cli == nil {
        return errors.New("Docker client is not initialized.")
    }

    // First, check if the volume already exists
    volumeExists, err := DockerVolumeStatus(volumeName, verbose)
    if err != nil {
        return err
    }
    if volumeExists {
        if err = cli.VolumeRemove(context.Background(), volumeName, false); err != nil {
            return fmt.Errorf("Could not remove docker volume %s.", volumeName)
        }
        return nil
    }

    return nil
}

func DockerVolumeStatus(volumeName string, verbose bool) (volumeExists bool, err error) {
    if cli == nil {
        return false, errors.New("Docker client is not initialized.")
    }

    volumeList, err := cli.VolumeList(context.Background(), filters.Args{})
    if err != nil {
        return false, errors.New("Could not read docker volume list")
    }

    for _, volumeInstance := range volumeList.Volumes {
        if volumeInstance.Name == volumeName {
            return true, nil
        }
    }

    return false, nil
}
