package container

import (
	ffactory "github.com/pip-services-samples/pip-samples-facade-go/build"
	cproc "github.com/pip-services3-go/pip-services3-container-go/container"
	mbuild "github.com/pip-services3-go/pip-services3-mongodb-go/build"
	rpcbuild "github.com/pip-services3-go/pip-services3-rpc-go/build"
)

type FacadeProcess struct {
	cproc.ProcessContainer
}

func NewBFacadeProcess() *FacadeProcess {

	bp := FacadeProcess{}
	bp.ProcessContainer = *cproc.NewEmptyProcessContainer()
	bp.AddFactory(ffactory.NewClientFacadeFactory())
	bp.AddFactory(ffactory.NewServiceFacadeFactory())
	bp.AddFactory(rpcbuild.NewDefaultRpcFactory())
	bp.AddFactory(mbuild.NewDefaultMongoDbFactory())

	return &bp
}
