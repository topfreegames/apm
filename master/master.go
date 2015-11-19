package master

import "github.com/topfreegames/apm/preparable"

type Master struct {
	SysFolder string	
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
	_, err := procPreparable.Start()
	if err != nil {
		return err
	}
	// TODO: Run watcher on this process.
	// TODO: Save process on process list.
	return nil
}
