package gh_context

import "github.com/watermint/toolbox/infra/api/api_context"

type Context interface {
	api_context.Context
	api_context.Post
	api_context.Get
	api_context.Put
	api_context.Upload
	api_context.Patch
}
