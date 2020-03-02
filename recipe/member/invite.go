package member

import (
	"errors"
	"github.com/watermint/toolbox/domain/model/mo_member"
	"github.com/watermint/toolbox/domain/service/sv_member"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/feed/fd_file"
	"github.com/watermint/toolbox/infra/recipe/rc_conn"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
)

type InviteRow struct {
	Email     string `json:"email"`
	GivenName string `json:"given_name"`
	Surname   string `json:"surname"`
}

func (z *InviteRow) Validate() error {
	if z.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

type Invite struct {
	File         fd_file.RowFeed
	Peer         rc_conn.ConnBusinessMgmt
	OperationLog rp_model.TransactionReport
	SilentInvite bool
}

func (z *Invite) Preset() {
	z.File.SetModel(&InviteRow{})
	z.OperationLog.SetModel(&InviteRow{}, &mo_member.Member{})
}

func (z *Invite) Test(c app_control.Control) error {
	return qt_errors.ErrorHumanInteractionRequired
}

func (z *Invite) msgFromTag(tag string) app_msg.Message {
	return app_msg.M("recipe.member.invite.tag." + tag)
}

func (z *Invite) Exec(c app_control.Control) error {
	ctx := z.Peer.Context()

	svm := sv_member.New(ctx)
	err := z.OperationLog.Open()
	if err != nil {
		return err
	}

	return z.File.EachRow(func(row interface{}, rowIndex int) error {
		m := row.(*InviteRow)
		if err = m.Validate(); err != nil {
			if rowIndex > 0 {
				z.OperationLog.Failure(err, m)
			}
			return nil
		}
		opts := make([]sv_member.AddOpt, 0)
		if m.GivenName != "" {
			opts = append(opts, sv_member.AddWithGivenName(m.GivenName))
		}
		if m.Surname != "" {
			opts = append(opts, sv_member.AddWithSurname(m.Surname))
		}
		if z.SilentInvite {
			opts = append(opts, sv_member.AddWithoutSendWelcomeEmail())
		}

		r, err := svm.Add(m.Email, opts...)
		switch {
		case err != nil:
			z.OperationLog.Failure(err, m)
			return nil

		case r.Tag == "success":
			z.OperationLog.Success(m, r)
			return nil

		case r.Tag == "user_already_on_team":
			z.OperationLog.Skip(z.msgFromTag(r.Tag), m)
			return nil

		default:
			// TODO: i18n
			z.OperationLog.Failure(errors.New("failure due to "+r.Tag), m)
			return nil
		}
	})
}
