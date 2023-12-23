package catalogue

import (
	"github.com/watermint/toolbox/infra/recipe/rc_spec"
	"testing"
)

func TestAutoDetectedRecipesClassic(t *testing.T) {
	ad := AutoDetectedRecipesClassic()
	for _, a := range ad {
		s := rc_spec.New(a)
		s.Name()
	}
}
