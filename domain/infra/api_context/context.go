package api_context

import (
	"github.com/watermint/toolbox/app/app_ui"
	"github.com/watermint/toolbox/domain/infra/api_async"
	"github.com/watermint/toolbox/domain/infra/api_list"
	"github.com/watermint/toolbox/domain/infra/api_rpc"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Context interface {
	Log() *zap.Logger
	Msg(key string) app_ui.UIMessage
	ErrorMsg(err error) app_ui.UIMessage

	Request(endpoint string) api_rpc.Caller
	List(endpoint string) api_list.List
	Async(endpoint string) api_async.Async

	AsMemberId(teamMemberId string) Context
	AsAdminId(teamMemberId string) Context
	WithPath(pathRoot PathRoot) Context
	ResetLastErrors()
}

type RetryContext interface {
	AddError(err error)
	LastErrors() []error
	RetryAfter() time.Time
	UpdateRetryAfter(after time.Time)
}

type ClientContext interface {
	DoRequest(req api_rpc.Request) (code int, header http.Header, body []byte, err error)
}

type PathRoot interface {
	Header() string
}

func Home() PathRoot {
	return &homePathRoot{Tag: "home"}
}
func Root(namespaceId string) PathRoot {
	return &rootPathRoot{Tag: "root", Root: namespaceId}
}
func Namespace(namespaceId string) PathRoot {
	return &namespacePathRoot{Tag: "namespace_id", NamespaceId: namespaceId}
}

type homePathRoot struct {
	Tag string `json:".tag"`
}

func (*homePathRoot) Header() string {
	return "{\".tag\":\"home\"}"
}

type rootPathRoot struct {
	Tag  string `json:".tag"`
	Root string `json:"root"`
}

func (z rootPathRoot) Header() string {
	return "{\".tag\":\"root\",\"root\":\"" + z.Root + "\"}"
}

type namespacePathRoot struct {
	Tag         string `json:".tag"`
	NamespaceId string `json:"namespace_id"`
}

func (z namespacePathRoot) Header() string {
	return "{\".tag\":\"namespace_id\",\"namespace_id\":\"" + z.NamespaceId + "\"}"
}
