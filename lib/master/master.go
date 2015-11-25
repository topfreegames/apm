package master

import "path"
import "errors"
import "sync"

import "time"

import "github.com/topfreegames/apm/preparable"
import "github.com/topfreegames/apm/process"
import "github.com/topfreegames/apm/utils"
import "github.com/topfreegames/apm/watcher"

import log "github.com/Sirupsen/logrus"

type Master struct {
	sync.Mutex
	SysFolder string	
	PidFile string
	OutFile string
	ErrFile string
	Watcher *watcher.Watcher
	
	Procs map[string]*process.Proc
}

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

// It will compile the source code into a binary and return a preparable
// ready to be executed.
func (master *Master) Prepare(sourcePath string, name string, language string, keepAlive bool, args []string) (*preparable.ProcPreparable, error) {
	procPreparable := &preparable.ProcPreparable {
		Name: name,
		SourcePath: sourcePath,
		SysFolder: master.SysFolder,
		Language: language,
		KeepAlive: keepAlive,
		Args: args,
	}
	_, err := procPreparable.PrepareBin()
	return procPreparable, err
}

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

func (master *Master) ListProcs() []*process.Proc {
	procsList := []*process.Proc{}
	for _, v := range master.Procs {
		procsList = append(procsList, v)
	}
	return procsList
}

func (master *Master) StartProcess(name string) error {
	master.Lock()
	defer master.Unlock()
	if proc, ok := master.Procs[name]; ok {
		return master.start(proc)
	}
	return errors.New("Unknown process.")
}

func (master *Master) StopProcess(name string) error {
	master.Lock()
	defer master.Unlock()
	if proc, ok := master.Procs[name]; ok {
		return master.stop(proc)
	}
	return errors.New("Unknown process.")
}

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

func (master *Master) Revive() {
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
			log.Infof("Failed to revive proc %s due to %s", proc.Name, err)
		}
	}
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
			<- waitStop
			proc.Pid = -1
			proc.Status.SetStatus("stopped")
		}
		log.Infof("Proc %s sucessfully stopped.", proc.Name)
	}
	return nil
}

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

func (master *Master) SaveProcs() {
	for {
		log.Infof("Saving list of procs.")
		master.Lock()
		master.saveProcs()
		master.Unlock()
		time.Sleep(5 * time.Minute)
	}
}

// NOT thread safe method. Lock should be acquire before calling it.
func (master *Master) saveProcs() error {
	configPath := master.getConfigPath()
	return utils.SafeWriteTomlFile(master, configPath)
}

func (master *Master) getConfigPath() string {
	return path.Join(master.SysFolder, "config.toml")
}
