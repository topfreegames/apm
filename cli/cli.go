package cli

import "github.com/topfreegames/apm/master"

import "log"
import "time"
import "fmt"

type Cli struct {
	remoteClient *master.RemoteClient
}

func InitCli(dsn string, timeout time.Duration) *Cli {
	client, err := master.StartRemoteClient(dsn, timeout)
	if err != nil {
		log.Fatal("Failed to start remote client due to: %+v\n", err)
	}
	return &Cli {
		remoteClient: client,
	}
}

func (cli *Cli) StartGoBin(sourcePath string, name string, keepAlive bool, args []string) {
	err := cli.remoteClient.StartGoBin(sourcePath, name, keepAlive, args)
	if err != nil {
		log.Fatal("Failed to start go bin due to: %+v\n", err)
	}
}

func (cli *Cli) StartProcess(procName string) {
	err := cli.remoteClient.StartProcess(procName)
	if err != nil {
		log.Fatal("Failed to start process due to: %+v\n", err)
	}
}

func (cli *Cli) StopProcess(procName string) {
	err := cli.remoteClient.StopProcess(procName)
	if err != nil {
		log.Fatal("Failed to stop process due to: %+v\n", err)
	}
}

func (cli *Cli) Status() {
	procs, err := cli.remoteClient.MonitStatus()
	if err != nil {
		log.Fatal("Failed to get status due to: %+v\n", err)
	}
	fmt.Printf("-----------------------------------------------------------------------------------\n")
	fmt.Printf("|     pid     |             name             |     status     |     keep-alive     |\n")  
	for id := range procs {
		proc := procs[id]
		kp := "True"
		if !proc.KeepAlive {
			kp = "False"
		}
		fmt.Printf("|%s|%s|%s|%s|\n",
			PadString(fmt.Sprintf("%d", proc.Pid), 13),
			PadString(proc.Name, 30),
			PadString(proc.Status.Status, 16),
			PadString(kp, 20))
	}
	fmt.Printf("-----------------------------------------------------------------------------------\n")
}

func PadString(str string, totalSize int) string {
	turn := 0
	for {
		if len(str) >= totalSize {
			return str
		}
		if turn == 0 {
			str = " " + str
			turn ^= 1
		} else {
			str = str + " "
			turn ^= 1
		}
	}
	return str
}
