package member

import (
	"errors"
	"github.com/watermint/toolbox/domain/model/mo_member"
	"github.com/watermint/toolbox/domain/service/sv_member"
	"github.com/watermint/toolbox/infra/api/api_util"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/quality/qt_test"
	"github.com/watermint/toolbox/infra/recpie/app_conn"
	"github.com/watermint/toolbox/infra/recpie/app_file"
	"github.com/watermint/toolbox/infra/recpie/app_kitchen"
	"github.com/watermint/toolbox/infra/recpie/app_vo"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/report/rp_spec"
	"github.com/watermint/toolbox/infra/report/rp_spec_impl"
	"github.com/watermint/toolbox/infra/ui/app_msg"
)

type InviteRow struct {
	Email     string
	GivenName string
	Surname   string
}

func (z *InviteRow) Validate() error {
	if z.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

type InviteVO struct {
	File app_file.Data
	Peer app_conn.ConnBusinessMgmt
}

const (
	reportInvite = "invite"
)

type Invite struct {
}

func (z *Invite) Reports() []rp_spec.ReportSpec {
	return []rp_spec.ReportSpec{
		rp_spec_impl.Spec(reportInvite, rp_model.TransactionHeader(&InviteRow{}, &mo_member.Member{})),
	}
}

func (z *Invite) Test(c app_control.Control) error {
	return qt_test.HumanInteractionRequired()
}

func (z *Invite) Console() {
}

func (z *Invite) Requirement() app_vo.ValueObject {
	return &InviteVO{}
}

func (z *Invite) msgFromTag(tag string) app_msg.Message {
	return app_msg.M("recipe.member.invite.tag." + tag)
}

func (z *Invite) Exec(k app_kitchen.Kitchen) error {
	var vo interface{} = k.Value()
	mvo := vo.(*InviteVO)

	connMgmt, err := mvo.Peer.Connect(k.Control())
	if err != nil {
		return err
	}

	svm := sv_member.New(connMgmt)
	rep, err := rp_spec_impl.New(z, k.Control()).Open(reportInvite)
	if err != nil {
		return err
	}
	defer rep.Close()

	if err := mvo.File.Model(k.Control(), &InviteRow{}); err != nil {
		return err
	}

	return mvo.File.EachRow(func(row interface{}, rowIndex int) error {
		m := row.(*InviteRow)
		if err = m.Validate(); err != nil {
			if rowIndex > 0 {
				rep.Failure(rp_model.MsgInvalidData, m, nil)
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

		r, err := svm.Add(m.Email, opts...)
		switch {
		case err != nil:
			rep.Failure(api_util.MsgFromError(err), m, nil)
			return nil

		case r.Tag == "success":
			rep.Success(m, r)
			return nil

		case r.Tag == "user_already_on_team":
			rep.Skip(z.msgFromTag(r.Tag), m, nil)
			return nil

		default:
			rep.Failure(z.msgFromTag(r.Tag), m, nil)
			return nil
		}
	})
}
