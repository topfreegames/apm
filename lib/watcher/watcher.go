package watcher

import "os"
import "sync"

import "github.com/topfreegames/apm/process"
import log "github.com/Sirupsen/logrus"

type ProcStatus struct {
	state *os.ProcessState
	err error
}

type ProcWatcher struct {
	procStatus chan *ProcStatus
	proc *process.Proc
	stopWatcher chan bool
}

type Watcher struct {
	sync.Mutex
	restartProc chan *process.Proc
	watchProcs map[string]*ProcWatcher
}

func InitWatcher() *Watcher {
	watcher := &Watcher {
		restartProc: make(chan *process.Proc),
		watchProcs: make(map[string]*ProcWatcher),
	}
	return watcher
}

func (watcher *Watcher) RestartProc() chan *process.Proc {
	return watcher.restartProc
}

func (watcher *Watcher) AddProcWatcher(proc *process.Proc) {
	watcher.Lock()
	defer watcher.Unlock()
	if _, ok := watcher.watchProcs[proc.Name]; ok {
		log.Warnf("A watcher for this process already exists.")
		return
	}
	procWatcher := &ProcWatcher{
		procStatus: make(chan *ProcStatus, 1),
		proc: proc,
		stopWatcher: make(chan bool, 1),
	}
	watcher.watchProcs[proc.Name] = procWatcher
	go func() {
		log.Infof("Starting watcher on proc %s", proc.Name)
		state, err := proc.Watch()
		procWatcher.procStatus <- &ProcStatus {
			state: state,
			err: err,
		}
	}()
	go func() {
		defer delete(watcher.watchProcs, procWatcher.proc.Name)
		select {
		case procStatus := <- procWatcher.procStatus:
			log.Infof("Proc %s is dead, advising master...", procWatcher.proc.Name)
			log.Infof("State is %s", procStatus.state.String())
			watcher.restartProc <- procWatcher.proc
			break
		case <- procWatcher.stopWatcher:
			break
		}
	}()
}

func (watcher *Watcher) StopWatcher(procName string) chan bool {
	if watcher, ok := watcher.watchProcs[procName]; ok {
		log.Infof("Stopping watcher on proc %s", procName)
		watcher.stopWatcher <- true
		waitStop := make(chan bool, 1)
		go func() {
			<- watcher.procStatus
			waitStop <- true
		}()
		return waitStop
	}
	return nil
}


