package namespace

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_namespace"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_namespace"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
)

type List struct {
	Peer      dbx_conn.ConnScopedTeam
	Namespace rp_model.RowReport
}

func (z *List) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeTeamDataMember,
	)
	z.Namespace.SetModel(&mo_namespace.Namespace{},
		rp_model.HiddenColumns(
			"namespace_id",
		),
	)
}

func (z *List) Exec(c app_control.Control) error {
	if err := z.Namespace.Open(); err != nil {
		return err
	}

	namespaces, err := sv_namespace.New(z.Peer.Context()).List()
	if err != nil {
		return err
	}

	for _, namespace := range namespaces {
		z.Namespace.Row(namespace)
	}
	return nil
}

func (z *List) Test(c app_control.Control) error {
	if err := rc_exec.Exec(c, &List{}, rc_recipe.NoCustomValues); err != nil {
		return err
	}
	return qtr_endtoend.TestRows(c, "namespace", func(cols map[string]string) error {
		if _, ok := cols["name"]; !ok {
			return errors.New("`name` is not found")
		}
		return nil
	})
}
