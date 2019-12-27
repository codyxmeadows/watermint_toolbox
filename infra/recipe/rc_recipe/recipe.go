package rc_recipe

import (
	"flag"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/feed/fd_file"
	"github.com/watermint/toolbox/infra/recipe/rc_vo"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/report/rp_spec"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_ui"
)

const (
	BasePackage = app.Pkg + "/recipe"
)

type Recipe interface {
	Exec(c app_control.Control) error
	Test(c app_control.Control) error
}

type SelfContainedRecipe interface {
	Recipe
	Preset()
}

// SecretRecipe will not be listed in available commands.
type SecretRecipe interface {
	Hidden()
}

// Console only recipe will not be listed in web console.
type ConsoleRecipe interface {
	Console()
}

type SpecValue interface {
	// Array of value names
	ValueNames() []string

	// Value description for the name
	ValueDesc(name string) app_msg.Message

	// Value default for the name
	ValueDefault(name string) interface{}

	// Value for the name
	Value(name string) Value

	// Customized value default for the name
	ValueCustomDefault(name string) app_msg.MessageOptional

	// Configure CLI flags
	SetFlags(f *flag.FlagSet, ui app_ui.UI)
}

type Spec interface {
	SpecValue

	// Recipe name
	Name() string

	// Recipe title
	Title() app_msg.Message

	// Recipe description
	Desc() app_msg.MessageOptional

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

	// Messages used by this recipe
	Messages() []app_msg.Message

	// True if this recipe use connection to the Dropbox Personal account
	ConnUsePersonal() bool

	// True if this recipe use connection to the Dropbox Business account
	ConnUseBusiness() bool

	// Returns array of scope of connections to Dropbox account(s)
	ConnScopes() []string

	// Field name and scope label map
	ConnScopeMap() map[string]string

	// Apply values to the new recipe instance
	SpinUp(ctl app_control.Control, custom func(r Recipe)) (rcp Recipe, err error)

	// Serialize values
	Debug() map[string]interface{}

	// SpinDown
	SpinDown(ctl app_control.Control) error
}

func NoCustomValues(r Recipe) {}
