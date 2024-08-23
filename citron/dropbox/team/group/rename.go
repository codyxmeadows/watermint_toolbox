package group

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_group"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_group"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
)

type Rename struct {
	rc_recipe.RemarkIrreversible
	Peer         dbx_conn.ConnScopedTeam
	CurrentName  string
	NewName      string
	OperationLog rp_model.TransactionReport
}

type RenameRow struct {
	CurrentName string `json:"current_name"`
	NewName     string `json:"new_name"`
}

func (z *Rename) Exec(c app_control.Control) error {
	row := &RenameRow{
		CurrentName: z.CurrentName,
		NewName:     z.NewName,
	}
	if err := z.OperationLog.Open(); err != nil {
		return err
	}

	current, err := sv_group.New(z.Peer.Client()).ResolveByName(z.CurrentName)
	if err != nil {
		z.OperationLog.Failure(err, row)
		return err
	}
	current.GroupName = z.NewName
	updated, err := sv_group.New(z.Peer.Client()).Update(current)
	if err != nil {
		z.OperationLog.Failure(err, row)
		return err
	}

	z.OperationLog.Success(row, updated)
	return nil
}

func (z *Rename) Test(c app_control.Control) error {
	err := rc_exec.ExecMock(c, &Rename{}, func(r rc_recipe.Recipe) {
		m := r.(*Rename)
		m.CurrentName = "Marketing"
		m.NewName = "Marketing (Archived)"
	})
	if err, _ = qt_errors.ErrorsForTest(c.Log(), err); err != nil && !errors.Is(err, sv_group.ErrorGroupNotFoundForName) {
		return err
	}
	return qt_errors.ErrorScenarioTest
}

func (z *Rename) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeGroupsWrite,
	)
	z.OperationLog.SetModel(
		&RenameRow{},
		&mo_group.Group{},
		rp_model.HiddenColumns(
			"result.group_id",
			"result.group_external_id",
		),
	)
}
