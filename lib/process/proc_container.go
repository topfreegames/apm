package process

import "os"
import "syscall"
import "errors"
import "strconv"

import "github.com/topfreegames/apm/lib/utils"

type ProcContainer interface {
	Start() error
	ForceStop() error
	GracefullyStop() error
	Restart() error
	Delete() error
	IsAlive() bool
	Identifier() string
	ShouldKeepAlive() bool
	AddRestart()
	NotifyStopped()
	SetStatus(status string)
	GetPid() int
	GetStatus() *ProcStatus
	Watch() (*os.ProcessState, error)
	release()
}

// Proc is a os.Process wrapper with Status and more info that will be used on Master to maintain
// the process health.
type Proc struct {
	Name      string
	Cmd       string
	Args      []string
	Path      string
	Pidfile   string
	Outfile   string
	Errfile   string
	KeepAlive bool
	Pid       int
	Status    *ProcStatus
	process   *os.Process
}

// Start will execute the command Cmd that should run the process. It will also create an out, err and pidfile
// in case they do not exist yet.
// Returns an error in case there's any.
func (proc *Proc) Start() error {
	outFile, err := utils.GetFile(proc.Outfile)
	if err != nil {
		return err
	}
	errFile, err := utils.GetFile(proc.Errfile)
	if err != nil {
		return err
	}
	wd, _ := os.Getwd()
	procAtr := &os.ProcAttr{
		Dir: wd,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			outFile,
			errFile,
		},
	}
	args := append([]string{proc.Name}, proc.Args...)
	process, err := os.StartProcess(proc.Cmd, args, procAtr)
	if err != nil {
		return err
	}
	proc.process = process
	proc.Pid = proc.process.Pid
	err = utils.WriteFile(proc.Pidfile, []byte(strconv.Itoa(proc.process.Pid)))
	if err != nil {
		return err
	}

	proc.Status.SetStatus("started")
	return nil
}

// ForceStop will forcefully send a SIGKILL signal to process killing it instantly.
// Returns an error in case there's any.
func (proc *Proc) ForceStop() error {
	if proc.process != nil {
		err := proc.process.Signal(syscall.SIGKILL)
		proc.Status.SetStatus("stopped")
		proc.release()
		return err
	}
	return errors.New("Process does not exist.")
}

// GracefullyStop will send a SIGTERM signal asking the process to terminate.
// The process may choose to die gracefully or ignore this signal completely. In that case
// the process will keep running unless you call ForceStop()
// Returns an error in case there's any.
func (proc *Proc) GracefullyStop() error {
	if proc.process != nil {
		err := proc.process.Signal(syscall.SIGTERM)
		proc.Status.SetStatus("asked to stop")
		return err
	}
	return errors.New("Process does not exist.")
}

// Restart will try to gracefully stop the process and then Start it again.
// Returns an error in case there's any.
func (proc *Proc) Restart() error {
	if proc.IsAlive() {
		err := proc.GracefullyStop()
		if err != nil {
			return err
		}
	}
	return proc.Start()
}

// Delete will delete everything created by this process, including the out, err and pid file.
// Returns an error in case there's any.
func (proc *Proc) Delete() error {
	proc.release()
	err := utils.DeleteFile(proc.Outfile)
	if err != nil {
		return err
	}
	err = utils.DeleteFile(proc.Errfile)
	if err != nil {
		return err
	}
	return os.RemoveAll(proc.Path)
}

// IsAlive will check if the process is alive or not.
// Returns true if the process is alive or false otherwise.
func (proc *Proc) IsAlive() bool {
	p, err := os.FindProcess(proc.Pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}

// Watch will stop execution and wait until the process change its state. Usually changing state, means that the process died.
// Returns a tuple with the new process state and an error in case there's any.
func (proc *Proc) Watch() (*os.ProcessState, error) {
	return proc.process.Wait()
}

// Will release the process and remove its PID file
func (proc *Proc) release() {
	if proc.process != nil {
		proc.process.Release()
	}
	utils.DeleteFile(proc.Pidfile)
}

// Notify that process was stopped so we can set its PID to -1
func (proc *Proc) NotifyStopped() {
	proc.Pid = -1;
}

// Add one restart to proc status
func (proc *Proc) AddRestart() {
	proc.Status.AddRestart()
}

// Return proc current PID
func (proc *Proc) GetPid() int {
	return proc.Pid;
}

// Return proc current status
func (proc *Proc) GetStatus() *ProcStatus {
	return proc.Status;
}

// Set proc status
func (proc *Proc) SetStatus(status string) {
	proc.Status.SetStatus(status);
}

// Proc identifier that will be used by watcher to keep track of its processes
func (proc *Proc) Identifier() string {
	return proc.Name;
}

// Returns true if the process should be kept alive or not
func (proc *Proc) ShouldKeepAlive() bool {
	return proc.KeepAlive;
}
