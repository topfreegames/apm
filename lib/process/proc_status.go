package process

// ProcStatus is a wrapper with the process current status.
type ProcStatus struct {
	Status   string
	Restarts int
}

// SetStatus will set the process string status.
func (proc_status *ProcStatus) SetStatus(status string) {
	proc_status.Status = status
}

// AddRestart will add one restart to the process status.
func (proc_status *ProcStatus) AddRestart() {
	proc_status.Restarts++
}
