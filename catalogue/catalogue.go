package catalogue

import (
	infra_api_api_api_auth_impl "github.com/watermint/toolbox/infra/api/api_auth_impl"
	infra_recipe_rc_conn_impl "github.com/watermint/toolbox/infra/recipe/rc_conn_impl"
	infra_recipe_rc_group "github.com/watermint/toolbox/infra/recipe/rc_group"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	infra_recipe_rcvalue "github.com/watermint/toolbox/infra/recipe/rc_value"
	infra_report_rpmodelimpl "github.com/watermint/toolbox/infra/report/rp_model_impl"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	infra_ui_appui "github.com/watermint/toolbox/infra/ui/app_ui"
	infra_util_ut_doc "github.com/watermint/toolbox/infra/util/ut_doc"
	ingredientfile "github.com/watermint/toolbox/ingredient/file"
	ingredientteamnamespacefile "github.com/watermint/toolbox/ingredient/team/namespace/file"
	ingredientteamfolder "github.com/watermint/toolbox/ingredient/teamfolder"
	"github.com/watermint/toolbox/recipe"
	recipedev "github.com/watermint/toolbox/recipe/dev"
	recipedevtest "github.com/watermint/toolbox/recipe/dev/test"
	recipefile "github.com/watermint/toolbox/recipe/file"
	recipefilecompare "github.com/watermint/toolbox/recipe/file/compare"
	recipefileimport "github.com/watermint/toolbox/recipe/file/import"
	recipefileimportbatch "github.com/watermint/toolbox/recipe/file/import/batch"
	recipefilesync "github.com/watermint/toolbox/recipe/file/sync"
	recipefilesyncpreflight "github.com/watermint/toolbox/recipe/file/sync/preflight"
	recipegroup "github.com/watermint/toolbox/recipe/group"
	recipegroupbatch "github.com/watermint/toolbox/recipe/group/batch"
	recipegroupmember "github.com/watermint/toolbox/recipe/group/member"
	recipemember "github.com/watermint/toolbox/recipe/member"
	recipememberquota "github.com/watermint/toolbox/recipe/member/quota"
	recipememberupdate "github.com/watermint/toolbox/recipe/member/update"
	recipesharedfolder "github.com/watermint/toolbox/recipe/sharedfolder"
	recipesharedfoldermember "github.com/watermint/toolbox/recipe/sharedfolder/member"
	recipesharedlink "github.com/watermint/toolbox/recipe/sharedlink"
	recipeteam "github.com/watermint/toolbox/recipe/team"
	recipeteamactivity "github.com/watermint/toolbox/recipe/team/activity"
	recipeteamactivitydaily "github.com/watermint/toolbox/recipe/team/activity/daily"
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

func NewCatalogue() infra_recipe_rc_group.Catalogue {
	return infra_recipe_rc_group.NewCatalogue(Recipes(), Ingredients(), Messages())
}

func Recipes() []rc_recipe.Recipe {
	cat := []rc_recipe.Recipe{
		&recipe.License{},
		&recipedev.Async{},
		&recipedev.Doc{},
		&recipedev.Dummy{},
		&recipedev.Preflight{},
		&recipedevtest.Auth{},
		&recipedevtest.Recipe{},
		&recipedevtest.Resources{},
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
		&recipefileimport.Url{},
		&recipefileimportbatch.Url{},
		&recipefilesync.Up{},
		&recipefilesyncpreflight.Up{},
		&recipegroup.Delete{},
		&recipegroup.List{},
		&recipegroupbatch.Delete{},
		&recipegroupmember.List{},
		&recipemember.Delete{},
		&recipemember.Detach{},
		&recipemember.Invite{},
		&recipemember.List{},
		&recipemember.Replication{},
		&recipememberquota.List{},
		&recipememberquota.Update{},
		&recipememberquota.Usage{},
		&recipememberupdate.Email{},
		&recipememberupdate.Externalid{},
		&recipememberupdate.Profile{},
		&recipesharedfolder.List{},
		&recipesharedfoldermember.List{},
		&recipesharedlink.Create{},
		&recipesharedlink.Delete{},
		&recipesharedlink.List{},
		&recipeteam.Feature{},
		&recipeteam.Info{},
		&recipeteamactivity.Event{},
		&recipeteamactivity.User{},
		&recipeteamactivitydaily.Event{},
		&recipeteamdevice.List{},
		&recipeteamdevice.Unlink{},
		&recipeteamdiag.Explorer{},
		&recipeteamfilerequest.List{},
		&recipeteamfolder.Archive{},
		&recipeteamfolder.List{},
		&recipeteamfolder.Permdelete{},
		&recipeteamfolder.Replication{},
		&recipeteamfolderbatch.Archive{},
		&recipeteamfolderbatch.Permdelete{},
		&recipeteamfolderbatch.Replication{},
		&recipeteamfolderfile.List{},
		&recipeteamfolderfile.Size{},
		&recipeteamlinkedapp.List{},
		&recipeteamnamespace.List{},
		&recipeteamnamespacefile.List{},
		&recipeteamnamespacefile.Size{},
		&recipeteamnamespacemember.List{},
		&recipeteamsharedlink.List{},
		&recipeteamsharedlinkupdate.Expiry{},
		//		&recipe.Web{},
	}
	return cat
}

func Ingredients() []rc_recipe.Recipe {
	cat := []rc_recipe.Recipe{
		&ingredientfile.Upload{},
		&ingredientteamfolder.Replication{},
		&ingredientteamnamespacefile.List{},
		&ingredientteamnamespacefile.Size{},
	}
	return cat
}

func Messages() []interface{} {
	msgs := []interface{}{
		infra_api_api_api_auth_impl.MCcAuth,
		infra_recipe_rc_group.MHeader,
		infra_recipe_rc_conn_impl.MConnect,
		infra_recipe_rcvalue.MRepository,
		infra_recipe_rcvalue.MValFdFileRowFeed,
		infra_report_rpmodelimpl.MTransactionReport,
		infra_report_rpmodelimpl.MXlsxWriter,
		infra_ui_appui.MConsole,
		infra_util_ut_doc.MDoc,
	}
	for _, m := range msgs {
		app_msg.Apply(m)
	}
	return msgs
}
