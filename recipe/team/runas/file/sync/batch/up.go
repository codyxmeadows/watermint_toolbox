package batch

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_file_filter"
	mo_path2 "github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_member"
	"github.com/watermint/toolbox/essentials/model/mo_filter"
	"github.com/watermint/toolbox/essentials/model/mo_int"
	"github.com/watermint/toolbox/essentials/model/mo_path"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/feed/fd_file"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/ingredient/file"
	"github.com/watermint/toolbox/quality/infra/qt_file"
	"os"
)

type Up struct {
	rc_recipe.RemarkIrreversible
	Peer          dbx_conn.ConnScopedTeam
	File          fd_file.RowFeed
	Upload        *file.Upload
	OperationLog  rp_model.TransactionReport
	Overwrite     bool
	Delete        bool
	BatchSize     mo_int.RangeInt
	ExitOnFailure bool
	Name          mo_filter.Filter
}

func (z *Up) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeFilesContentRead,
		dbx_auth.ScopeFilesContentWrite,
		dbx_auth.ScopeMembersRead,
		dbx_auth.ScopeTeamDataMember,
	)
	z.BatchSize.SetRange(1, 1000, 250)
	z.File.SetModel(&UpMapping{})
	z.OperationLog.SetModel(&UpMapping{}, nil)
	z.Name.SetOptions(
		mo_filter.NewNameFilter(),
		mo_filter.NewNameSuffixFilter(),
		mo_filter.NewNamePrefixFilter(),
		mo_file_filter.NewIgnoreFileFilter(),
	)
}

func (z *Up) Exec(c app_control.Control) error {
	if err := z.OperationLog.Open(); err != nil {
		return err
	}

	svm := sv_member.NewCached(z.Peer.Context())

	return z.File.EachRow(func(m interface{}, rowIndex int) error {
		um := m.(*UpMapping)
		member, err := svm.ResolveByEmail(um.MemberEmail)
		if err != nil {
			z.OperationLog.Failure(err, um)
			if z.ExitOnFailure {
				return err
			}
			return nil
		}

		err = rc_exec.Exec(c, z.Upload, func(r rc_recipe.Recipe) {
			ru := r.(*file.Upload)
			ru.LocalPath = mo_path.NewFileSystemPath(um.LocalPath)
			ru.DropboxPath = mo_path2.NewDropboxPath(um.DropboxPath)
			ru.Overwrite = z.Overwrite
			ru.Name = z.Name
			ru.Context = z.Peer.Context().AsMemberId(member.TeamMemberId)
			ru.BatchSize = z.BatchSize.Value()
			ru.Delete = z.Delete
		})
		if err != nil {
			z.OperationLog.Failure(err, um)
		} else {
			z.OperationLog.Success(um, nil)
		}
		return err
	})
}

func (z *Up) Test(c app_control.Control) error {
	f, err := qt_file.MakeTestFile("share", "john@example.com,/local/john,/file_server\nemma@example.com,/local/emma,/file_server")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(f)
	}()

	return rc_exec.ExecMock(c, &Up{}, func(r rc_recipe.Recipe) {
		m := r.(*Up)
		m.File.SetFilePath(f)
	})
}
