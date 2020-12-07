package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/containerd/cgroups"
	"github.com/fsnotify/fsnotify"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/shirou/gopsutil/process"
	"github.com/wylswz/battery-saver/xmbs"
	"gopkg.in/yaml.v2"
)

const (
	path          = "xm-cg"
	defaultConfig = "/etc/xm/battery-saver.yaml"

	//CPU usage of CPU should be a float number
	//Denoting the percentage of CPU accessable from
	//Certain process
	CPU = "cpu"
	//Memory Valid units for memory usages are
	// GiB
	Memory = "memory"
	//Network is Not used currently
	Network = "network"
)

// PMap Mapping pid to corresponding cgroups
var PMap map[int32]interface{} = make(map[int32]interface{})

var CGroup *cgroups.Cgroup

func _findProcs(procRe string) []int32 {
	procs, err := process.Processes()
	pids := []int32{}
	if err != nil {
		log.Println(err)
	}
	for _, p := range procs {

		pName, err := p.Name()
		if err != nil {
			log.Println(err)
		}

		match, err := regexp.Match(
			procRe, []byte(pName),
		)
		if err != nil {
			log.Println(err)
		}
		if match {
			pids = append(pids, p.Pid)
		}
	}

	return pids
}

func _handleCPU(r xmbs.Rule) {
	procs := _findProcs(r.Process)
	enter := xmbs.Enter(xmbs.Keys(PMap), procs)
	exit := xmbs.Exit(xmbs.Keys(PMap), procs)

	log.Println("New processes found: ", enter)
	log.Println("Processes to remove: ", exit)

	amount := uint64(xmbs.ParseInt(r.Amount))

	for _, pidExit := range exit {
		cgExit, ok := PMap[pidExit]
		if ok {
			cgExit.(cgroups.Cgroup).Delete()
		}
	}

	var cpuPeriod uint64 = 100000
	var cpuQuota int64 = int64(cpuPeriod * amount / 100)

	for _, pidEnter := range enter {
		cgEnter, err := cgroups.New(
			cgroups.V1,
			cgroups.StaticPath(fmt.Sprintf("%v/%v", path, pidEnter)),
			&specs.LinuxResources{
				CPU: &specs.LinuxCPU{
					Quota:  &cpuQuota,
					Period: &cpuPeriod,
				},
			},
		)
		if err != nil {
			log.Println(err)
		} else {
			cgEnter.Add(cgroups.Process{Pid: int(pidEnter)})
		}
	}

}

func handleConfigEvent(e fsnotify.Event) {
	configPath := e.Name
	config := &xmbs.Config{}
	err := yaml.Unmarshal(xmbs.MustReadFileAsBytes(configPath), config)
	if err != nil {
		log.Println(err)
	}

	// Do something with config
	// Modifying cgroups
	for _, r := range config.Rules {
		switch r.Resource {
		case CPU:
			_handleCPU(r)
		}
	}
}

func monitorProcesses() {
	for k, v := range PMap {
		exists, err := process.PidExists(k)
		if err == nil && !exists {
			v.(cgroups.Cgroup).Delete()
			delete(PMap, k)
		}
	}
}

func main() {

	fs := flag.NewFlagSet("", flag.ExitOnError)
	configPath := fs.String("config", defaultConfig, "Config file for battery saver")
	fs.Parse(os.Args[1:])

	watcher, err := fsnotify.NewWatcher()
	xmbs.CheckErr(err)

	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Do something
				handleConfigEvent(event)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println(err)
			}
		}
	}()

	go func() {
		for {
			monitorProcesses()
			time.Sleep(time.Second * 10)
		}
	}()

	log.Println("Watching: ", *configPath)
	err = watcher.Add(*configPath)
	xmbs.CheckErr(err)
	<-done

	for _, v := range PMap {
		err := v.(cgroups.Cgroup).Delete()
		if err != nil {
			log.Println("Err deleting cgroups", err)
		}
	}
}
