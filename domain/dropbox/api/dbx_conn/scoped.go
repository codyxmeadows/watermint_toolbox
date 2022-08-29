package dbx_conn

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_client"
	"github.com/watermint/toolbox/infra/api/api_conn"
)

type ConnScopedDropboxApi interface {
	api_conn.ScopedConnection

	Context() dbx_client.Client
}

type ConnScopedTeam interface {
	ConnScopedDropboxApi
	IsTeam() bool
}

type ConnScopedIndividual interface {
	ConnScopedDropboxApi
	IsIndividual() bool
}
