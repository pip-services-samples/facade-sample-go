package build

import (
	bbuild "github.com/pip-services-samples/pip-clients-beacons-go/build"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
)

type ClientFacadeFactory struct {
	cbuild.CompositeFactory
}

func NewClientFacadeFactory() *ClientFacadeFactory {
	c := &ClientFacadeFactory{
		CompositeFactory: *cbuild.NewCompositeFactory(),
	}

	// c.Add(accounts1.NewAccountsClientFactory())
	// c.Add(NewSessionsClientFactory())
	// c.Add(NewPasswordsClientFactory())
	// c.Add(NewRolesClientFactory())
	c.Add(bbuild.NewBeaconsClientFactory())
}
