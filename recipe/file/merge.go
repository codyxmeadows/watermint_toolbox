package file

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/usecase/uc_file_merge"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
)

type Merge struct {
	rc_recipe.RemarkIrreversible
	Peer                dbx_conn.ConnScopedIndividual
	From                mo_path.DropboxPath
	To                  mo_path.DropboxPath
	DryRun              bool
	KeepEmptyFolder     bool
	WithinSameNamespace bool
}

func (z *Merge) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeFilesContentRead,
		dbx_auth.ScopeFilesContentWrite,
	)
	z.DryRun = true
}

func (z *Merge) Exec(c app_control.Control) error {
	ctx := z.Peer.Context()

	ufm := uc_file_merge.New(ctx, c)
	opts := make([]uc_file_merge.MergeOpt, 0)
	if z.DryRun {
		opts = append(opts, uc_file_merge.DryRun())
	}
	if !z.KeepEmptyFolder {
		opts = append(opts, uc_file_merge.ClearEmptyFolder())
	}
	if z.WithinSameNamespace {
		opts = append(opts, uc_file_merge.WithinSameNamespace())
	}

	return ufm.Merge(z.From, z.To, opts...)
}

func (z *Merge) Test(c app_control.Control) error {
	err := rc_exec.ExecMock(c, &Merge{}, func(r rc_recipe.Recipe) {
		m := r.(*Merge)
		m.DryRun = true
		m.KeepEmptyFolder = true
		m.WithinSameNamespace = true
		m.From = qtr_endtoend.NewTestDropboxFolderPath("from")
		m.To = qtr_endtoend.NewTestDropboxFolderPath("to")
	})
	if err, _ = qt_errors.ErrorsForTest(c.Log(), err); err != nil {
		return err
	}
	err = rc_exec.ExecMock(c, &Merge{}, func(r rc_recipe.Recipe) {
		m := r.(*Merge)
		m.DryRun = false
		m.KeepEmptyFolder = false
		m.WithinSameNamespace = false
		m.From = qtr_endtoend.NewTestDropboxFolderPath("from")
		m.To = qtr_endtoend.NewTestDropboxFolderPath("to")
	})
	if err, _ = qt_errors.ErrorsForTest(c.Log(), err); err != nil {
		return err
	}
	return qt_errors.ErrorScenarioTest
}
