package rc_spec

import (
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
)

func New(rcp rc_recipe.Recipe) rc_recipe.Spec {
	return newSelfContained(rcp)
}
