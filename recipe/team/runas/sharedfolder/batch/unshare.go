package batch

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_sharedfolder"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_member"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_sharedfolder"
	"github.com/watermint/toolbox/domain/dropbox/usecase/uc_sharedfolder"
	"github.com/watermint/toolbox/essentials/go/es_lang"
	"github.com/watermint/toolbox/essentials/queue/eq_sequence"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/feed/fd_file"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
)

type Unshare struct {
	Peer                dbx_conn.ConnScopedTeam
	File                fd_file.RowFeed
	OperationLog        rp_model.TransactionReport
	LeaveCopy           bool
	SkipNotSharedFolder app_msg.Message
}

func (z *Unshare) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeFilesContentRead,
		dbx_auth.ScopeMembersRead,
		dbx_auth.ScopeSharingRead,
		dbx_auth.ScopeSharingWrite,
		dbx_auth.ScopeTeamDataMember,
	)
	z.File.SetModel(&MemberFolder{})
	z.OperationLog.SetModel(&MemberFolder{}, &mo_sharedfolder.SharedFolder{})
}

func (z *Unshare) unshare(mf *MemberFolder, svm sv_member.Member, c app_control.Control) error {
	member, err := svm.ResolveByEmail(mf.MemberEmail)
	if err != nil {
		z.OperationLog.Failure(err, mf)
		return err
	}

	cm := z.Peer.Client().AsMemberId(member.TeamMemberId)
	sf, err := uc_sharedfolder.NewResolver(cm).Resolve(mo_path.NewDropboxPath(mf.Path))
	switch err {
	case nil:
		// fall through
	case uc_sharedfolder.ErrorNotSharedFolder:
		z.OperationLog.Skip(z.SkipNotSharedFolder, mf)
		return nil

	default:
		z.OperationLog.Failure(err, mf)
		return err
	}

	err = sv_sharedfolder.New(cm).Remove(sf, sv_sharedfolder.LeaveACopy(z.LeaveCopy))
	if err != nil {
		z.OperationLog.Failure(err, sf)
		return err
	}
	z.OperationLog.Success(mf, sf)
	return nil
}

func (z *Unshare) Exec(c app_control.Control) error {
	if err := z.OperationLog.Open(); err != nil {
		return err
	}
	svm := sv_member.NewCached(z.Peer.Client())

	var lastErr, listErr error

	c.Sequence().Do(func(s eq_sequence.Stage) {
		s.Define("unshare", z.unshare, svm, c)
		q := s.Get("unshare")

		listErr = z.File.EachRow(func(m interface{}, rowIndex int) error {
			q.Enqueue(m)
			return nil
		})
	}, eq_sequence.ErrorHandler(func(err error, mouldId, batchId string, p interface{}) {
		lastErr = err
	}))

	return es_lang.NewMultiErrorOrNull(lastErr, listErr)
}

func (z *Unshare) Test(c app_control.Control) error {
	return qt_errors.ErrorHumanInteractionRequired
}
