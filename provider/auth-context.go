package provider

import (
	"main/users"
)

type AuthContext struct {
	HConf *HydraConfig
	Db    users.UserStore
}
