package build

import (
	bbuild "github.com/pip-services-samples/service-beacons-go/build"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
)

type ServiceFacadeFactory struct {
	cbuild.CompositeFactory
}

func NewServiceFacadeFactory() *ServiceFacadeFactory {
	c := &ServiceFacadeFactory{
		CompositeFactory: *cbuild.NewCompositeFactory(),
	}

	c.Add(bbuild.NewBeaconsServiceFactory())
	return c
}
