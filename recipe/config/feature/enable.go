package feature

import (
	"errors"
	"github.com/watermint/toolbox/infra/control/app_catalogue"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_feature"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/ui/app_msg"
)

type SampleFeature struct {
	app_feature.OptInStatus
}
type SampleFeatureNotInCatalogue struct {
	app_feature.OptInStatus
}

var (
	ErrorCatalogueIsNotAvailable = errors.New("catalogue is not available")
	ErrorInvalidKey              = errors.New("invalid key")
)

type Enable struct {
	Key                        string
	ErrorInvalidKey            app_msg.Message
	ErrorUnableToEnableFeature app_msg.Message
	InfoCancelled              app_msg.Message
	InfoOptIn                  app_msg.Message
}

func (z *Enable) Preset() {
}

func (z *Enable) Exec(c app_control.Control) error {
	ui := c.UI()
	cat := app_catalogue.Current()
	features := cat.Features()
	if c.Feature().IsTest() {
		features = append(features, &SampleFeature{})
	}
	var optIn app_feature.OptIn = nil
	for _, f := range features {
		if app_feature.OptInName(f) == z.Key {
			optIn = f
		}
	}
	if optIn == nil {
		ui.Error(z.ErrorInvalidKey.With("Key", z.Key))
		return ErrorInvalidKey
	}

	ui.Info(app_feature.OptInDescription(optIn))
	cont := ui.AskCont(app_feature.OptInAgreement(optIn))
	if !cont {
		ui.Info(z.InfoCancelled)
		return nil
	}
	optIn.OptInCommit(true)
	if err := c.Feature().OptInUpdate(optIn); err != nil {
		ui.Error(z.ErrorUnableToEnableFeature.With("Key", z.Key))
		return err
	}
	ui.Info(z.InfoOptIn.With("Key", z.Key))
	return nil
}

func (z *Enable) Test(c app_control.Control) error {
	if err := rc_exec.Exec(c, &Enable{}, func(r rc_recipe.Recipe) {
		f := &SampleFeatureNotInCatalogue{}
		m := r.(*Enable)
		m.Key = app_feature.OptInName(f)
	}); err != ErrorInvalidKey {
		return ErrorInvalidKey
	}

	if err := rc_exec.Exec(c, &Enable{}, func(r rc_recipe.Recipe) {
		f := &SampleFeature{}
		m := r.(*Enable)
		m.Key = app_feature.OptInName(f)
	}); err != nil {
		return err
	}
	return nil
}
