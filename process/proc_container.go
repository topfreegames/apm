package process

import "os"
import "syscall"
import "errors"
import "strconv"

import "github.com/topfreegames/apm/utils"

type Proc struct {
	Name string
	Cmd string
	Args []string
	Path string
	Pidfile string
	Outfile string
	Errfile string
	KeepAlive bool
	status *ProcStatus
	process *os.Process
}

func (proc *Proc) Start() error {
	proc.status = &ProcStatus{}
	outFile, err := utils.GetFile(proc.Outfile)
	if err != nil {
		return err
	}
	errFile, err := utils.GetFile(proc.Errfile)
	if err != nil {
		return err
	}
	wd, _ := os.Getwd()
	procAtr := &os.ProcAttr {
		Dir: wd,
		Env: os.Environ(),
		Files: []*os.File {
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
	err = utils.WriteFile(proc.Pidfile, []byte(strconv.Itoa(proc.process.Pid)))
	if err != nil {
		return err
	}

	proc.status.SetStatus("started")
	return nil
}

// Will forcefully kill the process
func (proc *Proc) ForceStop() error {
	if proc.process != nil {
		err := proc.process.Signal(syscall.SIGKILL)
		proc.status.SetStatus("stopped")
		proc.release()
		return err
	}
	return errors.New("Process does not exist.")
}

// Will send a SIGTERM signal asking the process
// to terminate. Note that the process may choose to ignore it.
func (proc *Proc) GracefullyStop() error {
	if proc.process != nil {
		err := proc.process.Signal(syscall.SIGTERM)
		proc.status.SetStatus("asked to stop")
		return err
	}
	return errors.New("Process does not exist.")
}

func (proc *Proc) Restart() error {
	if proc.IsAlive() {
		err := proc.GracefullyStop()
		if err != nil {
			return err
		}
	}
	return proc.Start()
}

func (proc *Proc) IsAlive() bool {
	return proc.process.Signal(syscall.Signal(0)) == nil
}

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

