package work_conn_impl

import (
	"github.com/watermint/toolbox/domain/slack/api/work_auth"
	"github.com/watermint/toolbox/domain/slack/api/work_client"
	"github.com/watermint/toolbox/domain/slack/api/work_client_impl"
	"github.com/watermint/toolbox/domain/slack/api/work_conn"
	api_auth2 "github.com/watermint/toolbox/essentials/api/api_auth"
	"github.com/watermint/toolbox/essentials/api/api_conn"
	"github.com/watermint/toolbox/essentials/api/api_conn_impl"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/control/app_control"
)

func NewSlackApi(peerName string) work_conn.ConnSlackApi {
	return &connSlackApi{
		peerName: peerName,
		scopes:   make([]string, 0),
	}
}

type connSlackApi struct {
	peerName string
	client   work_client.Client
	scopes   []string
}

func (z *connSlackApi) Connect(ctl app_control.Control) (err error) {
	session := api_auth2.OAuthSessionData{
		AppData:  work_auth.Slack,
		PeerName: z.peerName,
		Scopes:   z.scopes,
	}
	entity, useMock, err := api_conn_impl.ConnectByRedirect(session, ctl)
	if useMock {
		z.client = work_client_impl.NewMock(z.peerName, ctl)
		return nil
	}
	if err != nil {
		return err
	}

	z.client = work_client_impl.New(z.peerName, ctl, entity)
	return nil
}

func (z *connSlackApi) PeerName() string {
	return z.peerName
}

func (z *connSlackApi) SetPeerName(name string) {
	z.peerName = name
}

func (z *connSlackApi) ScopeLabel() string {
	return app.ServiceSlack
}

func (z *connSlackApi) ServiceName() string {
	return api_conn.ServiceSlack
}

func (z *connSlackApi) SetScopes(scopes ...string) {
	z.scopes = scopes
}

func (z *connSlackApi) Scopes() []string {
	return z.scopes
}

func (z *connSlackApi) Context() work_client.Client {
	return z.client
}
