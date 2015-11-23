package main

import "time"

import "github.com/topfreegames/apm/master"
import "github.com/topfreegames/apm/process"
import "github.com/topfreegames/apm/watcher"

func main() {
	master := &master.Master {
		SysFolder: "/Users/mdantas/testaAPM",
		Procs: make(map[string]*process.Proc),
		Watcher: watcher.InitWatcher(),
	}
	go master.WatchProcs()
	args := []string{"--prod", "--debug", "aguia_instance", "--app_name=colorfy", "--config_file=/Users/mdantas/go/src/git.topfreegames.com/topfreegames/aguia/lib/config.json", "--replicas=5"}
	preparable, err := master.Prepare(
		"git.topfreegames.com/topfreegames/aguia/lib",
		"colorfy",
		"go",
		true,
		args)
	if err != nil {
		panic(err)
	}
	err = master.RunPreparable(preparable)
	if err != nil {
		panic(err)
	}
	time.Sleep(60 * time.Second)
	master.StopProcess("colorfy")
	for{}
}
