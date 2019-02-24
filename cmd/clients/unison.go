package clients

import (
    "fmt"
    "os"
    "os/exec"
    "errors"
    "io/ioutil"
    "strconv"
    "syscall"
    "strings"
)

import "team-git.sancare.fr/dev/osmosis/cmd/tools"


type OsmosisUnisonInstance struct {
    Pid int
    Running bool
}

func GetUnisonInstance(serviceName string) (instance *OsmosisUnisonInstance) {
    if piddata, err := ioutil.ReadFile("/tmp/osmosis/"+serviceName+".pid"); err == nil {
        if pid, err := strconv.Atoi(string(piddata)); err == nil {
            if process, err := os.FindProcess(pid); err == nil {
                if err := process.Signal(syscall.Signal(0)); err == nil {
                    return &OsmosisUnisonInstance{Pid: pid, Running: true}
                }
            }
            return &OsmosisUnisonInstance{Pid: pid, Running: false}
        }
    }

    return &OsmosisUnisonInstance{Pid: -1, Running: false}
}


func StartUnisonInstance(serviceName string, config tools.OsmosisServiceConfig, listenPort int, verbose bool) (instance *OsmosisUnisonInstance, err error) {
    instance = GetUnisonInstance(serviceName)
    if instance.Pid != -1 && instance.Running {
        return nil, fmt.Errorf("A unison instance is already running for %s", serviceName)
    }

    instance = &OsmosisUnisonInstance{Pid: -1, Running: false}

    _, err = exec.LookPath("unison")
    if err != nil {
        return nil, errors.New("unison could not be found.\nCheck it is installed and in the $PATH before running osmosis.")
    }
    _, err = exec.LookPath("unison-fsmonitor")
    if err != nil {
        return nil, errors.New("unison-fsmonitor could not be found.\nCheck it is installed and in the $PATH before running osmosis.")
    }

    // Start a first unison command that is
    // aimed at initializing the archives.
    initUnisonCmd := buildUnisonCommand(config, listenPort)
    stderr, _ := initUnisonCmd.StderrPipe()

    if err = initUnisonCmd.Start(); err != nil {
        return nil, errors.New("Could not start initial sync with unison")
    }

    parsed_stderr, _ := ioutil.ReadAll(stderr)
    initUnisonCmd.Wait()

    if strings.Contains(string(parsed_stderr), "-ignorearchives") {
        // We are in an inconsistent state,
        // the unison archives must be rebuilt to make unison ok
        initUnisonCmd = buildUnisonCommand(config, listenPort)
        initUnisonCmd.Args = append(initUnisonCmd.Args, "-ignorearchives")
        if err = initUnisonCmd.Run(); err != nil {
            return nil, errors.New("An error occured trying syncing the two directories")
        }
    }

    watcherUnisonCmd := buildUnisonCommand(config, listenPort)
    watcherUnisonCmd.Args = append(watcherUnisonCmd.Args, "-silent", "-repeat", "watch")
    if err = watcherUnisonCmd.Start(); err != nil {
        return nil, fmt.Errorf("The two directories were synced, but the unison process failed to start, with the following error:\n  %s", err)
    }

    // Make directory /tmp/osmosis
    ioutil.WriteFile("/tmp/osmosis/"+serviceName+".pid", []byte(fmt.Sprintf("%d", watcherUnisonCmd.Process.Pid)), 0664)

    instance.Pid = watcherUnisonCmd.Process.Pid
    instance.Running = true

    if err = watcherUnisonCmd.Process.Release(); err != nil {
        return nil, errors.New("The unison process could not be detached")
    }

    return instance, nil
}


func StopUnisonInstance(serviceName string) (err error) {
    return nil
}


// Build an *exec.Cmd for unison with the apropriate properties
func buildUnisonCommand(config tools.OsmosisServiceConfig, listenPort int) (unisonCmd *exec.Cmd) {
    unisonCmd = exec.Command("unison")
    unisonCmd.Args = append(
        unisonCmd.Args,
        config.Src,
        fmt.Sprintf("socket://localhost:%d//sync", listenPort),
        "-batch",
        "-prefer",
        config.Src,
        "-copyonconflict",
    )
    for _, excludePath := range config.Excludes {
        unisonCmd.Args = append(unisonCmd.Args, "-ignore", fmt.Sprintf("BelowPath %s", excludePath))
    }
    unisonCmd.Env = os.Environ()

    return unisonCmd
}
