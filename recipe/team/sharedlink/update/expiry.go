package update

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_util"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_member"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_sharedlink"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_time"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_member"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_sharedlink"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/essentials/model/mo_int"
	"github.com/watermint/toolbox/essentials/model/mo_string"
	"github.com/watermint/toolbox/essentials/time/ut_format"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"math"
	"time"
)

type MsgExpiry struct {
	ProgressScanning      app_msg.Message
	ProgressUpdating      app_msg.Message
	ErrorUnableScanMember app_msg.Message
}

var (
	MExpiry = app_msg.Apply(&MsgExpiry{}).(*MsgExpiry)
)

type ExpiryScanWorker struct {
	ctl        app_control.Control
	ctx        dbx_context.Context
	rep        rp_model.TransactionReport
	repSkipped rp_model.RowReport
	member     *mo_member.Member
	newExpiry  time.Time
	visibility string
}

func (z *ExpiryScanWorker) Exec() error {
	ui := z.ctl.UI()
	l := z.ctl.Log().With(esl.Any("member", z.member))

	l.Debug("Scanning member shared links")
	ui.Progress(MExpiry.ProgressScanning.With("MemberEmail", z.member.Email))

	ctxMember := z.ctx.AsMemberId(z.member.TeamMemberId)
	links, err := sv_sharedlink.New(ctxMember).List()
	if err != nil {
		l.Debug("Unable to scan shared link", esl.Error(err))
		ui.Error(MExpiry.ErrorUnableScanMember.
			With("Member", z.member.Email).
			With("Error", err))

		return err
	}

	q := z.ctl.NewLegacyQueue()

	for _, link := range links {
		ll := l.With(esl.Any("link", link))
		if link.LinkVisibility() != z.visibility {
			ll.Debug("Skip link", esl.String("targetVisibility", z.visibility))
			z.repSkipped.Row(mo_sharedlink.NewSharedLinkMember(link, z.member))
			continue
		}

		update := false

		switch {
		case link.LinkExpires() == "":
			ll.Debug("The link doesn't have expiration")
			update = true

		default:
			le, v := ut_format.ParseTimestamp(link.LinkExpires())
			if !v {
				ll.Warn("Invalid timestamp format from API response")
				continue
			}

			if le.IsZero() || le.After(z.newExpiry) {
				ll.Debug("The link have long or no expiration")
				update = true
			}
		}

		if !update {
			z.repSkipped.Row(mo_sharedlink.NewSharedLinkMember(link, z.member))
			ll.Debug("Skip")
			continue
		}

		q.Enqueue(&ExpiryWorker{
			ctl:       z.ctl,
			ctx:       ctxMember,
			rep:       z.rep,
			member:    z.member,
			link:      link,
			newExpiry: z.newExpiry,
		})
	}
	q.Wait()

	return nil
}

type ExpiryWorker struct {
	ctl       app_control.Control
	ctx       dbx_context.Context
	rep       rp_model.TransactionReport
	member    *mo_member.Member
	link      mo_sharedlink.SharedLink
	newExpiry time.Time
}

func (z *ExpiryWorker) Exec() error {
	ui := z.ctl.UI()
	l := z.ctl.Log().With(esl.Any("link", z.link.Metadata()))

	ui.Progress(MExpiry.ProgressUpdating.With("MemberEmail", z.member.Email).
		With("Url", z.link.LinkUrl()).
		With("CurrentExpiry", z.link.LinkExpires()).
		With("NewExpiry", dbx_util.ToApiTimeString(z.newExpiry)))

	updated, err := sv_sharedlink.New(z.ctx).Update(z.link, sv_sharedlink.Expires(z.newExpiry))
	if err != nil {
		l.Debug("Unable to update expiration")
		z.rep.Failure(err, mo_sharedlink.NewSharedLinkMember(z.link, z.member))
		return err
	}

	l.Debug("Updated", esl.Any("updated", updated))
	z.rep.Success(
		mo_sharedlink.NewSharedLinkMember(z.link, z.member),
		updated,
	)

	return nil
}

type Expiry struct {
	rc_recipe.RemarkIrreversible
	Peer                       dbx_conn.ConnBusinessFile
	Days                       mo_int.RangeInt
	At                         mo_time.TimeOptional
	Visibility                 mo_string.SelectString
	Updated                    rp_model.TransactionReport
	Skipped                    rp_model.RowReport
	ErrorPleaseSpecifyDaysOrAt app_msg.Message
	ErrorInvalidDateTime       app_msg.Message
}

func (z *Expiry) Preset() {
	z.Days.SetRange(0, math.MaxInt32, 0)
	z.Visibility.SetOptions("public", "public", "team_only", "password", "team_and_password", "shared_folder_only")
	z.Skipped.SetModel(&mo_sharedlink.SharedLinkMember{}, rp_model.HiddenColumns(
		"shared_link_id",
		"account_id",
		"team_member_id",
	))
	z.Updated.SetModel(&mo_sharedlink.SharedLinkMember{}, &mo_sharedlink.Metadata{}, rp_model.HiddenColumns(
		"input.shared_link_id",
		"input.account_id",
		"input.team_member_id",
		"result.tag",
		"result.id",
		"result.url",
		"result.name",
		"result.path_lower",
		"result.visibility",
	))
}

func (z *Expiry) Exec(c app_control.Control) error {
	ui := c.UI()
	l := c.Log()
	var newExpiry time.Time
	if z.Days.Value() > 0 && z.At.Ok() {
		l.Debug("Both Days/At specified", esl.Int("evo.Days", z.Days.Value()), esl.String("evo.At", z.At.Value()))
		ui.Error(z.ErrorPleaseSpecifyDaysOrAt)
		return errors.New("please specify one of `-days` or `-at`")
	}

	switch {
	case z.Days.Value() > 0:
		newExpiry = dbx_util.RebaseTime(time.Now().Add(time.Duration(z.Days.Value()*24) * time.Hour))
		l.Debug("New expiry", esl.Int("evo.Days", z.Days.Value()), esl.String("newExpiry", newExpiry.String()))

	default:
		if !z.At.Ok() {
			l.Debug("Invalid date/time format for at option", esl.String("evo.At", z.At.Value()))
			ui.Error(z.ErrorInvalidDateTime.With("Time", z.At.Value()))
			return errors.New("invalid date/time format for `at`")
		}
		newExpiry = z.At.Time()
	}

	l = l.With(esl.String("newExpiry", newExpiry.String()))

	if err := z.Updated.Open(); err != nil {
		return err
	}
	if err := z.Skipped.Open(); err != nil {
		return err
	}

	members, err := sv_member.New(z.Peer.Context()).List()
	if err != nil {
		return err
	}

	q := c.NewLegacyQueue()

	for _, member := range members {
		q.Enqueue(&ExpiryScanWorker{
			ctl:        c,
			ctx:        z.Peer.Context(),
			rep:        z.Updated,
			repSkipped: z.Skipped,
			member:     member,
			newExpiry:  newExpiry,
			visibility: z.Visibility.Value(),
		})
	}
	q.Wait()

	return nil
}

func (z *Expiry) Test(c app_control.Control) error {
	// should fail
	{
		err := rc_exec.Exec(c, &Expiry{}, func(r rc_recipe.Recipe) {
			rc := r.(*Expiry)
			rc.Days.SetValue(1)
			rc.At = mo_time.NewOptional(time.Now().Add(1 * 1000 * time.Millisecond))
		})
		if err == nil {
			return errors.New("days and at should not be accepted same time")
		}
	}

	{
		err := rc_exec.ExecMock(c, &Expiry{}, func(r rc_recipe.Recipe) {
			m := r.(*Expiry)
			m.Days.SetValue(7)
		})
		if e, _ := qt_errors.ErrorsForTest(c.Log(), err); e != nil {
			return e
		}
	}

	{
		err := rc_exec.ExecMock(c, &Expiry{}, func(r rc_recipe.Recipe) {
			m := r.(*Expiry)
			m.At = mo_time.NewOptional(time.Now().Add(1 * 1000 * time.Millisecond))
		})
		if e, _ := qt_errors.ErrorsForTest(c.Log(), err); e != nil {
			return e
		}
	}

	return nil
}
