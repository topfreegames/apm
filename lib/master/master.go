/*
Master package is the main package that keeps everything running as it should be. It's responsible for starting, stopping and deleting processes. It also will keep an eye on the Watcher in case a process dies so it can restart it again.

RemoteMaster is responsible for exporting the main APM operations as HTTP requests. If you want to start a Remote Server, run:

- remoteServer := master.StartRemoteMasterServer(dsn, configFile)

It will start a remote master and return the instance.

To make remote requests, use the Remote Client by instantiating using:

- remoteClient, err := master.StartRemoteClient(dsn, timeout)

It will start the remote client and return the instance so you can use to initiate requests, such as:

- remoteClient.StartGoBin(sourcePath, name, keepAlive, args)
*/

package master

import "path"
import "errors"
import "fmt"
import "sync"

import "time"

import "github.com/topfreegames/apm/lib/preparable"
import "github.com/topfreegames/apm/lib/process"
import "github.com/topfreegames/apm/lib/utils"
import "github.com/topfreegames/apm/lib/watcher"

import log "github.com/Sirupsen/logrus"

// Master is the main module that keeps everything in place and execute
// the necessary actions to keep the process running as they should be.
type Master struct {
	sync.Mutex

	SysFolder string           // SysFolder is the main APM folder where the necessary config files will be stored.
	PidFile   string           // PidFille is the APM pid file path.
	OutFile   string           // OutFile is the APM output log file path.
	ErrFile   string           // ErrFile is the APM err log file path.
	Watcher   *watcher.Watcher // Watcher is a watcher instance.

	Procs map[string]*process.Proc // Procs is a map containing all procs started on APM.
}

// InitMaster will start a master instance with configFile.
// It returns a Master instance.
func InitMaster(configFile string) *Master {
	watcher := watcher.InitWatcher()
	master := &Master{}
	master.Procs = make(map[string]*process.Proc)
	err := utils.SafeReadTomlFile(configFile, master)
	if err != nil {
		panic(err)
	}
	master.Watcher = watcher
	master.Revive()
	log.Infof("All procs revived...")
	go master.WatchProcs()
	go master.SaveProcs()
	go master.UpdateStatus()
	return master
}

// WatchProcs will keep the procs running forever.
func (master *Master) WatchProcs() {
	for proc := range master.Watcher.RestartProc() {
		if !proc.KeepAlive {
			master.Lock()
			master.updateStatus(proc)
			master.Unlock()
			log.Infof("Proc %s does not have keep alive set. Will not be restarted.", proc.Name)
			continue
		}
		log.Infof("Restarting proc %s.", proc.Name)
		if proc.IsAlive() {
			log.Warnf("Proc %s was supposed to be dead, but it is alive.", proc.Name)
		}
		master.Lock()
		proc.Status.AddRestart()
		err := master.restart(proc)
		master.Unlock()
		if err != nil {
			log.Warnf("Could not restart process %s due to %s.", proc.Name, err)
		}
	}
}

// Prepare will compile the source code into a binary and return a preparable
// ready to be executed.
func (master *Master) Prepare(sourcePath string, name string, language string, keepAlive bool, args []string) (*preparable.ProcPreparable, []byte, error) {
	procPreparable := &preparable.ProcPreparable{
		Name:       name,
		SourcePath: sourcePath,
		SysFolder:  master.SysFolder,
		Language:   language,
		KeepAlive:  keepAlive,
		Args:       args,
	}
	output, err := procPreparable.PrepareBin()
	return procPreparable, output, err
}

// RunPreparable will run procPreparable and add it to the watch list in case everything goes well.
func (master *Master) RunPreparable(procPreparable *preparable.ProcPreparable) error {
	master.Lock()
	defer master.Unlock()
	if _, ok := master.Procs[procPreparable.Name]; ok {
		log.Warnf("Proc %s already exist.", procPreparable.Name)
		return errors.New("Trying to start a process that already exist.")
	}
	proc, err := procPreparable.Start()
	if err != nil {
		return err
	}
	master.Procs[proc.Name] = proc
	master.saveProcs()
	master.Watcher.AddProcWatcher(proc)
	proc.Status.SetStatus("running")
	return nil
}

// ListProcs will return a list of all procs.
func (master *Master) ListProcs() []*process.Proc {
	procsList := []*process.Proc{}
	for _, v := range master.Procs {
		procsList = append(procsList, v)
	}
	return procsList
}

// RestartProcess will restart a process.
func (master *Master) RestartProcess(name string) error {
	err := master.StopProcess(name)
	if err != nil {
		return err
	}
	return master.StartProcess(name)
}

// StartProcess will a start a process.
func (master *Master) StartProcess(name string) error {
	master.Lock()
	defer master.Unlock()
	if proc, ok := master.Procs[name]; ok {
		return master.start(proc)
	}
	return errors.New("Unknown process.")
}

// StopProcess will stop a process with the given name.
func (master *Master) StopProcess(name string) error {
	master.Lock()
	defer master.Unlock()
	if proc, ok := master.Procs[name]; ok {
		return master.stop(proc)
	}
	return errors.New("Unknown process.")
}

// DeleteProcess will delete a process and all its files and childs forever.
func (master *Master) DeleteProcess(name string) error {
	master.Lock()
	defer master.Unlock()
	log.Infof("Trying to delete proc %s", name)
	if proc, ok := master.Procs[name]; ok {
		err := master.stop(proc)
		if err != nil {
			return err
		}
		delete(master.Procs, name)
		err = master.delete(proc)
		if err != nil {
			return err
		}
		log.Infof("Successfully deleted proc %s", name)
	}
	return nil
}

// Revive will revive all procs listed on ListProcs. This should ONLY be called
// during Master startup.
func (master *Master) Revive() error {
	master.Lock()
	defer master.Unlock()
	procs := master.ListProcs()
	log.Info("Reviving all processes")
	for id := range procs {
		proc := procs[id]
		if !proc.KeepAlive {
			log.Infof("Proc %s does not have KeepAlive set. Will not revive it.", proc.Name)
			continue
		}
		log.Infof("Reviving proc %s", proc.Name)
		err := master.start(proc)
		if err != nil {
			return fmt.Errorf("Failed to revive proc %s due to %s", proc.Name, err)
		}
	}
	return nil
}

// NOT thread safe method. Lock should be acquire before calling it.
func (master *Master) start(proc *process.Proc) error {
	if !proc.IsAlive() {
		err := proc.Start()
		if err != nil {
			return err
		}
		master.Watcher.AddProcWatcher(proc)
		proc.Status.SetStatus("running")
	}
	return nil
}

func (master *Master) delete(proc *process.Proc) error {
	return proc.Delete()
}

// NOT thread safe method. Lock should be acquire before calling it.
func (master *Master) stop(proc *process.Proc) error {
	if proc.IsAlive() {
		waitStop := master.Watcher.StopWatcher(proc.Name)
		err := proc.GracefullyStop()
		if err != nil {
			return err
		}
		if waitStop != nil {
			<-waitStop
			proc.Pid = -1
			proc.Status.SetStatus("stopped")
		}
		log.Infof("Proc %s sucessfully stopped.", proc.Name)
	}
	return nil
}

// UpdateStatus will update a process status every 30s.
func (master *Master) UpdateStatus() {
	for {
		master.Lock()
		procs := master.ListProcs()
		for id := range procs {
			proc := procs[id]
			master.updateStatus(proc)
		}
		master.Unlock()
		time.Sleep(30 * time.Second)
	}
}

func (master *Master) updateStatus(proc *process.Proc) {
	if proc.IsAlive() {
		proc.Status.SetStatus("running")
	} else {
		proc.Pid = -1
		proc.Status.SetStatus("stopped")
	}
}

// NOT thread safe method. Lock should be acquire before calling it.
func (master *Master) restart(proc *process.Proc) error {
	err := master.stop(proc)
	if err != nil {
		return err
	}
	return master.start(proc)
}

// SaveProcs will save the list of procs onto the proc file.
func (master *Master) SaveProcs() {
	for {
		log.Infof("Saving list of procs.")
		master.Lock()
		master.saveProcs()
		master.Unlock()
		time.Sleep(5 * time.Minute)
	}
}

// Stop will stop APM and all of its running procs.
func (master *Master) Stop() error {
	log.Info("Stopping APM...")
	procs := master.ListProcs()
	for id := range procs {
		proc := procs[id]
		log.Info("Stopping proc %s", proc.Name)
		master.stop(proc)
	}
	log.Info("Saving and returning list of procs.")
	return master.saveProcs()
}

// NOT thread safe method. Lock should be acquire before calling it.
func (master *Master) saveProcs() error {
	configPath := master.getConfigPath()
	return utils.SafeWriteTomlFile(master, configPath)
}

func (master *Master) getConfigPath() string {
	return path.Join(master.SysFolder, "config.toml")
}
