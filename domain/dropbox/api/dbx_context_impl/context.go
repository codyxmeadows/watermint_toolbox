package dbx_context_impl

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_async"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_async_impl"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_list"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_list_impl"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_request"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_response"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_response_impl"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/essentials/network/nw_client"
	"github.com/watermint/toolbox/essentials/network/nw_replay"
	"github.com/watermint/toolbox/essentials/network/nw_rest"
	"github.com/watermint/toolbox/essentials/network/nw_simulator"
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/api/api_request"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"net/http"
)

func NewMock(ctl app_control.Control) dbx_context.Context {
	client := nw_rest.New(
		nw_rest.Assert(dbx_response_impl.AssertResponse),
		nw_rest.Mock())
	return &ctxImpl{
		client:  client,
		ctl:     ctl,
		builder: dbx_request.NewBuilder(ctl, nil),
	}
}

func NewReplayMock(ctl app_control.Control, rr []nw_replay.Response) dbx_context.Context {
	client := nw_rest.New(
		nw_rest.Assert(dbx_response_impl.AssertResponse),
		nw_rest.ReplayMock(rr))
	return &ctxImpl{
		client:  client,
		ctl:     ctl,
		builder: dbx_request.NewBuilder(ctl, nil),
	}
}

func New(ctl app_control.Control, token api_auth.Context) dbx_context.Context {
	l := ctl.Log()
	opts := make([]nw_rest.ClientOpt, 0)
	opts = append(opts, nw_rest.Assert(dbx_response_impl.AssertResponse))
	if ctl.Feature().Experiment(app.ExperimentDbxClientConditionerNarrow20) {
		l.Debug("Experiment: Network conditioner enabled: 20%")
		opts = append(opts, nw_rest.Conditioner(20, nw_simulator.RetryAfterHeaderRetryAfter, decorateRateLimit))
	} else if ctl.Feature().Experiment(app.ExperimentDbxClientConditionerNarrow40) {
		l.Debug("Experiment: Network conditioner enabled: 40%")
		opts = append(opts, nw_rest.Conditioner(40, nw_simulator.RetryAfterHeaderRetryAfter, decorateRateLimit))
	} else if ctl.Feature().Experiment(app.ExperimentDbxClientConditionerNarrow100) {
		l.Debug("Experiment: Network conditioner enabled: 100%")
		opts = append(opts, nw_rest.Conditioner(100, nw_simulator.RetryAfterHeaderRetryAfter, decorateRateLimit))
	}
	client := nw_rest.New(opts...)
	return &ctxImpl{
		client:  client,
		ctl:     ctl,
		builder: dbx_request.NewBuilder(ctl, token),
	}
}

func decorateRateLimit(res *http.Response) {

}

type ctxImpl struct {
	client  nw_client.Rest
	ctl     app_control.Control
	builder dbx_request.Builder
}

func (z ctxImpl) UI() app_ui.UI {
	return z.ctl.UI()
}

func (z ctxImpl) ClientHash() string {
	return z.builder.ContentHash()
}

func (z ctxImpl) Log() esl.Logger {
	return z.builder.Log()
}

func (z ctxImpl) Capture() esl.Logger {
	return z.ctl.Capture()
}

func (z ctxImpl) Async(endpoint string, d ...api_request.RequestDatum) dbx_async.Async {
	return dbx_async_impl.New(&z, endpoint, d)
}

func (z ctxImpl) List(endpoint string, d ...api_request.RequestDatum) dbx_list.List {
	return dbx_list_impl.New(&z, endpoint, d)
}

func (z ctxImpl) Post(endpoint string, d ...api_request.RequestDatum) dbx_response.Response {
	b := z.builder.With(
		http.MethodPost,
		RpcRequestUrl(RpcEndpoint, endpoint),
		api_request.Combine(d),
	)
	return dbx_response_impl.New(z.client.Call(&z, b))
}

func (z ctxImpl) Upload(endpoint string, d ...api_request.RequestDatum) dbx_response.Response {
	b := z.builder.With(
		http.MethodPost,
		ContentRequestUrl(endpoint),
		api_request.Combine(d),
	)
	return dbx_response_impl.New(z.client.Call(&z, b))
}

func (z ctxImpl) Download(endpoint string, d ...api_request.RequestDatum) dbx_response.Response {
	b := z.builder.With(
		http.MethodPost,
		ContentRequestUrl(endpoint),
		api_request.Combine(d),
	)
	return dbx_response_impl.New(z.client.Call(&z, b))
}

func (z ctxImpl) Notify(endpoint string, d ...api_request.RequestDatum) dbx_response.Response {
	b := z.builder.With(
		http.MethodPost,
		RpcRequestUrl(NotifyEndpoint, endpoint),
		api_request.Combine(d),
	)
	return dbx_response_impl.New(z.client.Call(&z, b))
}

func (z ctxImpl) AsMemberId(teamMemberId string) dbx_context.Context {
	z.builder = z.builder.AsMemberId(teamMemberId)
	return z
}

func (z ctxImpl) AsAdminId(teamMemberId string) dbx_context.Context {
	z.builder = z.builder.AsAdminId(teamMemberId)
	return z
}

func (z ctxImpl) WithPath(pathRoot dbx_context.PathRoot) dbx_context.Context {
	z.builder = z.builder.WithPath(pathRoot)
	return z
}

func (z ctxImpl) NoAuth() dbx_context.Context {
	z.builder = z.builder.NoAuth()
	return z
}
