package gh_context_impl

import (
	"crypto/sha256"
	"fmt"
	"github.com/watermint/toolbox/domain/github/api/gh_context"
	"github.com/watermint/toolbox/domain/github/api/gh_request"
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/infra/api/api_request"
	"github.com/watermint/toolbox/infra/api/api_response"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/network/nw_monitor"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"strings"
)

func NewNoAuth(ctl app_control.Control) gh_context.Context {
	return &Context{
		scope:     "",
		token:     nil,
		ctl:       ctl,
		isNoRetry: false,
	}
}

type Context struct {
	scope     string
	token     *oauth2.Token
	ctl       app_control.Control
	isNoRetry bool
}

func (z *Context) ClientHash() string {
	tok := ""
	if z.token != nil {
		tok = z.token.AccessToken
	}
	seeds := []string{
		"s", z.scope,
		"t", tok,
	}
	return fmt.Sprintf("%x", sha256.Sum224([]byte(strings.Join(seeds, ","))))
}

func (z *Context) Log() *zap.Logger {
	return z.ctl.Log()
}

func (z *Context) Capture() *zap.Logger {
	return z.ctl.Capture()
}

func (z *Context) NoRetryOnError() api_context.Context {
	return &Context{
		ctl:       z.ctl,
		isNoRetry: true,
	}
}

func (z *Context) IsNoRetry() bool {
	return z.isNoRetry
}

func (z *Context) MakeResponse(req *http.Request, res *http.Response) (api_response.Response, error) {
	return NewResponse(z, req, res)
}

func (z *Context) Post(endpoint string) api_request.Request {
	return gh_request.New(z, z.scope, z.token, endpoint, "POST")
}

func (z *Context) Get(endpoint string) api_request.Request {
	return gh_request.New(z, z.scope, z.token, endpoint, "GET")
}

func NewResponse(ctx api_context.Context, req *http.Request, res *http.Response) (api_response.Response, error) {
	l := ctx.Log()
	defer nw_monitor.Log(req, res)
	if res == nil {
		l.Debug("Null response")
		return nil, api_response.ErrorNoResponse
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		l.Debug("Unable to read body", zap.Error(err))
		return nil, err
	}
	res.ContentLength = int64(len(body))

	return api_response.New(res, body), nil
}
