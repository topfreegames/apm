package master

import "net"
import "net/rpc"
import "log"
import "time"

import "github.com/topfreegames/apm/process"

type RemoteMaster struct {
	master *Master
}

type RemoteClient struct {
	conn *rpc.Client
}

type GoBin struct {
	SourcePath string
	Name string
	KeepAlive bool
	Args []string
}

func (remote_master *RemoteMaster) StartGoBin(goBin *GoBin, ack *bool) error {
	preparable, err := remote_master.master.Prepare(goBin.SourcePath, goBin.Name, "go", goBin.KeepAlive, goBin.Args)
	*ack = true
	if err != nil {
		return err
	}
	return remote_master.master.RunPreparable(preparable)
}

func (remote_master *RemoteMaster) StartProcess(procName string, ack *bool) error {
	*ack = true
	return remote_master.master.StartProcess(procName)
}

func (remote_master *RemoteMaster) StopProcess(procName string, ack *bool) error {
	*ack = true
	return remote_master.master.StopProcess(procName)
}

func (remote_master *RemoteMaster) MonitStatus(req string, procs *[]*process.Proc) error {
	req = ""
	*procs = remote_master.master.ListProcs()
	return nil
}

func (remote_master *RemoteMaster) DeleteProcess(procName string, ack *bool) error {
	*ack = true
	return remote_master.master.DeleteProcess(procName)
}

func StartRemoteMasterServer(dsn string, configFile string) {
	remoteMaster := &RemoteMaster{
		master: InitMaster(configFile),
	}
	rpc.Register(remoteMaster)
	l, e := net.Listen("tcp", dsn)
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	rpc.Accept(l)
}

func StartRemoteClient(dsn string, timeout time.Duration) (*RemoteClient, error) {
	conn, err := net.DialTimeout("tcp", dsn, timeout)
	if err != nil {
		return nil, err
	}
	return &RemoteClient{conn: rpc.NewClient(conn)}, nil
}

func (client *RemoteClient) StartGoBin(sourcePath string, name string, keepAlive bool, args []string) error {
	goBin := &GoBin {
		SourcePath: sourcePath,
		Name: name,
		KeepAlive: keepAlive,
		Args: args,
	}
	var started bool
	return client.conn.Call("RemoteMaster.StartGoBin", goBin, &started)
}

func (client *RemoteClient) StartProcess(procName string) error {
	var started bool
	return client.conn.Call("RemoteMaster.StartProcess", procName, &started)
}

func (client *RemoteClient) StopProcess(procName string) error {
	var stopped bool
	return client.conn.Call("RemoteMaster.StopProcess", procName, &stopped)
}

func (client *RemoteClient) DeleteProcess(procName string) error {
	var deleted bool
	return client.conn.Call("RemoteMaster.DeleteProcess", procName, &deleted)
}

func (client *RemoteClient) MonitStatus() ([]*process.Proc, error) {
	var procs []*process.Proc
	err := client.conn.Call("RemoteMaster.MonitStatus", "", &procs)
	return procs, err
}



