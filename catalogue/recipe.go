package catalogue

// Code generated by dev catalogue command DO NOT EDIT

import (
	infra_recipe_rc_recipe "github.com/watermint/toolbox/infra/recipe/rc_recipe"
	recipe "github.com/watermint/toolbox/recipe"
	recipeconfig "github.com/watermint/toolbox/recipe/config"
	recipeconnect "github.com/watermint/toolbox/recipe/connect"
	recipedevbenchmark "github.com/watermint/toolbox/recipe/dev/benchmark"
	recipedevbuild "github.com/watermint/toolbox/recipe/dev/build"
	recipedevciartifact "github.com/watermint/toolbox/recipe/dev/ci/artifact"
	recipedevciauth "github.com/watermint/toolbox/recipe/dev/ci/auth"
	recipedevdiag "github.com/watermint/toolbox/recipe/dev/diag"
	recipedevkvs "github.com/watermint/toolbox/recipe/dev/kvs"
	recipedevrelease "github.com/watermint/toolbox/recipe/dev/release"
	recipedevreplay "github.com/watermint/toolbox/recipe/dev/replay"
	recipedevspec "github.com/watermint/toolbox/recipe/dev/spec"
	recipedevstage "github.com/watermint/toolbox/recipe/dev/stage"
	recipedevtest "github.com/watermint/toolbox/recipe/dev/test"
	recipedevutil "github.com/watermint/toolbox/recipe/dev/util"
	recipedevutilimage "github.com/watermint/toolbox/recipe/dev/util/image"
	recipefile "github.com/watermint/toolbox/recipe/file"
	recipefilecompare "github.com/watermint/toolbox/recipe/file/compare"
	recipefiledispatch "github.com/watermint/toolbox/recipe/file/dispatch"
	recipefileexport "github.com/watermint/toolbox/recipe/file/export"
	recipefileimport "github.com/watermint/toolbox/recipe/file/import"
	recipefileimportbatch "github.com/watermint/toolbox/recipe/file/import/batch"
	recipefilemount "github.com/watermint/toolbox/recipe/file/mount"
	recipefilesearch "github.com/watermint/toolbox/recipe/file/search"
	recipefilesync "github.com/watermint/toolbox/recipe/file/sync"
	recipefilerequest "github.com/watermint/toolbox/recipe/filerequest"
	recipefilerequestdelete "github.com/watermint/toolbox/recipe/filerequest/delete"
	recipegroup "github.com/watermint/toolbox/recipe/group"
	recipegroupbatch "github.com/watermint/toolbox/recipe/group/batch"
	recipegroupfolder "github.com/watermint/toolbox/recipe/group/folder"
	recipegroupmember "github.com/watermint/toolbox/recipe/group/member"
	recipeimage "github.com/watermint/toolbox/recipe/image"
	recipejobhistory "github.com/watermint/toolbox/recipe/job/history"
	recipejoblog "github.com/watermint/toolbox/recipe/job/log"
	recipemember "github.com/watermint/toolbox/recipe/member"
	recipememberclear "github.com/watermint/toolbox/recipe/member/clear"
	recipememberfile "github.com/watermint/toolbox/recipe/member/file"
	recipememberfolder "github.com/watermint/toolbox/recipe/member/folder"
	recipememberquota "github.com/watermint/toolbox/recipe/member/quota"
	recipememberupdate "github.com/watermint/toolbox/recipe/member/update"
	recipeservicesasanateam "github.com/watermint/toolbox/recipe/services/asana/team"
	recipeservicesasanateamproject "github.com/watermint/toolbox/recipe/services/asana/team/project"
	recipeservicesasanateamtask "github.com/watermint/toolbox/recipe/services/asana/team/task"
	recipeservicesasanaworkspace "github.com/watermint/toolbox/recipe/services/asana/workspace"
	recipeservicesasanaworkspaceproject "github.com/watermint/toolbox/recipe/services/asana/workspace/project"
	recipeservicesgithub "github.com/watermint/toolbox/recipe/services/github"
	recipeservicesgithubcontent "github.com/watermint/toolbox/recipe/services/github/content"
	recipeservicesgithubissue "github.com/watermint/toolbox/recipe/services/github/issue"
	recipeservicesgithubrelease "github.com/watermint/toolbox/recipe/services/github/release"
	recipeservicesgithubreleaseasset "github.com/watermint/toolbox/recipe/services/github/release/asset"
	recipeservicesgithubtag "github.com/watermint/toolbox/recipe/services/github/tag"
	recipeservicesgooglemailfilter "github.com/watermint/toolbox/recipe/services/google/mail/filter"
	recipeservicesgooglemailfilterbatch "github.com/watermint/toolbox/recipe/services/google/mail/filter/batch"
	recipeservicesgooglemaillabel "github.com/watermint/toolbox/recipe/services/google/mail/label"
	recipeservicesgooglemailmessage "github.com/watermint/toolbox/recipe/services/google/mail/message"
	recipeservicesgooglemailmessagelabel "github.com/watermint/toolbox/recipe/services/google/mail/message/label"
	recipeservicesgooglemailmessageprocessed "github.com/watermint/toolbox/recipe/services/google/mail/message/processed"
	recipeservicesgooglemailthread "github.com/watermint/toolbox/recipe/services/google/mail/thread"
	recipeservicesslackconversation "github.com/watermint/toolbox/recipe/services/slack/conversation"
	recipesharedfolder "github.com/watermint/toolbox/recipe/sharedfolder"
	recipesharedfoldermember "github.com/watermint/toolbox/recipe/sharedfolder/member"
	recipesharedlink "github.com/watermint/toolbox/recipe/sharedlink"
	recipesharedlinkfile "github.com/watermint/toolbox/recipe/sharedlink/file"
	recipeteam "github.com/watermint/toolbox/recipe/team"
	recipeteamactivity "github.com/watermint/toolbox/recipe/team/activity"
	recipeteamactivitybatch "github.com/watermint/toolbox/recipe/team/activity/batch"
	recipeteamactivitydaily "github.com/watermint/toolbox/recipe/team/activity/daily"
	recipeteamcontentmember "github.com/watermint/toolbox/recipe/team/content/member"
	recipeteamcontentmount "github.com/watermint/toolbox/recipe/team/content/mount"
	recipeteamcontentpolicy "github.com/watermint/toolbox/recipe/team/content/policy"
	recipeteamdevice "github.com/watermint/toolbox/recipe/team/device"
	recipeteamdiag "github.com/watermint/toolbox/recipe/team/diag"
	recipeteamfilerequest "github.com/watermint/toolbox/recipe/team/filerequest"
	recipeteamlinkedapp "github.com/watermint/toolbox/recipe/team/linkedapp"
	recipeteamnamespace "github.com/watermint/toolbox/recipe/team/namespace"
	recipeteamnamespacefile "github.com/watermint/toolbox/recipe/team/namespace/file"
	recipeteamnamespacemember "github.com/watermint/toolbox/recipe/team/namespace/member"
	recipeteamreport "github.com/watermint/toolbox/recipe/team/report"
	recipeteamsharedlink "github.com/watermint/toolbox/recipe/team/sharedlink"
	recipeteamsharedlinkupdate "github.com/watermint/toolbox/recipe/team/sharedlink/update"
	recipeteamfolder "github.com/watermint/toolbox/recipe/teamfolder"
	recipeteamfolderbatch "github.com/watermint/toolbox/recipe/teamfolder/batch"
	recipeteamfolderfile "github.com/watermint/toolbox/recipe/teamfolder/file"
	recipeteamfoldermember "github.com/watermint/toolbox/recipe/teamfolder/member"
	recipeteamfolderpolicy "github.com/watermint/toolbox/recipe/teamfolder/policy"
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
		&recipedevbenchmark.Local{},
		&recipedevbenchmark.Upload{},
		&recipedevbuild.Catalogue{},
		&recipedevbuild.Doc{},
		&recipedevbuild.License{},
		&recipedevbuild.Preflight{},
		&recipedevbuild.Readme{},
		&recipedevciartifact.Connect{},
		&recipedevciartifact.Up{},
		&recipedevciauth.Connect{},
		&recipedevciauth.Export{},
		&recipedevciauth.Import{},
		&recipedevdiag.Endpoint{},
		&recipedevdiag.Throughput{},
		&recipedevkvs.Dump{},
		&recipedevrelease.Candidate{},
		&recipedevrelease.Publish{},
		&recipedevreplay.Approve{},
		&recipedevreplay.Bundle{},
		&recipedevreplay.Recipe{},
		&recipedevreplay.Remote{},
		&recipedevspec.Diff{},
		&recipedevspec.Doc{},
		&recipedevstage.Gmail{},
		&recipedevstage.Gui{},
		&recipedevstage.Scoped{},
		&recipedevstage.Teamfolder{},
		&recipedevtest.Echo{},
		&recipedevtest.Kvsfootprint{},
		&recipedevtest.Monkey{},
		&recipedevtest.Recipe{},
		&recipedevtest.Resources{},
		&recipedevutil.Anonymise{},
		&recipedevutil.Curl{},
		&recipedevutil.Wait{},
		&recipedevutilimage.Jpeg{},
		&recipefile.Copy{},
		&recipefile.Delete{},
		&recipefile.Download{},
		&recipefile.List{},
		&recipefile.Merge{},
		&recipefile.Move{},
		&recipefile.Replication{},
		&recipefile.Restore{},
		&recipefile.Size{},
		&recipefile.Watch{},
		&recipefilecompare.Account{},
		&recipefilecompare.Local{},
		&recipefiledispatch.Local{},
		&recipefileexport.Doc{},
		&recipefileimport.Url{},
		&recipefileimportbatch.Url{},
		&recipefilemount.List{},
		&recipefilesearch.Content{},
		&recipefilesearch.Name{},
		&recipefilesync.Down{},
		&recipefilesync.Online{},
		&recipefilesync.Up{},
		&recipefilerequest.Create{},
		&recipefilerequest.List{},
		&recipefilerequestdelete.Closed{},
		&recipefilerequestdelete.Url{},
		&recipegroup.Add{},
		&recipegroup.Delete{},
		&recipegroup.List{},
		&recipegroup.Rename{},
		&recipegroupbatch.Delete{},
		&recipegroupfolder.List{},
		&recipegroupmember.Add{},
		&recipegroupmember.Delete{},
		&recipegroupmember.List{},
		&recipeimage.Info{},
		&recipejobhistory.Archive{},
		&recipejobhistory.Delete{},
		&recipejobhistory.List{},
		&recipejobhistory.Ship{},
		&recipejoblog.Jobid{},
		&recipejoblog.Kind{},
		&recipejoblog.Last{},
		&recipemember.Delete{},
		&recipemember.Detach{},
		&recipemember.Invite{},
		&recipemember.List{},
		&recipemember.Reinvite{},
		&recipemember.Replication{},
		&recipememberclear.Externalid{},
		&recipememberfile.Permdelete{},
		&recipememberfolder.List{},
		&recipememberquota.List{},
		&recipememberquota.Update{},
		&recipememberquota.Usage{},
		&recipememberupdate.Email{},
		&recipememberupdate.Externalid{},
		&recipememberupdate.Profile{},
		&recipeservicesasanateam.List{},
		&recipeservicesasanateamproject.List{},
		&recipeservicesasanateamtask.List{},
		&recipeservicesasanaworkspace.List{},
		&recipeservicesasanaworkspaceproject.List{},
		&recipeservicesgithub.Profile{},
		&recipeservicesgithubcontent.Get{},
		&recipeservicesgithubcontent.Put{},
		&recipeservicesgithubissue.List{},
		&recipeservicesgithubrelease.Draft{},
		&recipeservicesgithubrelease.List{},
		&recipeservicesgithubreleaseasset.Download{},
		&recipeservicesgithubreleaseasset.List{},
		&recipeservicesgithubreleaseasset.Upload{},
		&recipeservicesgithubtag.Create{},
		&recipeservicesgooglemailfilter.Add{},
		&recipeservicesgooglemailfilter.Delete{},
		&recipeservicesgooglemailfilter.List{},
		&recipeservicesgooglemailfilterbatch.Add{},
		&recipeservicesgooglemaillabel.Add{},
		&recipeservicesgooglemaillabel.Delete{},
		&recipeservicesgooglemaillabel.List{},
		&recipeservicesgooglemaillabel.Rename{},
		&recipeservicesgooglemailmessage.List{},
		&recipeservicesgooglemailmessagelabel.Add{},
		&recipeservicesgooglemailmessagelabel.Delete{},
		&recipeservicesgooglemailmessageprocessed.List{},
		&recipeservicesgooglemailthread.List{},
		&recipeservicesslackconversation.List{},
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
		&recipeteamcontentmember.List{},
		&recipeteamcontentmount.List{},
		&recipeteamcontentpolicy.List{},
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
		&recipeteamreport.Activity{},
		&recipeteamreport.Devices{},
		&recipeteamreport.Membership{},
		&recipeteamreport.Storage{},
		&recipeteamsharedlink.List{},
		&recipeteamsharedlinkupdate.Expiry{},
		&recipeteamfolder.Add{},
		&recipeteamfolder.Archive{},
		&recipeteamfolder.List{},
		&recipeteamfolder.Permdelete{},
		&recipeteamfolder.Replication{},
		&recipeteamfolderbatch.Archive{},
		&recipeteamfolderbatch.Permdelete{},
		&recipeteamfolderbatch.Replication{},
		&recipeteamfolderfile.List{},
		&recipeteamfolderfile.Size{},
		&recipeteamfoldermember.List{},
		&recipeteamfolderpolicy.List{},
	}
}
