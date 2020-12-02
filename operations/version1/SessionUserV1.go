package operations1

import (
	"time"
)

type SessionUserV1 struct {
	/* Identification */
	Id         string
	Login      string
	Name       string
	CreateTime time.Time

	/* User information */
	TimeZone string
	Language string
	Theme    string

	/* Security info **/
	Roles         []string
	ChangePwdTime time.Time
	Sites         Site
	Settings      interface{}

	/* Custom fields */
	CustomHdr interface{}
	CustomDat interface{}
}

type Site struct {
	Id   string
	name string
}
