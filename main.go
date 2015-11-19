package main

import "github.com/topfreegames/apm/master"

func main() {
	master := master.Master {
		SysFolder: "/Users/mdantas/testaAPM",		
	}
	args := []string{"--prod", "--debug", "aguia_instance", "--app_name=colorfy", "--config_file=/Users/mdantas/go/src/git.topfreegames.com/topfreegames/aguia/lib/config.json", "--replicas=5", "--massive"}
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
}
