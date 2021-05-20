package file

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_file"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_file"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/report/rp_writer_impl"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
)

type Watch struct {
	Peer      dbx_conn.ConnScopedIndividual
	Path      mo_path.DropboxPath
	Recursive bool
}

func (z *Watch) Exec(c app_control.Control) error {
	ctx := z.Peer.Context()
	opts := make([]sv_file.ListOpt, 0)
	opts = append(opts, sv_file.Recursive(z.Recursive))
	w := rp_writer_impl.NewJsonWriter("entries", c, true)
	if err := w.Open(c, &mo_file.ConcreteEntry{}); err != nil {
		return err
	}
	defer w.Close()

	return sv_file.NewFiles(ctx).Poll(z.Path, func(entry mo_file.Entry) {
		w.Row(entry.Concrete())
	}, opts...)
}

func (z *Watch) Test(c app_control.Control) error {
	return qt_errors.ErrorNoTestRequired
}

func (z *Watch) Preset() {
	z.Peer.SetScopes(dbx_auth.ScopeFilesContentRead)
}
