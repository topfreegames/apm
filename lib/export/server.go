package export

import "fmt"

import  "git.apache.org/thrift.git/lib/go/thrift"

import "github.com/topfreegames/apm/lib/master"
import "github.com/topfreegames/apm/lib/export/gen-go/apm"

type ApmHandler struct {
	master *master.Master
}

func (apmHandler *ApmHandler) Ping() (err error) {
	return nil
}

func (apmHandler *ApmHandler) Save() (err error) {
	return apmHandler.master.SaveProcs()
}

func (apmHandler *ApmHandler) Resurrect() (err error) {
	return apmHandler.master.Revive()
}

func (apmHandler *ApmHandler) Gobin(goBin *apm.GoBin) (str string, err error) {
	preparable, output, err := apmHandler.master.Prepare(goBin.SourcePath, goBin.Name, "go", goBin.KeepAlive, goBin.Args_)
	str = fmt.Sprintf("ERR: %s, OUT: %s", err, output)
	if err != nil {
		return str, err
	}
	return str, apmHandler.master.RunPreparable(preparable)
}

func (apmHandler *ApmHandler) StartProc(procName string) (err error) {
	return apmHandler.master.StartProcess(procName)
}

func (apmHandler *ApmHandler) StopProc(procName string) (err error) {
	return apmHandler.master.StopProcess(procName)
}

func (apmHandler *ApmHandler) RestartProc(procName string) (err error) {
	return apmHandler.master.RestartProcess(procName)
}

func (apmHandler *ApmHandler) DeleteProc(procName string) (err error) {
	return apmHandler.master.DeleteProcess(procName)
}

func (apmHandler *ApmHandler) Monit() (r []*apm.Proc, err error) {
	fmt.Println("OLA......")
	procs := apmHandler.master.ListProcs()
	r = []*apm.Proc{}
	for id := range procs {
		proc := procs[id]
		procConvert := &apm.Proc {
			Name: proc.Name,
			Cmd: proc.Cmd,
			Args_: proc.Args,
			Path: proc.Path,
			Pidfile: proc.Pidfile,
			Outfile: proc.Outfile,
			Errfile: proc.Errfile,
			KeepAlive: proc.KeepAlive,
			Pid: int32(proc.Pid),
			Status: &apm.ProcStatus {
				Status: proc.Status.Status,
				Restarts: int32(proc.Status.Restarts),
			},
		}
		r = append(r, procConvert)
	}
	return r, nil
}

func RunServer(addr string, configFile string) (*master.Master, error) {
	fmt.Printf("ADDR IS: %s\n", addr)
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	
	serverTransport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		return nil, err
	}

	masterHandler := master.InitMaster(configFile)
	handler := &ApmHandler{
		master: masterHandler,
	}
	processor := apm.NewApmProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	go func() {
		fmt.Printf("Serving...\n")
		err := server.Serve()
		fmt.Printf("ERROR IS: %s\n", err)
	}()
	
	return masterHandler, nil
}
