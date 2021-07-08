package quota

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_member"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_member_quota"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_member"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_member_quota"
	"github.com/watermint/toolbox/essentials/go/es_goroutine"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/essentials/model/mo_int"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/feed/fd_file"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/infra/qt_file"
	"math"
)

type MsgUpdate struct {
	ProgressUpdate app_msg.Message
}

var (
	MUpdate = app_msg.Apply(&MsgUpdate{}).(*MsgUpdate)
)

type UpdateWorker struct {
	member *mo_member.Member
	quota  int

	ctl app_control.Control
	ctx dbx_context.Context
	rep rp_model.TransactionReport
}

func (z *UpdateWorker) Exec() error {
	l := z.ctl.Log()
	z.ctl.UI().Progress(MUpdate.ProgressUpdate.With("MemberEmail", z.member.Email).With("Quota", z.quota))
	l.Debug("Updating quota", esl.String("Routine", es_goroutine.GetGoRoutineName()), esl.Any("Member", z.member))

	q := &mo_member_quota.Quota{
		TeamMemberId: z.member.TeamMemberId,
		Quota:        z.quota,
	}
	in := mo_member_quota.MemberQuota{
		Email: z.member.Email,
		Quota: z.quota,
	}

	newQuota, err := sv_member_quota.NewQuota(z.ctx).Update(q)
	if err != nil {
		z.rep.Failure(err, in)
	} else {
		z.rep.Success(in, mo_member_quota.NewMemberQuota(z.member, newQuota))
	}
	return nil
}

type Update struct {
	Peer         dbx_conn.ConnScopedTeam
	File         fd_file.RowFeed
	OperationLog rp_model.TransactionReport
	Quota        mo_int.RangeInt
}

func (z *Update) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeMembersRead,
		dbx_auth.ScopeMembersWrite,
	)
	z.File.SetModel(&mo_member_quota.MemberQuota{})
	z.OperationLog.SetModel(&mo_member_quota.MemberQuota{}, &mo_member_quota.MemberQuota{})
	z.Quota.SetRange(0, math.MaxInt32, 0)
}

func (z *Update) Exec(c app_control.Control) error {
	ctx := z.Peer.Context()

	members, err := sv_member.New(ctx).List()
	if err != nil {
		return err
	}
	emailToMember := mo_member.MapByEmail(members)

	err = z.OperationLog.Open()
	if err != nil {
		return err
	}

	q := c.NewLegacyQueue()

	err = z.File.EachRow(func(m interface{}, rowIndex int) error {
		mq := m.(*mo_member_quota.MemberQuota)
		member, ok := emailToMember[mq.Email]
		if !ok {
			z.OperationLog.Failure(errors.New("member not found for an email"), mq)
			return nil
		}
		quota := z.Quota.Value()
		if mq.Quota != 0 {
			quota = mq.Quota
		}

		q.Enqueue(&UpdateWorker{
			member: member,
			quota:  quota,
			ctl:    c,
			ctx:    ctx,
			rep:    z.OperationLog,
		})
		return nil
	})
	q.Wait()

	return err
}

func (z *Update) Test(c app_control.Control) error {
	err := rc_exec.ExecMock(c, &Update{}, func(r rc_recipe.Recipe) {
		f, err := qt_file.MakeTestFile("update-quota", "john@example.com,10")
		if err != nil {
			return
		}
		m := r.(*Update)
		m.Quota.SetValue(150)
		m.File.SetFilePath(f)
	})
	if e, _ := qt_errors.ErrorsForTest(c.Log(), err); e != nil {
		return e
	}
	return qt_errors.ErrorHumanInteractionRequired
}
