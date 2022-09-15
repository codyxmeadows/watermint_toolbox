package rc_recipe

import (
	"github.com/watermint/toolbox/essentials/encoding/es_json"
	"github.com/watermint/toolbox/essentials/lang"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/data/da_griddata"
	"github.com/watermint/toolbox/infra/data/da_json"
	"github.com/watermint/toolbox/infra/data/da_text"
	"github.com/watermint/toolbox/infra/doc/dc_index"
	"github.com/watermint/toolbox/infra/doc/dc_recipe"
	"github.com/watermint/toolbox/infra/feed/fd_file"
	"github.com/watermint/toolbox/infra/recipe/rc_error_handler"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_ui"
)

type Spec interface {
	SpecValue

	// Recipe name
	Name() string

	// Recipe title
	Title() app_msg.Message

	// Recipe description
	Desc() app_msg.MessageOptional

	// Recipe remarks
	Remarks() app_msg.MessageOptional

	// Path signature of the recipe
	Path() (path []string, name string)

	// Id of the recipe spec. Id format is path + name of Path() connected with `-` (dash).
	SpecId() string

	// Command name and link to the document
	CliNameRef(media dc_index.MediaType, lg lang.Lang, relPath string) app_msg.Message

	// Recipe path on cli
	CliPath() string

	// Recipe argument on cli
	CliArgs() app_msg.MessageOptional

	// Notes for the recipe on cli
	CliNote() app_msg.MessageOptional

	// Spec of reports generated by this recipe
	Reports() []rp_model.Spec

	// Spec of feeds
	Feeds() map[string]fd_file.Spec

	// Spec of grid data input
	GridDataInput() map[string]da_griddata.GridDataInputSpec

	// Spec of grid data output
	GridDataOutput() map[string]da_griddata.GridDataOutputSpec

	// Spec of text input
	TextInput() map[string]da_text.TextInputSpec

	// Spec of json input
	JsonInput() map[string]da_json.JsonInputSpec

	// Messages used by this recipe
	Messages() []app_msg.Message

	// Returns a list of services used by this recipe
	Services() []string

	// True if this recipe use connection to the Dropbox Personal account
	ConnUsePersonal() bool

	// True if this recipe use connection to the Dropbox Business account
	ConnUseBusiness() bool

	// Returns array of scope of connections to Dropbox account(s)
	ConnScopes() []string

	// Field name and scope label map
	ConnScopeMap() map[string]string

	// Serialize
	Capture(ctl app_control.Control) (v interface{}, err error)

	// Deserialize & spin up
	Restore(j es_json.Json, ctl app_control.Control) (rcp Recipe, err error)

	// Apply values to the new recipe instance
	SpinUp(ctl app_control.Control, custom func(r Recipe)) (rcp Recipe, err error)

	// Serialize values
	Debug() map[string]interface{}

	// SpinDown
	SpinDown(ctl app_control.Control) error

	// True if the recipe is not for general usage.
	IsSecret() bool

	// True if the recipe is not designed for non-console UI.
	IsConsole() bool

	// True if the recipe is in experimental phase.
	IsExperimental() bool

	// True if the operation is irreversible.
	IsIrreversible() bool

	// True if the operation is transient.
	IsTransient() bool

	// IsDeprecated returns true if the operation is deprecated.
	IsDeprecated() bool

	// Print usage
	PrintUsage(ui app_ui.UI)

	// Create new spec
	New() Spec

	// Specification document
	Doc(ui app_ui.UI) *dc_recipe.Recipe

	// Error handlers for the recipe.
	ErrorHandlers() []rc_error_handler.ErrorHandler
}
