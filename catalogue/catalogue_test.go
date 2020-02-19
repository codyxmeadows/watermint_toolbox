package catalogue

import (
	"flag"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/recipe/rc_group"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/recipe/rc_spec"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
	"os"
	"testing"
)

func TestCatalogue(t *testing.T) {
	_, _, _, ui := qt_recipe.Resources(t)
	cat := NewCatalogue()
	testGroup(cat.RootGroup(), ui)
	for _, r := range cat.Ingredients() {
		spec := rc_spec.New(r)
		for _, m := range spec.Messages() {
			ui.Info(m)
		}
	}
}

func testGroup(g rc_group.Group, ui app_ui.UI) {
	g.PrintUsage(ui, os.Args[0], app.Version)
	for _, sg := range g.SubGroups() {
		testGroup(sg, ui)
	}
	for _, r := range g.Recipes() {
		testRecipe(r, ui)
	}
}

func testRecipe(r rc_recipe.Spec, ui app_ui.UI) {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	r.SetFlags(f, ui)
	r.PrintUsage(ui)

	for _, m := range r.Messages() {
		ui.Info(m)
	}
}
