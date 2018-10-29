package cmd_sharedlink

import (
	"flag"
	"github.com/watermint/toolbox/cmdlet"
	"github.com/watermint/toolbox/dbx_api"
	"github.com/watermint/toolbox/dbx_api/dbx_profile"
	"github.com/watermint/toolbox/dbx_api/dbx_sharing"
	"github.com/watermint/toolbox/dbx_api/dbx_team"
)

type CmdTeamSharedLinkList struct {
	*cmdlet.SimpleCommandlet

	apiContext *dbx_api.Context
	report     cmdlet.Report
	filter     cmdlet.SharedLinkFilter
}

func (CmdTeamSharedLinkList) Name() string {
	return "list"
}

func (CmdTeamSharedLinkList) Desc() string {
	return "List all shared links of the team members' accounts"
}

func (CmdTeamSharedLinkList) Usage() string {
	return ""
}

func (c *CmdTeamSharedLinkList) FlagConfig(f *flag.FlagSet) {
	c.report.FlagConfig(f)
	c.filter.FlagConfig(f)
}

func (c *CmdTeamSharedLinkList) Exec(args []string) {
	apiMgmt, err := c.ExecContext.LoadOrAuthBusinessFile()
	if err != nil {
		return
	}
	c.report.Open(c)
	defer c.report.Close()

	ml := dbx_team.MembersList{
		OnError: c.DefaultErrorHandler,
		OnEntry: func(member *dbx_profile.Member) bool {
			sl := dbx_sharing.SharedLinkList{
				AsMemberId:    member.Profile.TeamMemberId,
				AsMemberEmail: member.Profile.Email,
				OnError:       c.DefaultErrorHandler,
				OnEntry: func(link *dbx_sharing.SharedLink) bool {
					if c.filter.IsAcceptable(link) {
						c.report.Report(link)
					}
					return true
				},
			}
			sl.List(apiMgmt)
			return true
		},
	}
	ml.List(apiMgmt, false)
}
