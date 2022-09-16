package file

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_client"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_file"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_file"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_profile"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
)

type List struct {
	Peer                         dbx_conn.ConnScopedIndividual
	Path                         mo_path.DropboxPath
	Recursive                    bool
	IncludeDeleted               bool
	IncludeMountedFolders        bool
	IncludeExplicitSharedMembers bool
	FileList                     rp_model.RowReport
}

func (z *List) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeAccountInfoRead,
		dbx_auth.ScopeFilesContentRead,
	)
	z.FileList.SetModel(
		&mo_file.ConcreteEntry{},
		rp_model.HiddenColumns(
			"id",
			"path_lower",
			"revision",
			"content_hash",
			"shared_folder_id",
			"parent_shared_folder_id",
		),
	)
}

func (z *List) Exec(c app_control.Control) error {
	profile, err := sv_profile.NewProfile(z.Peer.Client()).Current()
	if err != nil {
		return err
	}

	opts := make([]sv_file.ListOpt, 0)
	opts = append(opts, sv_file.IncludeDeleted(z.IncludeDeleted))
	opts = append(opts, sv_file.Recursive(z.Recursive))
	opts = append(opts, sv_file.IncludeHasExplicitSharedMembers(z.IncludeExplicitSharedMembers))
	opts = append(opts, sv_file.IncludeMountedFolders(z.IncludeMountedFolders))

	if err := z.FileList.Open(); err != nil {
		return err
	}
	client := z.Peer.Client().WithPath(
		dbx_client.Root(profile.RootNamespaceId),
	)

	return sv_file.NewFiles(client).ListEach(z.Path, func(entry mo_file.Entry) {
		z.FileList.Row(entry.Concrete())
	}, opts...)
}

func (z *List) Test(c app_control.Control) error {
	return rc_exec.ExecMock(c, &List{}, func(r rc_recipe.Recipe) {
		m := r.(*List)
		m.Path = qtr_endtoend.NewTestDropboxFolderPath("list")
	})
}
