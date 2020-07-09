package provider

import (
	"main/users"
)

type Provider struct {
	HConf *HydraConfig
	Db    users.UserStore
}
