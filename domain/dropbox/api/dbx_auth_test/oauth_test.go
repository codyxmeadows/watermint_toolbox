package dbx_auth_test

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/infra/api/api_auth_impl"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

func TestOAuth_Auth(t *testing.T) {
	qtr_endtoend.TestWithControl(t, func(ctl app_control.Control) {
		a := api_auth_impl.NewConsoleOAuth(ctl, "test-oauth-auth", dbx_auth.NewLegacyApp(ctl))
		if a.PeerName() != "test-oauth-auth" {
			t.Error(a.PeerName())
		}
		_, err := a.Start([]string{"test-scope"})
		if err != app.ErrorUserCancelled {
			t.Error(err)
		}
	})
}
