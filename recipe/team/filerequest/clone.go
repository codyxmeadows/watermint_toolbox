package filerequest

import (
	"errors"
	"github.com/watermint/toolbox/domain/model/mo_filerequest"
	"github.com/watermint/toolbox/domain/model/mo_member"
	"github.com/watermint/toolbox/domain/model/mo_path"
	"github.com/watermint/toolbox/domain/service/sv_filerequest"
	"github.com/watermint/toolbox/domain/service/sv_member"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/feed/fd_file"
	"github.com/watermint/toolbox/infra/recipe/rc_conn"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/quality/infra/qt_endtoend"
	"strings"
)

type Clone struct {
	File         fd_file.RowFeed
	Peer         rc_conn.ConnBusinessFile
	OperationLog rp_model.TransactionReport
}

func (z *Clone) Preset() {
	z.File.SetModel(&mo_filerequest.MemberFileRequest{})
	z.OperationLog.SetModel(&mo_filerequest.MemberFileRequest{}, &mo_filerequest.MemberFileRequest{})
}

func (z *Clone) Hidden() {
}

func (z *Clone) Exec(c app_control.Control) error {
	members, err := sv_member.New(z.Peer.Context()).List()
	if err != nil {
		return err
	}
	emailToMember := mo_member.MapByEmail(members)

	if err := z.OperationLog.Open(); err != nil {
		return err
	}

	return z.File.EachRow(func(m interface{}, rowIndex int) error {
		fm := m.(*mo_filerequest.MemberFileRequest)
		if fm.Email == "" || fm.Destination == "" || fm.Title == "" {
			z.OperationLog.Failure(errors.New("invalid data"), fm)
			return nil
		}
		member, ok := emailToMember[strings.ToLower(fm.Email)]
		if !ok {
			z.OperationLog.Failure(errors.New("entry not found for the id"), fm)
			return nil
		}

		opts := make([]sv_filerequest.UpdateOpt, 0)
		if fm.Deadline != "" {
			opts = append(opts, sv_filerequest.OptDeadline(fm.Deadline))
		}
		if fm.DeadlineAllowLateUploads != "" {
			opts = append(opts, sv_filerequest.OptAllowLateUploads(fm.DeadlineAllowLateUploads))
		}
		req, err := sv_filerequest.New(z.Peer.Context().AsMemberId(member.TeamMemberId)).Create(
			fm.Title,
			mo_path.NewDropboxPath(fm.Destination),
			opts...,
		)
		if err != nil {
			z.OperationLog.Failure(err, fm)
		} else {
			z.OperationLog.Success(fm, req)
		}
		return nil
	})
}

func (z *Clone) Test(c app_control.Control) error {
	return qt_endtoend.HumanInteractionRequired()
}
