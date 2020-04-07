package dbx_auth

import (
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/control/app_control"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"strings"
)

func NewConsoleGenerated(c app_control.Control, peerName string) api_auth.Console {
	return &Generated{
		ctl:      c,
		peerName: peerName,
	}
}

type Generated struct {
	ctl      app_control.Control
	peerName string
}

func (z *Generated) PeerName() string {
	return z.peerName
}

func (z *Generated) Auth(scope string) (tc api_auth.Context, err error) {
	token, err := z.generatedToken(scope)
	if err != nil {
		return nil, err
	}
	return api_auth.NewContext(token, z.peerName, scope), nil
}

func (z *Generated) generatedTokenInstruction(scope string) {
	ui := z.ctl.UI()
	api := ""
	toa := ""

	switch scope {
	case api_auth.DropboxTokenFull:
		api = "Dropbox API"
		toa = "Full Dropbox"
	case api_auth.DropboxTokenApp:
		api = "Dropbox API"
		toa = "App folder"
	case api_auth.DropboxTokenBusinessInfo:
		api = "Dropbox Business API"
		toa = "Team information"
	case api_auth.DropboxTokenBusinessAudit:
		api = "Dropbox Business API"
		toa = "Team auditing"
	case api_auth.DropboxTokenBusinessFile:
		api = "Dropbox Business API"
		toa = "Team member file access"
	case api_auth.DropboxTokenBusinessManagement:
		api = "Dropbox Business API"
		toa = "Team member management"
	default:
		z.ctl.Log().Fatal("Undefined token type", zap.String("type", scope))
	}

	ui.Info(MCcAuth.GeneratedToken1.With("API", api).With("TypeOfAccess", toa))
}

func (z *Generated) generatedToken(scope string) (*oauth2.Token, error) {
	ui := z.ctl.UI()
	z.generatedTokenInstruction(scope)
	for {
		code, cancel := ui.AskSecure(MCcAuth.GeneratedToken2)
		if cancel {
			return nil, ErrorUserCancelled
		}
		trim := strings.TrimSpace(code)
		if len(trim) > 0 {
			return &oauth2.Token{AccessToken: trim}, nil
		}
	}
}
