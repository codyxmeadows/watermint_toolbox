package file

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_team"
	"github.com/watermint/toolbox/essentials/model/mo_filter"
	"github.com/watermint/toolbox/essentials/model/mo_int"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	namespacefile "github.com/watermint/toolbox/ingredient/team/namespace/file"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"math"
)

type Size struct {
	Peer                       dbx_conn.ConnScopedTeam
	FileSize                   *namespacefile.Size
	Depth                      mo_int.RangeInt
	Folder                     mo_filter.Filter
	ErrorTeamSpaceNotSupported app_msg.Message
}

func (z *Size) Preset() {
	z.Peer.SetScopes(
		dbx_auth.ScopeFilesMetadataRead,
		dbx_auth.ScopeMembersRead,
		dbx_auth.ScopeTeamDataMember,
		dbx_auth.ScopeTeamDataTeamSpace,
		dbx_auth.ScopeTeamInfoRead,
	)
	z.Depth.SetRange(1, math.MaxInt32, 3)
	z.Folder.SetOptions(
		mo_filter.NewNameFilter(),
		mo_filter.NewNamePrefixFilter(),
		mo_filter.NewNameSuffixFilter(),
	)
}

func (z *Size) Exec(c app_control.Control) error {
	if ok, _ := sv_team.UnlessTeamFolderApiSupported(z.Peer.Client()); ok {
		c.UI().Error(z.ErrorTeamSpaceNotSupported)
		return errors.New("team space is not supported by this command")
	}

	return rc_exec.Exec(c, z.FileSize, func(r rc_recipe.Recipe) {
		rc := r.(*namespacefile.Size)
		rc.IncludeSharedFolder = false
		rc.IncludeTeamFolder = true
		rc.Depth = z.Depth
		rc.Folder = z.Folder
		rc.Peer = z.Peer
	})
}

func (z *Size) Test(c app_control.Control) error {
	return qt_errors.ErrorNoTestRequired
}
