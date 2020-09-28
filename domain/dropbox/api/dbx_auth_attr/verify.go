package dbx_auth_attr

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context_impl"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/control/app_control"
)

var (
	ErrorUnableToRetrieveResponse = errors.New("unable to retrieve json response")
	ErrorNoVerification           = errors.New("no verification")
	ErrorUnexpectedResponseFormat = errors.New("unexpected response format")
)

// Returns description of the account
func VerifyToken(name string, ctx api_auth.Context, ctl app_control.Control, app api_auth.App) (actx api_auth.Context, err error) {
	scopes := ctx.Scopes()
	l := ctl.Log().With(esl.String("peerName", ctx.PeerName()), esl.Strings("scopes", scopes))
	ui := ctl.UI()

	if len(scopes) != 1 {
		l.Debug("Skip verification")
		return ctx, nil
	}

	switch scopes[0] {
	case api_auth.DropboxTokenFull, api_auth.DropboxTokenApp:
		apiCtx := dbx_context_impl.New(name, ctl, ctx)
		res := apiCtx.Post("users/get_current_account")
		if err, fail := res.Failure(); fail {
			l.Debug("Unable to verify token", esl.Error(err))
			return nil, err
		}

		j := res.Success().Json()
		desc, found := j.FindString("name.display_name")
		if !found {
			return nil, ErrorUnexpectedResponseFormat
		}
		suppl, found := j.FindString("email")
		if !found {
			return nil, ErrorUnexpectedResponseFormat
		}
		l.Debug("Token Verified", esl.String("desc", desc))

		return api_auth.NewContextWithAttr(ctx, app.Config(scopes), desc, suppl), nil

	case api_auth.DropboxTokenBusinessInfo,
		api_auth.DropboxTokenBusinessManagement,
		api_auth.DropboxTokenBusinessFile,
		api_auth.DropboxTokenBusinessAudit:
		apiCtx := dbx_context_impl.New(name, ctl, ctx)
		res := apiCtx.Post("team/get_info")
		if err, fail := res.Failure(); fail {
			l.Debug("Unable to verify token", esl.Error(err))
			return nil, err
		}
		j := res.Success().Json()
		desc, found := j.FindString("name")
		if !found {
			return nil, ErrorUnexpectedResponseFormat
		}
		supplLic, found := j.FindNumber("num_licensed_users")
		if !found {
			return nil, ErrorUnexpectedResponseFormat
		}
		suppl := ui.Text(MAttr.AttrTeamLicenses.With("Licenses", supplLic))
		l.Debug("Token Verified", esl.String("desc", desc), esl.String("suppl", suppl))

		return api_auth.NewContextWithAttr(ctx, app.Config(scopes), desc, suppl), nil

	default:
		l.Debug("Skip verification")
		return ctx, nil
	}
}
