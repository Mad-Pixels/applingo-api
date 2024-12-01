package api

import "github.com/Mad-Pixels/applingo-api/pkg/auth"

type ReqCtx struct {
	permissionLevel int
	authType        string
	user            string
}

func (c ReqCtx) HasPermissions(requiredLevel int) bool {
	return c.permissionLevel >= requiredLevel
}

func (c ReqCtx) IsDevice() bool {
	return c.authType == string(auth.HMAC)
}

func (c ReqCtx) IsUser() bool {
	return c.authType == string(auth.JWT)
}

func (c ReqCtx) GetUser() string {
	return c.user
}
