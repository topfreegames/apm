package cli

import "github.com/topfreegames/apm/lib/export"
import "github.com/topfreegames/apm/lib/export/gen-go/apm"

import "log"
import "time"
import "fmt"

// Cli is the command line client.
type Cli struct {
	remoteClient *apm.ApmClient
}

// InitCli initiates a remote client connecting to dsn.
// Returns a Cli instance.
func InitCli(dsn string, timeout time.Duration) *Cli {
	client, err := export.NewClient(dsn, timeout)
	if err != nil {
		client.Transport.Close()
		log.Fatalf("Failed to start remote client due to: %+v\n", err)
	}
	return &Cli{
		remoteClient: client,
	}
}

func (cli *Cli) Close() {
	cli.remoteClient.Transport.Close()
}

// Save will save all previously saved processes onto a list.
// Display an error in case there's any.
func (cli *Cli) Save() {
	err := cli.remoteClient.Save()
	if err != nil {
		log.Printf("Failed to save list of processes due to: %+v\n", err)
	}
}

// Resurrect will restore all previously save processes.
// Display an error in case there's any.
func (cli *Cli) Resurrect() {
	err := cli.remoteClient.Resurrect()
	if err != nil {
		log.Printf("Failed to resurrect all previously save processes due to: %+v\n", err)
	}
}
// StartGoBin will try to start a go binary process.
// Returns a fatal error in case there's any.
func (cli *Cli) StartGoBin(sourcePath string, name string, keepAlive bool, args []string) {
	goBin := &apm.GoBin {
		SourcePath: sourcePath,
		Name: name,
		KeepAlive: keepAlive,
		Args_: args,
	}
	out, err := cli.remoteClient.Gobin(goBin)
	if err != nil {
		log.Printf("Failed to start go bin due to: %+v -> %s\n", err, out)
	}
}

// RestartProcess will try to restart a process with procName. Note that this process
// must have been already started through StartGoBin.
func (cli *Cli) RestartProcess(procName string) {
	err := cli.remoteClient.RestartProc(procName)
	if err != nil {
		log.Printf("Failed to restart process due to: %+v\n", err)
	}
}

// StartProcess will try to start a process with procName. Note that this process
// must have been already started through StartGoBin.
func (cli *Cli) StartProcess(procName string) {
	err := cli.remoteClient.StartProc(procName)
	if err != nil {
		log.Printf("Failed to start process due to: %+v\n", err)
	}
}

// StopProcess will try to stop a process named procName.
func (cli *Cli) StopProcess(procName string) {
	err := cli.remoteClient.StopProc(procName)
	if err != nil {
		log.Printf("Failed to stop process due to: %+v\n", err)
	}
}

// DeleteProcess will stop and delete all dependencies from process procName forever.
func (cli *Cli) DeleteProcess(procName string) {
	err := cli.remoteClient.DeleteProc(procName)
	if err != nil {
		log.Printf("Failed to delete process due to: %+v\n", err)
	}
}

// Status will display the status of all procs started through StartGoBin.
func (cli *Cli) Status() {
	procs, err := cli.remoteClient.Monit()
	if err != nil {
		log.Printf("Failed to get status due to: %+v\n", err)
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

// PadString will add totalSize spaces evenly to the right and left side of str.
// Returns str after applying the pad.
func PadString(str string, totalSize int) string {
	turn := 0
	for {
		if len(str) >= totalSize {
			break
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
