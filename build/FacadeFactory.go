package build

import (
	service1 "github.com/pip-services-samples/pip-samples-facade-go/services/version1"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
)

type FacadeFactory struct {
	cbuild.Factory
	NullClientDescriptor      *cref.Descriptor
	DirectClientDescriptor    *cref.Descriptor
	HttpClientDescriptor      *cref.Descriptor
	GrpcClientDescriptor      *cref.Descriptor
	FacadeServiceV1Descriptor *cref.Descriptor
}

func NewFacadeFactory() *FacadeFactory {

	bcf := FacadeFactory{}
	bcf.Factory = *cbuild.NewFactory()

	bcf.FacadeServiceV1Descriptor = cref.NewDescriptor("pip-facades-example", "service", "http", "*", "1.0")

	bcf.RegisterType(bcf.FacadeServiceV1Descriptor, service1.NewFacadeServiceV1)

	return &bcf
}
