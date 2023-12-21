package file

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/essentials/model/mo_filter"
	"github.com/watermint/toolbox/essentials/model/mo_int"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/ingredient/ig_dropbox/ig_team/ig_namespace/ig_file"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
)

type Size struct {
	IncludeSharedFolder bool
	IncludeTeamFolder   bool
	IncludeMemberFolder bool
	IncludeAppFolder    bool
	Folder              mo_filter.Filter
	Depth               mo_int.RangeInt
	NamespaceSize       *ig_file.Size
	Peer                dbx_conn.ConnScopedTeam
}

func (z *Size) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeFilesMetadataRead,
		dbx_auth.ScopeMembersRead,
		dbx_auth.ScopeTeamDataMember,
		dbx_auth.ScopeTeamDataTeamSpace,
		dbx_auth.ScopeTeamInfoRead,
	)
	z.IncludeSharedFolder = true
	z.IncludeTeamFolder = true
	z.Depth.SetRange(1, 300, 3)
	z.Folder.SetOptions(
		mo_filter.NewNameFilter(),
		mo_filter.NewNamePrefixFilter(),
		mo_filter.NewNameSuffixFilter(),
	)
}

func (z *Size) Exec(c app_control.Control) error {
	return rc_exec.Exec(c, z.NamespaceSize, func(r rc_recipe.Recipe) {
		rc := r.(*ig_file.Size)
		rc.IncludeSharedFolder = z.IncludeSharedFolder
		rc.IncludeTeamFolder = z.IncludeTeamFolder
		rc.IncludeMemberFolder = z.IncludeMemberFolder
		rc.IncludeAppFolder = z.IncludeAppFolder
		rc.Folder = z.Folder
		rc.Depth = z.Depth
		rc.Peer = z.Peer
	})
}

func (z *Size) Test(c app_control.Control) error {
	return qt_errors.ErrorNoTestRequired
}
