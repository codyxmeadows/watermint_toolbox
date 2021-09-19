package sharedfolder

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_sharedfolder"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_sharedfolder"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
)

type List struct {
	Peer         dbx_conn.ConnScopedIndividual
	SharedFolder rp_model.RowReport
}

func (z *List) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeSharingRead,
	)
	z.SharedFolder.SetModel(
		&mo_sharedfolder.SharedFolder{},
		rp_model.HiddenColumns(
			"shared_folder_id",
			"parent_shared_folder_id",
			"owner_team_id",
		),
	)
}

func (z *List) Test(c app_control.Control) error {
	if err := rc_exec.Exec(c, &List{}, rc_recipe.NoCustomValues); err != nil {
		return err
	}
	return qtr_endtoend.TestRows(c, "shared_folder", func(cols map[string]string) error {
		if _, ok := cols["name"]; !ok {
			return errors.New("name is not found")
		}
		return nil
	})
}

func (z *List) Exec(c app_control.Control) error {
	c.Log().Debug("Scanning folders")
	folders, err := sv_sharedfolder.New(z.Peer.Context()).List()
	if err != nil {
		return err
	}

	if err := z.SharedFolder.Open(); err != nil {
		return err
	}

	for _, folder := range folders {
		z.SharedFolder.Row(folder)
	}
	return nil
}
