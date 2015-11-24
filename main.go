package main

import "gopkg.in/alecthomas/kingpin.v2"
import "github.com/topfreegames/apm/cli"
import "github.com/topfreegames/apm/master"

import "os"

var (
	app = kingpin.New("apm", "Aguia Process Manager.")
	dns = app.Flag("dns", "TCP Dns host.").Default(":9876").String()
	timeout = app.Flag("timeout", "Timeout to connect to client").Default("30s").Duration()

	serve = app.Command("serve", "Create APM server instance.")
	serveConfigFile = serve.Flag("config-file", "Config file location").Required().String()
	
	bin = app.Command("bin", "Create bin process.")
	binSourcePath = bin.Flag("source", "Go project source path. (Ex: github.com/topfreegames/apm)").Required().String()
	binName = bin.Arg("name", "Process name.").Required().String()
	binKeepAlive = bin.Flag("keep-alive", "Keep process alive forever.").Required().Bool()
	binArgs = bin.Flag("args", "External args.").Strings()

	start = app.Command("start", "Start a process.")
	startName = start.Arg("name", "Process name.").Required().String()

	stop = app.Command("stop", "Stop a process.")
	stopName = stop.Arg("name", "Process name.").Required().String()

	delete = app.Command("delete", "Delete a process.")
	deleteName = delete.Arg("name", "Process name.").Required().String()
	
	status = app.Command("status", "Get APM status.")
)
	

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case serve.FullCommand():
		master.StartRemoteMasterServer(*dns, *serveConfigFile)
	case bin.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.StartGoBin(*binSourcePath, *binName, *binKeepAlive, *binArgs)
	case start.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.StartProcess(*startName)
	case stop.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.StopProcess(*stopName)
	case delete.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.DeleteProcess(*deleteName)
	case status.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.Status()
	}
}
