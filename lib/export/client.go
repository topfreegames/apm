package export

import "time"

import "git.apache.org/thrift.git/lib/go/thrift"
import "github.com/topfreegames/apm/lib/export/gen-go/apm"

func NewClient(addr string, timeout time.Duration) (*apm.ApmClient, error) {
	socketTransport, err := thrift.NewTSocket(addr)
	if err != nil {
		return nil, err
	}
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	
	transport := transportFactory.GetTransport(socketTransport)

	if err := transport.Open(); err != nil {
		return nil, err
	}
	
	return apm.NewApmClientFactory(transport, protocolFactory), nil
}
