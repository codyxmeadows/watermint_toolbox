package fg_client

import "github.com/watermint/toolbox/essentials/api/api_client"

type Client interface {
	api_client.Client
	api_client.Get
	api_client.Post
	api_client.Delete
}
