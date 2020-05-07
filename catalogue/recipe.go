package catalogue

// Code generated by dev catalogue command DO NOT EDIT

import (
	infra_recipe_rc_recipe "github.com/watermint/toolbox/infra/recipe/rc_recipe"
	recipe "github.com/watermint/toolbox/recipe"
	recipeconfig "github.com/watermint/toolbox/recipe/config"
	recipeconnect "github.com/watermint/toolbox/recipe/connect"
	recipedev "github.com/watermint/toolbox/recipe/dev"
	recipedevciartifact "github.com/watermint/toolbox/recipe/dev/ci/artifact"
	recipedevciauth "github.com/watermint/toolbox/recipe/dev/ci/auth"
	recipedevdesktop "github.com/watermint/toolbox/recipe/dev/desktop"
	recipedevkvs "github.com/watermint/toolbox/recipe/dev/kvs"
	recipedevrelease "github.com/watermint/toolbox/recipe/dev/release"
	recipedevspec "github.com/watermint/toolbox/recipe/dev/spec"
	recipedevtest "github.com/watermint/toolbox/recipe/dev/test"
	recipedevutil "github.com/watermint/toolbox/recipe/dev/util"
	recipefile "github.com/watermint/toolbox/recipe/file"
	recipefilecompare "github.com/watermint/toolbox/recipe/file/compare"
	recipefiledispatch "github.com/watermint/toolbox/recipe/file/dispatch"
	recipefileexport "github.com/watermint/toolbox/recipe/file/export"
	recipefileimport "github.com/watermint/toolbox/recipe/file/import"
	recipefileimportbatch "github.com/watermint/toolbox/recipe/file/import/batch"
	recipefilesearch "github.com/watermint/toolbox/recipe/file/search"
	recipefilesync "github.com/watermint/toolbox/recipe/file/sync"
	recipefilesyncpreflight "github.com/watermint/toolbox/recipe/file/sync/preflight"
	recipefilerequest "github.com/watermint/toolbox/recipe/filerequest"
	recipefilerequestdelete "github.com/watermint/toolbox/recipe/filerequest/delete"
	recipegroup "github.com/watermint/toolbox/recipe/group"
	recipegroupbatch "github.com/watermint/toolbox/recipe/group/batch"
	recipegroupmember "github.com/watermint/toolbox/recipe/group/member"
	recipejob "github.com/watermint/toolbox/recipe/job"
	recipejobhistory "github.com/watermint/toolbox/recipe/job/history"
	recipemember "github.com/watermint/toolbox/recipe/member"
	recipememberquota "github.com/watermint/toolbox/recipe/member/quota"
	recipememberupdate "github.com/watermint/toolbox/recipe/member/update"
	recipeservicesgithub "github.com/watermint/toolbox/recipe/services/github"
	recipeservicesgithubissue "github.com/watermint/toolbox/recipe/services/github/issue"
	recipeservicesgithubrelease "github.com/watermint/toolbox/recipe/services/github/release"
	recipeservicesgithubreleaseasset "github.com/watermint/toolbox/recipe/services/github/release/asset"
	recipeservicesgithubtag "github.com/watermint/toolbox/recipe/services/github/tag"
	recipesharedfolder "github.com/watermint/toolbox/recipe/sharedfolder"
	recipesharedfoldermember "github.com/watermint/toolbox/recipe/sharedfolder/member"
	recipesharedlink "github.com/watermint/toolbox/recipe/sharedlink"
	recipesharedlinkfile "github.com/watermint/toolbox/recipe/sharedlink/file"
	recipeteam "github.com/watermint/toolbox/recipe/team"
	recipeteamactivity "github.com/watermint/toolbox/recipe/team/activity"
	recipeteamactivitybatch "github.com/watermint/toolbox/recipe/team/activity/batch"
	recipeteamactivitydaily "github.com/watermint/toolbox/recipe/team/activity/daily"
	recipeteamcontent "github.com/watermint/toolbox/recipe/team/content"
	recipeteamdevice "github.com/watermint/toolbox/recipe/team/device"
	recipeteamdiag "github.com/watermint/toolbox/recipe/team/diag"
	recipeteamfilerequest "github.com/watermint/toolbox/recipe/team/filerequest"
	recipeteamlinkedapp "github.com/watermint/toolbox/recipe/team/linkedapp"
	recipeteamnamespace "github.com/watermint/toolbox/recipe/team/namespace"
	recipeteamnamespacefile "github.com/watermint/toolbox/recipe/team/namespace/file"
	recipeteamnamespacemember "github.com/watermint/toolbox/recipe/team/namespace/member"
	recipeteamsharedlink "github.com/watermint/toolbox/recipe/team/sharedlink"
	recipeteamsharedlinkupdate "github.com/watermint/toolbox/recipe/team/sharedlink/update"
	recipeteamfolder "github.com/watermint/toolbox/recipe/teamfolder"
	recipeteamfolderbatch "github.com/watermint/toolbox/recipe/teamfolder/batch"
	recipeteamfolderfile "github.com/watermint/toolbox/recipe/teamfolder/file"
)

func AutoDetectedRecipes() []infra_recipe_rc_recipe.Recipe {
	return []infra_recipe_rc_recipe.Recipe{
		&recipe.License{},
		&recipe.Version{},
		&recipeconfig.Disable{},
		&recipeconfig.Enable{},
		&recipeconfig.Features{},
		&recipeconnect.BusinessAudit{},
		&recipeconnect.BusinessFile{},
		&recipeconnect.BusinessInfo{},
		&recipeconnect.BusinessMgmt{},
		&recipeconnect.UserFile{},
		&recipedev.Async{},
		&recipedev.Catalogue{},
		&recipedev.Doc{},
		&recipedev.Dummy{},
		&recipedev.Echo{},
		&recipedev.Preflight{},
		&recipedevciartifact.Connect{},
		&recipedevciartifact.Up{},
		&recipedevciauth.Connect{},
		&recipedevciauth.Export{},
		&recipedevciauth.Import{},
		&recipedevdesktop.Install{},
		&recipedevdesktop.Start{},
		&recipedevdesktop.Stop{},
		&recipedevdesktop.Suspendupdate{},
		&recipedevkvs.Dump{},
		&recipedevrelease.Candidate{},
		&recipedevrelease.Publish{},
		&recipedevspec.Diff{},
		&recipedevspec.Doc{},
		&recipedevtest.Monkey{},
		&recipedevtest.Recipe{},
		&recipedevtest.Resources{},
		&recipedevutil.Curl{},
		&recipedevutil.Wait{},
		&recipefile.Copy{},
		&recipefile.Delete{},
		&recipefile.Download{},
		&recipefile.List{},
		&recipefile.Merge{},
		&recipefile.Move{},
		&recipefile.Replication{},
		&recipefile.Restore{},
		&recipefile.Upload{},
		&recipefile.Watch{},
		&recipefilecompare.Account{},
		&recipefilecompare.Local{},
		&recipefiledispatch.Local{},
		&recipefileexport.Doc{},
		&recipefileimport.Url{},
		&recipefileimportbatch.Url{},
		&recipefilesearch.Content{},
		&recipefilesearch.Name{},
		&recipefilesync.Up{},
		&recipefilesyncpreflight.Up{},
		&recipefilerequest.Create{},
		&recipefilerequest.List{},
		&recipefilerequestdelete.Closed{},
		&recipefilerequestdelete.Url{},
		&recipegroup.Add{},
		&recipegroup.Delete{},
		&recipegroup.List{},
		&recipegroup.Rename{},
		&recipegroupbatch.Delete{},
		&recipegroupmember.Add{},
		&recipegroupmember.Delete{},
		&recipegroupmember.List{},
		&recipejob.Loop{},
		&recipejob.Run{},
		&recipejobhistory.Archive{},
		&recipejobhistory.Delete{},
		&recipejobhistory.List{},
		&recipejobhistory.Ship{},
		&recipemember.Delete{},
		&recipemember.Detach{},
		&recipemember.Invite{},
		&recipemember.List{},
		&recipemember.Reinvite{},
		&recipemember.Replication{},
		&recipememberquota.List{},
		&recipememberquota.Update{},
		&recipememberquota.Usage{},
		&recipememberupdate.Email{},
		&recipememberupdate.Externalid{},
		&recipememberupdate.Profile{},
		&recipeservicesgithub.Profile{},
		&recipeservicesgithubissue.List{},
		&recipeservicesgithubrelease.Draft{},
		&recipeservicesgithubrelease.List{},
		&recipeservicesgithubreleaseasset.Download{},
		&recipeservicesgithubreleaseasset.List{},
		&recipeservicesgithubreleaseasset.Upload{},
		&recipeservicesgithubtag.Create{},
		&recipesharedfolder.List{},
		&recipesharedfoldermember.List{},
		&recipesharedlink.Create{},
		&recipesharedlink.Delete{},
		&recipesharedlink.List{},
		&recipesharedlinkfile.List{},
		&recipeteam.Feature{},
		&recipeteam.Info{},
		&recipeteamactivity.Event{},
		&recipeteamactivity.User{},
		&recipeteamactivitybatch.User{},
		&recipeteamactivitydaily.Event{},
		&recipeteamcontent.Member{},
		&recipeteamcontent.Policy{},
		&recipeteamdevice.List{},
		&recipeteamdevice.Unlink{},
		&recipeteamdiag.Explorer{},
		&recipeteamfilerequest.Clone{},
		&recipeteamfilerequest.List{},
		&recipeteamlinkedapp.List{},
		&recipeteamnamespace.List{},
		&recipeteamnamespacefile.List{},
		&recipeteamnamespacefile.Size{},
		&recipeteamnamespacemember.List{},
		&recipeteamsharedlink.List{},
		&recipeteamsharedlinkupdate.Expiry{},
		&recipeteamfolder.Archive{},
		&recipeteamfolder.List{},
		&recipeteamfolder.Permdelete{},
		&recipeteamfolder.Replication{},
		&recipeteamfolderbatch.Archive{},
		&recipeteamfolderbatch.Permdelete{},
		&recipeteamfolderbatch.Replication{},
		&recipeteamfolderfile.List{},
		&recipeteamfolderfile.Size{},
	}
}
