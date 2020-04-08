package dbx_auth_attr

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/control/app_control"
	"go.uber.org/zap"
)

func NewConsole(c app_control.Control, peerName string) api_auth.Console {
	l := c.Log().With(zap.String("peerName", peerName))
	var oa api_auth.Console

	// Make redirect impl. hidden for while
	if c.Feature().IsConfigEnabled(app_control.ConfigKeyDomainDropboxAuthUseRedirect) {
		oa = dbx_auth.NewConsoleRedirect(c, peerName)
	} else {
		oa = dbx_auth.NewConsoleOAuth(c, peerName)
	}
	aa := NewConsoleAttr(c, oa)
	if c.Feature().IsSecure() {
		l.Debug("Skip caching")
		return aa
	}
	l.Debug("Token cache enabled")
	ca := dbx_auth.NewConsoleCache(c, aa)
	return ca
}

func NewConsoleAttr(c app_control.Control, auth api_auth.Console) api_auth.Console {
	return &Attr{
		ctl:  c,
		auth: auth,
	}
}
