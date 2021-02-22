package test

import (
	"errors"
	"github.com/watermint/toolbox/essentials/ambient/ea_indicator"
	"github.com/watermint/toolbox/essentials/ambient/ea_notification"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/essentials/model/mo_string"
	"github.com/watermint/toolbox/infra/control/app_catalogue"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/recipe/rc_spec"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/quality/recipe/qtr_timeout"
	"time"
)

var (
	recipeTimeout = 30 * time.Second
)

type Recipe struct {
	rc_recipe.RemarkSecret
	All       bool
	Single    mo_string.OptionalString
	NoTimeout bool
	Verbose   bool
	Result    rp_model.RowReport
}

func (z *Recipe) Preset() {
	z.Result.SetModel(&qtr_timeout.RecipeTestResult{})
}

func (z *Recipe) runSingle(c app_control.Control, r rc_recipe.Recipe) error {
	rs := rc_spec.New(r)
	path, name := rs.Path()
	l := c.Log().With(esl.Strings("path", path), esl.String("name", name))

	if rs.IsSecret() {
		l.Info("Skip secret recipe")
		return nil
	}

	rtr, err := qtr_timeout.RunRecipeTestWithTimeout(c, r, !z.NoTimeout, false)
	z.Result.Row(&rtr)
	return err
}

func (z *Recipe) runAll(c app_control.Control) error {
	cat := app_catalogue.Current()
	l := c.Log()

	ea_notification.Global().SuppressForce()
	ea_indicator.SuppressIndicatorForce()
	for _, r := range cat.Recipes() {
		if err := z.runSingle(c, r); err != nil {
			return err
		}
	}
	l.Info("All tests passed without error")
	return nil
}

func (z *Recipe) Exec(c app_control.Control) error {
	cat := app_catalogue.Current()
	l := c.Log()
	if err := z.Result.Open(); err != nil {
		return err
	}

	switch {
	case z.All:
		if err := z.runAll(c); err != nil {
			return err
		}

	case z.Single.IsExists():
		for _, r := range cat.Recipes() {
			p := rc_recipe.Key(r)
			if p != z.Single.Value() {
				continue
			}
			if err := z.runSingle(c, r); err != nil {
				return err
			}
			l.Info("Recipe test success")
			return nil
		}
		l.Error("recipe not found", esl.String("vo.Recipe", z.Single.Value()))
		return errors.New("recipe not found")

	default:
		l.Error("require -all or -single option")
		return errors.New("missing option")
	}
	return nil
}

func (z *Recipe) Test(c app_control.Control) error {
	return rc_exec.Exec(c, &Recipe{}, func(r rc_recipe.Recipe) {
		m := r.(*Recipe)
		m.All = true
		m.NoTimeout = false
	})
}
