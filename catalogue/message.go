package catalogue

// Code generated by dev catalogue command DO NOT EDIT

import (
	domaindropboxapidbx_auth_attr "github.com/watermint/toolbox/domain/dropbox/api/dbx_auth_attr"
	domaindropboxapidbx_conn_impl "github.com/watermint/toolbox/domain/dropbox/api/dbx_conn_impl"
	domaindropboxapidbx_error "github.com/watermint/toolbox/domain/dropbox/api/dbx_error"
	domaindropboxapidbx_list_impl "github.com/watermint/toolbox/domain/dropbox/api/dbx_list_impl"
	domaindropboxfilesystem "github.com/watermint/toolbox/domain/dropbox/filesystem"
	domaindropboxmodelmo_file_filter "github.com/watermint/toolbox/domain/dropbox/model/mo_file_filter"
	domaindropboxmodelmo_sharedfolder_member "github.com/watermint/toolbox/domain/dropbox/model/mo_sharedfolder_member"
	domaindropboxusecaseuc_compare_local "github.com/watermint/toolbox/domain/dropbox/usecase/uc_compare_local"
	domaindropboxusecaseuc_compare_paths "github.com/watermint/toolbox/domain/dropbox/usecase/uc_compare_paths"
	domaindropboxusecaseuc_file_merge "github.com/watermint/toolbox/domain/dropbox/usecase/uc_file_merge"
	domaindropboxusecaseuc_file_relocation "github.com/watermint/toolbox/domain/dropbox/usecase/uc_file_relocation"
	domaingooglemailservicesv_label "github.com/watermint/toolbox/domain/google/mail/service/sv_label"
	domaingooglemailservicesv_message "github.com/watermint/toolbox/domain/google/mail/service/sv_message"
	essentialslogesl_rotate "github.com/watermint/toolbox/essentials/log/esl_rotate"
	essentialsmodelmo_filter "github.com/watermint/toolbox/essentials/model/mo_filter"
	essentialsnetworknw_diag "github.com/watermint/toolbox/essentials/network/nw_diag"
	infraapiapi_auth_oauth "github.com/watermint/toolbox/infra/api/api_auth_oauth"
	infraapiapi_callback "github.com/watermint/toolbox/infra/api/api_callback"
	infracontrolapp_error "github.com/watermint/toolbox/infra/control/app_error"
	infracontrolapp_job_impl "github.com/watermint/toolbox/infra/control/app_job_impl"
	infradatada_griddata "github.com/watermint/toolbox/infra/data/da_griddata"
	infradatada_json "github.com/watermint/toolbox/infra/data/da_json"
	infradocdc_license "github.com/watermint/toolbox/infra/doc/dc_license"
	infradocdc_options "github.com/watermint/toolbox/infra/doc/dc_options"
	infradocdc_supplemental "github.com/watermint/toolbox/infra/doc/dc_supplemental"
	infrafeedfd_file_impl "github.com/watermint/toolbox/infra/feed/fd_file_impl"
	infrareciperc_exec "github.com/watermint/toolbox/infra/recipe/rc_exec"
	infrareciperc_group "github.com/watermint/toolbox/infra/recipe/rc_group"
	infrareciperc_group_impl "github.com/watermint/toolbox/infra/recipe/rc_group_impl"
	infrareciperc_spec "github.com/watermint/toolbox/infra/recipe/rc_spec"
	infrareciperc_value "github.com/watermint/toolbox/infra/recipe/rc_value"
	infrareportrp_model_impl "github.com/watermint/toolbox/infra/report/rp_model_impl"
	infrareportrp_writer_impl "github.com/watermint/toolbox/infra/report/rp_writer_impl"
	infrauiapp_ui "github.com/watermint/toolbox/infra/ui/app_ui"
	ingredientfile "github.com/watermint/toolbox/ingredient/file"
	recipedevdiag "github.com/watermint/toolbox/recipe/dev/diag"
	recipefiledispatch "github.com/watermint/toolbox/recipe/file/dispatch"
	recipefileimportbatch "github.com/watermint/toolbox/recipe/file/import/batch"
	recipegroupmember "github.com/watermint/toolbox/recipe/group/member"
	recipegroupmemberbatch "github.com/watermint/toolbox/recipe/group/member/batch"
	recipemember "github.com/watermint/toolbox/recipe/member"
	recipememberquota "github.com/watermint/toolbox/recipe/member/quota"
	recipememberupdate "github.com/watermint/toolbox/recipe/member/update"
	recipeservicesgithubreleaseasset "github.com/watermint/toolbox/recipe/services/github/release/asset"
	recipeteamactivity "github.com/watermint/toolbox/recipe/team/activity"
	recipeteamactivitybatch "github.com/watermint/toolbox/recipe/team/activity/batch"
	recipeteamdevice "github.com/watermint/toolbox/recipe/team/device"
	recipeteamsharedlink "github.com/watermint/toolbox/recipe/team/sharedlink"
)

func AutoDetectedMessageObjects() []interface{} {
	return []interface{}{
		&domaindropboxapidbx_auth_attr.MsgAttr{},
		&domaindropboxapidbx_conn_impl.MsgConnect{},
		&domaindropboxapidbx_error.MsgHandler{},
		&domaindropboxapidbx_list_impl.MsgList{},
		&domaindropboxfilesystem.MsgFileSystemCached{},
		&domaindropboxmodelmo_file_filter.MsgFileFilterOpt{},
		&domaindropboxmodelmo_sharedfolder_member.MsgExternalOpt{},
		&domaindropboxmodelmo_sharedfolder_member.MsgInternalOpt{},
		&domaindropboxusecaseuc_compare_local.MsgCompare{},
		&domaindropboxusecaseuc_compare_paths.MsgCompare{},
		&domaindropboxusecaseuc_file_merge.MsgMerge{},
		&domaindropboxusecaseuc_file_relocation.MsgRelocation{},
		&domaingooglemailservicesv_label.MsgFindLabel{},
		&domaingooglemailservicesv_message.MsgProgress{},
		&essentialslogesl_rotate.MsgOut{},
		&essentialslogesl_rotate.MsgPurge{},
		&essentialslogesl_rotate.MsgRotate{},
		&essentialsmodelmo_filter.MsgFilter{},
		&essentialsnetworknw_diag.MsgNetwork{},
		&infraapiapi_auth_oauth.MsgApiAuth{},
		&infraapiapi_callback.MsgCallback{},
		&infracontrolapp_error.MsgErrorReport{},
		&infracontrolapp_job_impl.MsgLauncher{},
		&infradatada_griddata.MsgGridDataInput{},
		&infradatada_json.MsgJsonInput{},
		&infradocdc_license.MsgLicense{},
		&infradocdc_options.MsgDoc{},
		&infradocdc_supplemental.MsgDeveloper{},
		&infradocdc_supplemental.MsgDropboxBusiness{},
		&infradocdc_supplemental.MsgExperimentalFeature{},
		&infradocdc_supplemental.MsgPathVariable{},
		&infradocdc_supplemental.MsgTroubleshooting{},
		&infrafeedfd_file_impl.MsgRowFeed{},
		&infrareciperc_exec.MsgPanic{},
		&infrareciperc_group.MsgHeader{},
		&infrareciperc_group_impl.MsgGroup{},
		&infrareciperc_spec.MsgSelfContained{},
		&infrareciperc_value.MsgRepository{},
		&infrareciperc_value.MsgValFdFileRowFeed{},
		&infrareportrp_model_impl.MsgColumnSpec{},
		&infrareportrp_model_impl.MsgTransactionReport{},
		&infrareportrp_writer_impl.MsgSortedWriter{},
		&infrareportrp_writer_impl.MsgUIWriter{},
		&infrareportrp_writer_impl.MsgXlsxWriter{},
		&infrauiapp_ui.MsgConsole{},
		&infrauiapp_ui.MsgProgress{},
		&ingredientfile.MsgUpload{},
		&recipedevdiag.MsgLoader{},
		&recipefiledispatch.MsgLocal{},
		&recipefileimportbatch.MsgUrl{},
		&recipegroupmember.MsgList{},
		&recipegroupmemberbatch.MsgOperation{},
		&recipemember.MsgInvite{},
		&recipememberquota.MsgList{},
		&recipememberquota.MsgUpdate{},
		&recipememberquota.MsgUsage{},
		&recipememberupdate.MsgEmail{},
		&recipeservicesgithubreleaseasset.MsgUp{},
		&recipeteamactivity.MsgUser{},
		&recipeteamactivitybatch.MsgUser{},
		&recipeteamdevice.MsgUnlink{},
		&recipeteamsharedlink.MsgList{},
	}
}
