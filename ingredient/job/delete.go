package job

import (
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/essentials/model/mo_string"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_job"
	"github.com/watermint/toolbox/infra/control/app_job_impl"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/quality/infra/qt_file"
	"os"
	"time"
)

type Delete struct {
	Path              mo_string.OptionalString
	Days              int
	ProgressDeleting  app_msg.Message
	ErrorFailedDelete app_msg.Message
}

func (z *Delete) Exec(c app_control.Control) error {
	histories, err := app_job_impl.GetHistories(z.Path)
	if err != nil {
		return err
	}
	threshold := time.Now().Add(time.Duration(-z.Days*24) * time.Hour)
	l := c.Log()

	for _, h := range histories {
		ts, found := h.TimeStart()
		if !found {
			l.Debug("Skip: Time start not found for the job")
			continue
		}
		if h.JobId() == c.Workspace().JobId() {
			l.Debug("Skip current job")
			continue
		}
		if h.IsNested() {
			l.Debug("Skip nested job")
			continue
		}
		if ts.After(threshold) {
			l.Debug("Skip: Time start is in range of retain", esl.String("jobId", h.JobId()))
			continue
		}
		c.UI().Info(z.ProgressDeleting.With("JobId", h.JobId()))
		if ho, ok := h.(app_job.HistoryOperation); ok {
			err := ho.Delete()
			if err != nil {
				l.Debug("Unable to archive", esl.Error(err), esl.Any("history", h))
				c.UI().Error(z.ProgressDeleting.With("JobId", h.JobId()).With("Error", err.Error()))
			}
		} else {
			l.Warn("This history is not supported to delete", esl.String("jobId", h.JobId()))
		}
	}
	return nil
}

func (z *Delete) Test(c app_control.Control) error {
	workspace, err := qt_file.MakeTestFolder("delete", false)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(workspace)
	}()

	return rc_exec.Exec(c, &Delete{}, func(r rc_recipe.Recipe) {
		m := r.(*Delete)
		m.Days = 365
		m.Path = mo_string.NewOptional(workspace)
	})
}

func (z *Delete) Preset() {
	z.Days = 28
}
