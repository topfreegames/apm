package process

type ProcStatus struct {
	Status string
	Restarts int
}

func (proc_status *ProcStatus) SetStatus(status string) {
	proc_status.Status = status
}

func (proc_status *ProcStatus) AddRestart() {
	proc_status.Restarts += 1
}
