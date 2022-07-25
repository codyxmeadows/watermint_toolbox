package mount

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_sharedfolder"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_sharedfolder_mount"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
)

type Mountable struct {
	Peer           dbx_conn.ConnScopedIndividual
	Mountables     rp_model.RowReport
	IncludeMounted bool
}

func (z *Mountable) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeSharingRead,
	)
	z.Mountables.SetModel(
		&mo_sharedfolder.SharedFolder{},
		rp_model.HiddenColumns(
			"parent_shared_folder_id",
			"team_member_id",
			"namespace_id",
			"owner_team_id",
		),
	)
}

func (z *Mountable) Exec(c app_control.Control) error {
	if err := z.Mountables.Open(); err != nil {
		return err
	}

	mounts, err := sv_sharedfolder_mount.New(z.Peer.Context()).List()
	if err != nil {
		return err
	}

	for _, mount := range mounts {
		if mount.PathLower != "" {
			if !z.IncludeMounted {
				continue
			}
		}
		z.Mountables.Row(mount)
	}
	return nil
}

func (z *Mountable) Test(c app_control.Control) error {
	return rc_exec.ExecMock(c, &Mountable{}, rc_recipe.NoCustomValues)
}
