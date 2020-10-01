package provider

import (
	"main/users"

	"github.com/gorilla/sessions"
)

type Provider struct {
	HConf *HydraConfig
	Db    users.UserStore
	Store *sessions.CookieStore
}
